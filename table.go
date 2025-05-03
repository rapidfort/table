// Package table provides a utility for rendering ASCII tables with box-drawing characters
package table

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// Box drawing characters
const (
	TopLeft     = "┌"
	TopRight    = "┐"
	BottomLeft  = "└"
	BottomRight = "┘"
	HLine       = "─"
	VLine       = "│"
	LeftT       = "├"
	RightT      = "┤"
	TopT        = "┬"
	BottomT     = "┴"
	Cross       = "┼"
	padding     = 1
)

// Table represents a table with borders and alignment control
type Table struct {
	Headers      []string
	Rows         [][]string
	Descriptions map[int]string
	columnWidths []int
	alignments   []string // "left", "right", "center" for each column
	consoleWidth int      // Maximum width of the console
	fillWidth    bool
	maxWidths    map[int]int // Maximum width for specific columns
	// Reference to the table group this table belongs to (if any)
	group *TableGroup
}

// TableGroup manages multiple tables with consistent column widths
type TableGroup struct {
	tables       []*Table
	columnWidths []int
}

// NewGroup creates a new TableGroup for managing multiple tables
func NewGroup() *TableGroup {
	return &TableGroup{
		tables: []*Table{},
	}
}

// Add adds a table to the group
func (g *TableGroup) Add(table *Table) {
	g.tables = append(g.tables, table)
	table.group = g
}

// SyncColumnWidths ensures all tables in the group have consistent column widths
func (g *TableGroup) SyncColumnWidths() {
	if len(g.tables) == 0 {
		return
	}

	// Initialize with the first table's column count
	firstTable := g.tables[0]
	colCount := len(firstTable.columnWidths)

	// Initialize the group's column widths
	g.columnWidths = make([]int, colCount)

	// Find the maximum width for each column across all tables
	for _, table := range g.tables {
		table.calculateInitialColumnWidths()

		for i := 0; i < colCount && i < len(table.columnWidths); i++ {
			if table.columnWidths[i] > g.columnWidths[i] {
				// Check if this column has a max width constraint
				if maxWidth, exists := table.maxWidths[i]; exists && table.columnWidths[i] > maxWidth {
					g.columnWidths[i] = maxWidth
				} else {
					g.columnWidths[i] = table.columnWidths[i]
				}
			}
		}
	}

	// Apply the group's column widths to all tables
	for _, table := range g.tables {
		for i := 0; i < colCount && i < len(table.columnWidths); i++ {
			table.columnWidths[i] = g.columnWidths[i]
		}

		// Apply any final adjustments needed for console width
		table.adjustColumnWidthsToFit()
	}
}

// SetFillWidth sets whether the table should expand to fill the console width
func (t *Table) SetFillWidth(enabled bool) {
	t.fillWidth = enabled
}

// SetConsoleWidth sets the maximum width for the table
func (t *Table) SetConsoleWidth(width int) {
	t.consoleWidth = width
}

// SetAlignment sets the alignment for a specific column
func (t *Table) SetAlignment(columnIndex int, alignment string) {
	if columnIndex >= 0 && columnIndex < len(t.alignments) {
		t.alignments[columnIndex] = alignment
	}
}

// SetMaxWidth sets the maximum width for a specific column
func (t *Table) SetMaxWidth(columnIndex int, maxWidth int) {
	if columnIndex >= 0 && columnIndex < len(t.Headers) {
		t.maxWidths[columnIndex] = maxWidth
	}
}

// AddRow adds a new row to the table
func (t *Table) AddRow(row []string) {
	for len(row) < len(t.Headers) {
		row = append(row, "")
	}
	t.Rows = append(t.Rows, row)
}

// AddDescription adds a description for a specific row
func (t *Table) AddDescription(rowIndex int, description string) {
	if rowIndex >= 0 && rowIndex < len(t.Rows) {
		t.Descriptions[rowIndex] = description
	}
}

// calculateInitialColumnWidths computes the initial width for each column
func (t *Table) calculateInitialColumnWidths() {
	for i, header := range t.Headers {
		if l := utf8.RuneCountInString(header); l > t.columnWidths[i] {
			t.columnWidths[i] = l
		}
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(t.columnWidths) {
				if l := utf8.RuneCountInString(cell); l > t.columnWidths[i] {
					t.columnWidths[i] = l
				}
			}
		}
	}

	// Apply max widths if specified
	for i, width := range t.columnWidths {
		if maxWidth, exists := t.maxWidths[i]; exists && width > maxWidth {
			t.columnWidths[i] = maxWidth
		}
	}
}

// adjustColumnWidthsToFit adjusts column widths to fit the console
func (t *Table) adjustColumnWidthsToFit() {
	// --- a) SHRINK if too wide ---
	total := 1
	for _, w := range t.columnWidths {
		total += w + 2 + 1
	}
	if total > t.consoleWidth {
		excess := total - t.consoleWidth
		for excess > 0 {
			maxW, idx := 0, -1
			for i, w := range t.columnWidths {
				if w > maxW && w > 3 {
					maxW, idx = w, i
				}
			}
			if idx < 0 {
				break
			}
			t.columnWidths[idx]--
			excess--
		}
	}

	// --- b) EXPAND if fillWidth is set and still under consoleWidth ---
	if t.fillWidth {
		total = 1
		for _, w := range t.columnWidths {
			total += w + 2 + 1
		}
		if total < t.consoleWidth {
			extra := t.consoleWidth - total

			// Count columns that can be expanded (exclude those with max width)
			expandableCols := 0
			for i := range t.columnWidths {
				if maxWidth, exists := t.maxWidths[i]; !exists || t.columnWidths[i] < maxWidth {
					expandableCols++
				}
			}

			if expandableCols > 0 {
				per := extra / expandableCols
				rem := extra % expandableCols

				expandedCount := 0
				for i := range t.columnWidths {
					// Skip columns that have reached their max width
					if maxWidth, exists := t.maxWidths[i]; exists && t.columnWidths[i] >= maxWidth {
						continue
					}

					t.columnWidths[i] += per
					if expandedCount < rem {
						t.columnWidths[i]++
						expandedCount++
					}

					// Ensure we don't exceed max width after expansion
					if maxWidth, exists := t.maxWidths[i]; exists && t.columnWidths[i] > maxWidth {
						t.columnWidths[i] = maxWidth
					}
				}
			}
		}
	}
}

// Calculate total width of the table
func (t *Table) calculateTotalWidth() int {
	width := 1 // Start with left border
	for _, w := range t.columnWidths {
		width += w + 2 + 1 // column width + padding + separator
	}
	return width
}

// smartSplitCellContent intelligently splits cell content
func (t *Table) smartSplitCellContent(content string, colIndex int) []string {
	maxWidth := t.columnWidths[colIndex]
	if utf8.RuneCountInString(content) <= maxWidth {
		return []string{content}
	}
	if strings.Contains(content, "/") || strings.Contains(content, ".") {
		return t.splitPackageName(content, maxWidth)
	}
	if strings.Contains(content, ",") {
		return t.splitCommaSeparatedList(content, maxWidth)
	}
	return t.splitByWords(content, maxWidth)
}

// splitPackageName splits package names and URLs at boundaries
func (t *Table) splitPackageName(content string, maxWidth int) []string {
	parts := strings.Split(content, "/")
	var res []string
	line := ""
	for i, part := range parts {
		if i > 0 {
			part = "/" + part
		}
		if line != "" && utf8.RuneCountInString(line+part) > maxWidth {
			res = append(res, line)
			line = part
		} else {
			line += part
		}
		if utf8.RuneCountInString(line) > maxWidth {
			chunks := t.splitByWords(line, maxWidth)
			res = append(res, chunks...)
			line = ""
		}
	}
	if line != "" {
		res = append(res, line)
	}
	return res
}

// splitCommaSeparatedList breaks at commas, keeping items under maxWidth
func (t *Table) splitCommaSeparatedList(content string, maxWidth int) []string {
	parts := strings.Split(content, ",")
	var res []string
	line := ""
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if i > 0 {
			part = ", " + part
		}
		if line != "" && utf8.RuneCountInString(line+part) > maxWidth {
			res = append(res, line)
			line = strings.TrimPrefix(part, ", ")
		} else {
			line += part
		}
	}
	if line != "" {
		res = append(res, line)
	}
	return res
}

// splitByWords splits on spaces to keep each line under maxWidth
func (t *Table) splitByWords(content string, maxWidth int) []string {
	words := strings.Fields(content)
	var res []string
	line := ""
	for _, w := range words {
		test := w
		if line != "" {
			test = line + " " + w
		}
		if utf8.RuneCountInString(test) <= maxWidth {
			line = test
		} else {
			if line != "" {
				res = append(res, line)
			}
			line = w
		}
		if utf8.RuneCountInString(line) > maxWidth {
			part := line[:maxWidth]
			res = append(res, part)
			line = line[maxWidth:]
		}
	}
	if line != "" {
		res = append(res, line)
	}
	return res
}

// formatCellContent formats a cell's content with alignment and padding
func (t *Table) formatCellContent(content string, colIndex int) string {
	w := t.columnWidths[colIndex]
	switch t.alignments[colIndex] {
	case "right":
		return fmt.Sprintf(" %*s ", w, content)
	case "center":
		totalPad := w - utf8.RuneCountInString(content)
		left := totalPad / 2
		return fmt.Sprintf(" %s%s%s ",
			strings.Repeat(" ", left), content,
			strings.Repeat(" ", totalPad-left))
	default:
		return fmt.Sprintf(" %-*s ", w, content)
	}
}

// calculateTableWidth returns the total width of the table
func (t *Table) calculateTableWidth() int {
	total := 1
	for _, w := range t.columnWidths {
		total += w + 2 + 1
	}
	return total
}

// Border renderers

func (t *Table) renderTopBorder() string {
	var sb strings.Builder
	sb.WriteString(TopLeft)
	for i, w := range t.columnWidths {
		sb.WriteString(strings.Repeat(HLine, w+2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(TopT)
		}
	}
	sb.WriteString(TopRight + "\n")
	return sb.String()
}

func (t *Table) renderMiddleBorder() string {
	var sb strings.Builder
	sb.WriteString(LeftT)
	for i, w := range t.columnWidths {
		sb.WriteString(strings.Repeat(HLine, w+2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(Cross)
		}
	}
	sb.WriteString(RightT + "\n")
	return sb.String()
}

// renderDataToDescBorder renders the border between a data row and its description
// Uses BottomT (┴) for column separators instead of Cross (┼)
func (t *Table) renderDataToDescBorder() string {

	var sb strings.Builder
	sb.WriteString(LeftT)
	sb.WriteString(strings.Repeat(HLine, t.columnWidths[0]+2))
	sb.WriteString(Cross)
	merged := 0
	for i := 1; i < len(t.columnWidths); i++ {
		merged += t.columnWidths[i] + 2
	}
	sb.WriteString(strings.Repeat(HLine, merged))
	sb.WriteString(RightT + "\n")
	return sb.String()
}

// renderLastRowBottomBorder renders the bottom border for the last row when it has a description
// Uses continuous HLine (─) with no column separators
func (t *Table) renderLastRowBottomBorder() string {
	var sb strings.Builder
	sb.WriteString(BottomLeft)

	// Calculate total internal width (all columns + separators)
	totalWidth := 0
	for _, w := range t.columnWidths {
		totalWidth += w + 2 // width + padding on both sides
	}
	// Add separators between columns
	totalWidth += len(t.columnWidths) - 1

	// Draw a continuous line across the entire width with no column separators
	sb.WriteString(strings.Repeat(HLine, totalWidth))

	sb.WriteString(BottomRight + "\n")
	return sb.String()
}

// wrapDescriptionText wraps the description text to fit within the table width
func (t *Table) wrapDescriptionText(text string) []string {
	// Calculate available width for the description
	totalWidth := t.calculateTableWidth() - 4 // Subtract borders and padding

	// If the text contains newlines, split by newlines first
	if strings.Contains(text, "\n") {
		var result []string
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			wrapped := t.splitByWords(line, totalWidth)
			result = append(result, wrapped...)
		}
		return result
	}

	// Otherwise wrap the text to fit the width
	return t.splitByWords(text, totalWidth)
}

// renderDescriptionRow renders a description row with proper formatting
func (t *Table) renderDescriptionRow(text string, isFirstLine bool) string {
	var sb strings.Builder
	totalWidth := t.calculateTableWidth() - 2 // Width minus left and right border

	sb.WriteString(VLine)

	if isFirstLine {
		prefix := ""
		sb.WriteString(prefix)

		// Calculate remaining space and padding
		contentWidth := utf8.RuneCountInString(text)
		remainingSpace := totalWidth - utf8.RuneCountInString(prefix)

		if contentWidth <= remainingSpace {
			// Text fits on the line
			sb.WriteString(text)
			sb.WriteString(strings.Repeat(" ", remainingSpace-contentWidth))
		} else {
			// Text needs to be truncated
			sb.WriteString(text[:remainingSpace])
		}
	} else {
		prefix := "   • "
		sb.WriteString(prefix)

		// Calculate remaining space and padding
		contentWidth := utf8.RuneCountInString(text)
		remainingSpace := totalWidth - utf8.RuneCountInString(prefix)

		if contentWidth <= remainingSpace {
			// Text fits on the line
			sb.WriteString(text)
			sb.WriteString(strings.Repeat(" ", remainingSpace-contentWidth))
		} else {
			// Text needs to be truncated
			sb.WriteString(text[:remainingSpace])
		}
	}

	sb.WriteString(VLine + "\n")
	return sb.String()
}

func (t *Table) renderBottomBorder() string {
	var sb strings.Builder
	sb.WriteString(BottomLeft)
	for i, w := range t.columnWidths {
		sb.WriteString(strings.Repeat(HLine, w+2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(BottomT)
		}
	}
	sb.WriteString(BottomRight + "\n")
	return sb.String()
}

// // Helper function to check if a row has a description
// func (t *Table) hasDescription(rowIndex int) bool {
// 	_, hasDesc := t.Descriptions[rowIndex]
// 	return hasDesc
// }

// // Helper function to check if the next row exists and has a description
// func (t *Table) nextRowHasDescription(rowIndex int) bool {
// 	if rowIndex+1 >= len(t.Rows) {
// 		return false
// 	}
// 	_, hasDesc := t.Descriptions[rowIndex+1]
// 	return hasDesc
// }

// Add these functions to your existing table.go file

// Implementation of dynamic terminal width detection and text wrapping

// detectTerminalWidth gets the current terminal width or returns a default
func detectTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		// Default to 80 if we can't detect terminal width
		return 80
	}
	return width
}

// calculateOptimalColumnWidths distributes width to each column based on content
func (t *Table) calculateOptimalColumnWidths(maxWidth int) {
	// First, get minimum widths needed for each column
	t.calculateInitialColumnWidths()

	// Calculate total required width
	totalRequiredWidth := 1 // Start with left border
	for _, w := range t.columnWidths {
		// Add column width + padding + separator
		totalRequiredWidth += w + 2 + 1
	}

	// If total width exceeds available width, redistribute
	if totalRequiredWidth > maxWidth {
		// We need to shrink columns to fit
		excessWidth := totalRequiredWidth - maxWidth
		t.shrinkColumnsToFit(excessWidth)
	} else if t.fillWidth {
		// We have extra space and fillWidth is true, so expand columns
		extraWidth := maxWidth - totalRequiredWidth
		t.expandColumnsToFit(extraWidth)
	}
}

// shrinkColumnsToFit reduces column widths to fit within available space
func (t *Table) shrinkColumnsToFit(excessWidth int) {
	// Start by reducing the widest columns first
	for excessWidth > 0 {
		// Find the widest column that can be shrunk
		maxWidth, idx := 0, -1
		for i, w := range t.columnWidths {
			// Don't shrink below minimum usable width (3 chars)
			if w > maxWidth && w > 3 {
				maxWidth, idx = w, i
			}
		}

		if idx < 0 {
			// No more columns can be shrunk, we'll have to live with horizontal scrolling
			break
		}

		// Reduce the column width
		reduceBy := 1
		if excessWidth > 5 && t.columnWidths[idx] > 10 {
			// For large excesses, reduce by more to avoid many small reductions
			reduceBy = excessWidth / 5
			if reduceBy > (t.columnWidths[idx] - 3) {
				reduceBy = t.columnWidths[idx] - 3
			}
		}

		t.columnWidths[idx] -= reduceBy
		excessWidth -= reduceBy
	}
}

// expandColumnsToFit distributes extra space among columns
func (t *Table) expandColumnsToFit(extraWidth int) {
	// Count expandable columns (exclude those with max width constraints)
	expandableCols := 0
	for i := range t.columnWidths {
		if maxWidth, exists := t.maxWidths[i]; !exists || t.columnWidths[i] < maxWidth {
			expandableCols++
		}
	}

	if expandableCols > 0 {
		perColumn := extraWidth / expandableCols
		remainder := extraWidth % expandableCols

		expandedCount := 0
		for i := range t.columnWidths {
			// Skip columns that have reached their max width
			if maxWidth, exists := t.maxWidths[i]; exists && t.columnWidths[i] >= maxWidth {
				continue
			}

			t.columnWidths[i] += perColumn
			if expandedCount < remainder {
				t.columnWidths[i]++
				expandedCount++
			}

			// Ensure we don't exceed max width constraints
			if maxWidth, exists := t.maxWidths[i]; exists && t.columnWidths[i] > maxWidth {
				t.columnWidths[i] = maxWidth
			}
		}
	}
}

// wrapCellText wraps text to fit within the column width
func (t *Table) wrapCellText(text string, colIndex int) []string {
	maxWidth := t.columnWidths[colIndex]

	// If text is shorter than maxWidth, return as is
	if utf8.RuneCountInString(text) <= maxWidth {
		return []string{text}
	}

	// Handle special cases (URLs, package names, comma-separated lists)
	if strings.Contains(text, "/") || strings.Contains(text, ".") {
		return t.smartSplitPackageName(text, maxWidth)
	}
	if strings.Contains(text, ",") {
		return t.smartSplitCommaSeparatedList(text, maxWidth)
	}

	// Otherwise, split by words
	return t.smartSplitByWords(text, maxWidth)
}

// smartSplitPackageName splits package names and URLs at logical boundaries
func (t *Table) smartSplitPackageName(text string, maxWidth int) []string {
	parts := strings.Split(text, "/")
	var result []string
	currentLine := ""

	for i, part := range parts {
		// Add slash prefix for all but the first part
		if i > 0 {
			part = "/" + part
		}

		// Check if adding this part would exceed the max width
		if currentLine != "" && utf8.RuneCountInString(currentLine+part) > maxWidth {
			// Current line is full, add it to result and start a new one
			result = append(result, currentLine)
			currentLine = part
		} else {
			// Add to current line
			currentLine += part
		}

		// If current line itself exceeds max width, split it further
		if utf8.RuneCountInString(currentLine) > maxWidth {
			// Use word splitting as fallback
			wordSplit := t.smartSplitByWords(currentLine, maxWidth)

			// Add all but the last part to result
			if len(wordSplit) > 1 {
				result = append(result, wordSplit[:len(wordSplit)-1]...)
			}

			// Keep last part as current line
			currentLine = wordSplit[len(wordSplit)-1]
		}
	}

	// Add any remaining text
	if currentLine != "" {
		result = append(result, currentLine)
	}

	return result
}

// smartSplitCommaSeparatedList splits comma-separated lists at logical boundaries
func (t *Table) smartSplitCommaSeparatedList(text string, maxWidth int) []string {
	items := strings.Split(text, ",")
	var result []string
	currentLine := ""

	for i, item := range items {
		item = strings.TrimSpace(item)

		// Add comma prefix for all but the first item
		if i > 0 {
			item = ", " + item
		}

		// Check if adding this item would exceed the max width
		if currentLine != "" && utf8.RuneCountInString(currentLine+item) > maxWidth {
			// Current line is full, add it to result and start a new one
			result = append(result, currentLine)
			currentLine = strings.TrimPrefix(item, ", ")
		} else {
			// Add to current line
			currentLine += item
		}

		// If current line itself exceeds max width, split it further
		if utf8.RuneCountInString(currentLine) > maxWidth {
			// Use word splitting as fallback
			wordSplit := t.smartSplitByWords(currentLine, maxWidth)

			// Add all but the last part to result
			if len(wordSplit) > 1 {
				result = append(result, wordSplit[:len(wordSplit)-1]...)
			}

			// Keep last part as current line
			currentLine = wordSplit[len(wordSplit)-1]
		}
	}

	// Add any remaining text
	if currentLine != "" {
		result = append(result, currentLine)
	}

	return result
}

// smartSplitByWords splits text by words to fit within maxWidth
func (t *Table) smartSplitByWords(text string, maxWidth int) []string {
	words := strings.Fields(text)
	var result []string
	currentLine := ""

	for _, word := range words {
		// Check if adding this word would exceed the max width
		testLine := word
		if currentLine != "" {
			testLine = currentLine + " " + word
		}

		if utf8.RuneCountInString(testLine) <= maxWidth {
			// Word fits on current line
			currentLine = testLine
		} else {
			// Word doesn't fit, add current line to result and start new one
			if currentLine != "" {
				result = append(result, currentLine)
			}

			// If the word itself is longer than maxWidth, split it
			if utf8.RuneCountInString(word) > maxWidth {
				// Split the word into chunks of maxWidth
				for len(word) > 0 {
					if utf8.RuneCountInString(word) <= maxWidth {
						currentLine = word
						word = ""
					} else {
						// Find a good split point (max width characters)
						result = append(result, word[:maxWidth])
						word = word[maxWidth:]
					}
				}
			} else {
				currentLine = word
			}
		}
	}

	// Add any remaining text
	if currentLine != "" {
		result = append(result, currentLine)
	}

	return result
}

// Now modify the RapidFortTable constructor to use terminal width by default

// RapidFortTable creates a new Table with the given headers
func RapidFortTable(headers []string) *Table {
	// Auto-detect terminal width
	termWidth := detectTerminalWidth()

	table := &Table{
		Headers:      headers,
		Rows:         [][]string{},
		Descriptions: make(map[int]string),
		columnWidths: make([]int, len(headers)),
		alignments:   make([]string, len(headers)),
		consoleWidth: termWidth,
		fillWidth:    true, // Enable fill width by default
		maxWidths:    make(map[int]int),
	}

	// Set default left alignment for all columns
	for i := range headers {
		table.alignments[i] = "left"
	}

	return table
}

// Finally, update the Render method to use our new functions

// Render method with correct border management for merged descriptions
func (t *Table) Render() string {
	if t.group == nil {
		t.calculateOptimalColumnWidths(t.consoleWidth)
	}

	var sb strings.Builder

	// Top border
	sb.WriteString(t.renderTopBorder())

	// Headers
	headerLines := make([][]string, len(t.Headers))
	for i, h := range t.Headers {
		headerLines[i] = t.smartSplitCellContent(h, i)
	}
	maxH := 0
	for _, lines := range headerLines {
		if len(lines) > maxH {
			maxH = len(lines)
		}
	}
	for line := 0; line < maxH; line++ {
		sb.WriteString(VLine)
		for ci := range t.Headers {
			txt := ""
			if line < len(headerLines[ci]) {
				txt = headerLines[ci][line]
			}
			sb.WriteString(t.formatCellContent(txt, ci))
			sb.WriteString(VLine)
		}
		sb.WriteString("\n")
	}

	// Header/Data separator
	sb.WriteString(t.renderMiddleBorder())

	// Rows
	for ri, row := range t.Rows {
		// Print the data row
		rowLines := make([][]string, len(row))
		maxR := 0
		for ci, cell := range row {
			rowLines[ci] = t.smartSplitCellContent(cell, ci)
			if len(rowLines[ci]) > maxR {
				maxR = len(rowLines[ci])
			}
		}

		for line := 0; line < maxR; line++ {
			sb.WriteString(VLine)
			for ci := range row {
				txt := ""
				if line < len(rowLines[ci]) {
					txt = rowLines[ci][line]
				}
				sb.WriteString(t.formatCellContent(txt, ci))
				sb.WriteString(VLine)
			}
			sb.WriteString("\n")
		}

		// Check for description
		if desc, ok := t.Descriptions[ri]; ok {
			// Process the description
			bulletPoints := strings.Split(desc, "\n")
			if len(bulletPoints) == 1 && strings.Contains(desc, ",") {
				bulletPoints = strings.Split(desc, ",")
				for i := range bulletPoints {
					bulletPoints[i] = strings.TrimSpace(bulletPoints[i])
				}
			}

			// Calculate merged width (columns 2-n)
			mergedWidth := 0
			for i := 1; i < len(t.columnWidths); i++ {
				mergedWidth += t.columnWidths[i] + 2
				if i < len(t.columnWidths)-1 {
					mergedWidth += 1 // Add separator space
				}
			}

			// Render the border that merges columns 2-n
			sb.WriteString(VLine)
			sb.WriteString(t.formatCellContent("", 0))
			sb.WriteString(LeftT)

			// Draw border for all merged columns
			for i := 1; i < len(t.columnWidths); i++ {
				sb.WriteString(strings.Repeat(HLine, t.columnWidths[i]+2))
				if i < len(t.columnWidths)-1 {
					// Use BottomT to "absorb" the column separator
					sb.WriteString(BottomT)
				}
			}
			sb.WriteString(RightT)
			sb.WriteString("\n")

			// Render notes section
			sb.WriteString(VLine)
			sb.WriteString(t.formatCellContent("", 0)) // Empty first column
			sb.WriteString(VLine)
			notesHeader := "[ RF Advisory ]"
			paddingSpace := mergedWidth - utf8.RuneCountInString(notesHeader)
			sb.WriteString(notesHeader)
			sb.WriteString(strings.Repeat(" ", paddingSpace))
			sb.WriteString(VLine)
			sb.WriteString("\n")

			// Render bullet points
			for _, bp := range bulletPoints {
				bp = strings.TrimSpace(bp)
				if bp != "" {
					bulletLine := "   • " + bp
					wrappedLines := t.smartSplitByWords(bulletLine, mergedWidth-4)

					for _, wrappedLine := range wrappedLines {
						sb.WriteString(VLine)
						sb.WriteString(t.formatCellContent("", 0)) // Empty first column
						sb.WriteString(VLine)

						paddingSpace := mergedWidth - utf8.RuneCountInString(wrappedLine)
						sb.WriteString(wrappedLine)
						sb.WriteString(strings.Repeat(" ", paddingSpace))
						sb.WriteString(VLine)
						sb.WriteString("\n")
					}
				}
			}

			// Check if last row and if next row has description
			isLastRow := ri == len(t.Rows)-1
			nextHasDesc := false
			if !isLastRow {
				if _, ok := t.Descriptions[ri+1]; ok {
					nextHasDesc = true
				}
			}

			// Render appropriate border after description
			if isLastRow {
				// Last row with description - properly render bottom border
				sb.WriteString(BottomLeft)
				sb.WriteString(strings.Repeat(HLine, t.columnWidths[0]+2))
				sb.WriteString(BottomT)
				sb.WriteString(strings.Repeat(HLine, mergedWidth))
				sb.WriteString(BottomRight)
				sb.WriteString("\n")
			} else if nextHasDesc {
				sb.WriteString(t.renderDescToDataBorder())
				// Next row has description - special handling
				// sb.WriteString(VLine)
				// sb.WriteString(t.formatCellContent("", 0))
				// sb.WriteString(LeftT)

				// // Use HLine for the merged part, TopT for where columns will split
				// for i := 1; i < len(t.columnWidths); i++ {
				// 	sb.WriteString(strings.Repeat(HLine, t.columnWidths[i]+2))
				// 	if i < len(t.columnWidths)-1 {
				// 		sb.WriteString(TopT) // Columns will split at these positions
				// 	}
				// }
				// sb.WriteString(RightT)
				// sb.WriteString("\n")
			} else {
				// Next row is normal data - use Cross for normal separator
				sb.WriteString(LeftT)
				sb.WriteString(strings.Repeat(HLine, t.columnWidths[0]+2))
				sb.WriteString(Cross)

				for i := 1; i < len(t.columnWidths); i++ {
					sb.WriteString(strings.Repeat(HLine, t.columnWidths[i]+2))
					if i < len(t.columnWidths)-1 {
						sb.WriteString(Cross)
					}
				}
				sb.WriteString(RightT)
				sb.WriteString("\n")
			}
		} else {
			// No description - add separator if not last row
			if ri < len(t.Rows)-1 {
				sb.WriteString(t.renderMiddleBorder())
			}
		}
	}

	// Bottom border (only if last row has no description)
	if _, ok := t.Descriptions[len(t.Rows)-1]; !ok {
		sb.WriteString(t.renderBottomBorder())
	}

	return sb.String()
}

// Render a line of description text
func (t *Table) renderDescriptionLine(sb *strings.Builder, text string, mergedWidth int) {
	sb.WriteString(VLine)
	sb.WriteString(t.formatCellContent("", 0)) // Empty first column
	sb.WriteString(VLine)

	paddingSpace := mergedWidth - utf8.RuneCountInString(text) - 2
	sb.WriteString(" ")
	sb.WriteString(text)
	sb.WriteString(strings.Repeat(" ", paddingSpace))
	sb.WriteString(" ")
	sb.WriteString(VLine)
	sb.WriteString("\n")
}

// Render border from description to normal data row
func (t *Table) renderDescToDataBorder() string {

	var sb strings.Builder
	sb.WriteString(LeftT)
	sb.WriteString(strings.Repeat(HLine, t.columnWidths[0]+2))
	sb.WriteString(Cross)
	for i := 1; i < len(t.columnWidths); i++ {
		sb.WriteString(strings.Repeat(HLine, t.columnWidths[i]+2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(TopT)
		}
	}
	sb.WriteString(RightT + "\n")
	return sb.String()
}

// Helper method to render border after a description (going back to normal data row)
func (t *Table) renderDataToDataBorderAfterDesc() string {
	var sb strings.Builder
	sb.WriteString(LeftT)

	// First column
	sb.WriteString(strings.Repeat(HLine, t.columnWidths[0]+2))
	sb.WriteString(Cross)

	// Rest of columns normally
	for i := 1; i < len(t.columnWidths); i++ {
		sb.WriteString(strings.Repeat(HLine, t.columnWidths[i]+2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(Cross)
		}
	}

	sb.WriteString(RightT)
	sb.WriteString("\n")
	return sb.String()
}
