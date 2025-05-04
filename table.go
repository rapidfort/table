// Package table provides a utility for rendering ASCII tables with box-drawing characters
package table

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	// ANSI codes
	DimStyleStart  = "\x1b[2m"
	DimStyleEnd    = "\x1b[0m"
	BoldStyleStart = "\x1b[1m" // Bold style for highlighting
	BoldStyleEnd   = "\x1b[0m"

	// Box drawing characters
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

	// Cell sizing constants
	minTerminalWidth = 80
	maxColumnWidth   = 50
)

// Table represents a table with borders and alignment control
type Table struct {
	Headers            []string
	Rows               [][]string
	Descriptions       map[int]string
	columnWidths       []int
	alignments         []string // "left", "right", "center" for each column
	consoleWidth       int      // Maximum width of the console
	fillWidth          bool
	maxWidths          map[int]int // Maximum width for specific columns
	dimBorder          bool        // New field
	highlightHeaders   bool        // Always highlight headers
	highlightedHeaders []int       // Indices of headers to highlight
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

// GetTables returns the tables in the group (public method for testing)
func (g *TableGroup) GetTables() []*Table {
	return g.tables
}

func (t *Table) SetDimBorder(enabled bool) {
	t.dimBorder = enabled
}

// SetHeaderHighlighting enables/disables header highlighting
func (t *Table) SetHeaderHighlighting(enabled bool) {
	t.highlightHeaders = enabled
}

// SetHighlightedHeaders sets which headers should be highlighted
func (t *Table) SetHighlightedHeaders(indices []int) {
	t.highlightedHeaders = indices
}

// AddHighlightedHeader adds a header to the highlighted list
func (t *Table) AddHighlightedHeader(index int) {
	if index >= 0 && index < len(t.Headers) {
		t.highlightedHeaders = append(t.highlightedHeaders, index)
	}
}

// ClearHighlightedHeaders removes all header highlights
func (t *Table) ClearHighlightedHeaders() {
	t.highlightedHeaders = nil
}

// isHighlightedHeader checks if a header at given index should be highlighted
func (t *Table) isHighlightedHeader(index int) bool {
	if t.highlightHeaders {
		return true // Highlight all headers if this flag is set
	}
	for _, i := range t.highlightedHeaders {
		if i == index {
			return true
		}
	}
	return false
}

// getStyledChar returns a border character with optional dim styling
func (t *Table) getStyledChar(char string) string {
	if t.dimBorder {
		return DimStyleStart + char + DimStyleEnd
	}
	return char
}

// getStyledHLine returns a horizontal line string with optional dim styling
func (t *Table) getStyledHLine(width int) string {
	if t.dimBorder {
		return DimStyleStart + strings.Repeat(HLine, width) + DimStyleEnd
	}
	return strings.Repeat(HLine, width)
}

// getHighlightedText returns text with bold styling if it should be highlighted
func (t *Table) getHighlightedText(text string, headerIndex int) string {
	if t.isHighlightedHeader(headerIndex) {
		return BoldStyleStart + text + BoldStyleEnd
	}
	return text
}

// stripANSI removes ANSI codes from a string for length calculation
func stripANSI(str string) string {
	// This is a simple implementation that removes common ANSI codes
	str = strings.ReplaceAll(str, BoldStyleStart, "")
	str = strings.ReplaceAll(str, BoldStyleEnd, "")
	str = strings.ReplaceAll(str, DimStyleStart, "")
	str = strings.ReplaceAll(str, DimStyleEnd, "")
	return str
}

// Add adds a table to the group
func (g *TableGroup) Add(table *Table) {
	g.tables = append(g.tables, table)
	table.group = g
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

// formatCellContent formats a cell's content with alignment and padding
func (t *Table) formatCellContent(content string, colIndex int) string {
	w := t.columnWidths[colIndex]

	// Strip ANSI codes for length calculation
	strippedContent := stripANSI(content)
	contentLength := utf8.RuneCountInString(strippedContent)

	switch t.alignments[colIndex] {
	case "right":
		padding := w - contentLength
		if padding < 0 {
			padding = 0
		}
		return fmt.Sprintf(" %s%s ", strings.Repeat(" ", padding), content)
	case "center":
		totalPad := w - contentLength
		if totalPad < 0 {
			totalPad = 0
		}
		left := totalPad / 2
		return fmt.Sprintf(" %s%s%s ",
			strings.Repeat(" ", left), content,
			strings.Repeat(" ", totalPad-left))
	default:
		padding := w - contentLength
		if padding < 0 {
			padding = 0
		}
		return fmt.Sprintf(" %s%s ", content, strings.Repeat(" ", padding))
	}
}

// calculateInitialColumnWidths computes the initial width for each column
func (t *Table) calculateInitialColumnWidths() {
	// Ensure columnWidths is properly initialized
	if len(t.columnWidths) != len(t.Headers) {
		t.columnWidths = make([]int, len(t.Headers))
	}

	// Reset all column widths to 0 to start fresh
	for i := range t.columnWidths {
		t.columnWidths[i] = 0
	}

	// Calculate minimum width needed for headers
	for i, header := range t.Headers {
		if l := utf8.RuneCountInString(header); l > t.columnWidths[i] {
			t.columnWidths[i] = l
		}
	}

	// Calculate minimum width needed for each cell
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
		// Apply global max column width
		if width > maxColumnWidth {
			t.columnWidths[i] = maxColumnWidth
		}
		// Apply column-specific max width
		if maxWidth, exists := t.maxWidths[i]; exists && t.columnWidths[i] > maxWidth {
			t.columnWidths[i] = maxWidth
		}
	}
}

// adjustColumnWidthsToFit adjusts column widths to fit the console
func (t *Table) adjustColumnWidthsToFit() {
	// Calculate current table width including borders and padding
	total := 1 // Left border
	for _, w := range t.columnWidths {
		total += w + 2 + 1 // Content + padding + separator
	}

	// If table exceeds terminal width, shrink columns
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

	// REMOVED: Do not expand columns to fill width when not needed
	// The table will use only the minimum required width for each column
}

func (t *Table) smartSplitCellContent(content string, colIndex int) []string {
	maxWidth := t.columnWidths[colIndex]
	if utf8.RuneCountInString(content) <= maxWidth {
		return []string{content}
	}

	// Try splitting by comma first
	if strings.Contains(content, ",") {
		return t.splitCommaSeparatedList(content, maxWidth)
	}
	// Then try slash
	if strings.Contains(content, "/") || strings.Contains(content, ".") {
		return t.splitLongString(content, maxWidth)
	}
	// Fall back to word splitting
	return t.splitByWords(content, maxWidth)
}

func (t *Table) splitLongString(content string, maxWidth int) []string {
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
			// Handle case where single word is too long
			for utf8.RuneCountInString(line) > maxWidth {
				part := line[:maxWidth]
				res = append(res, part)
				line = line[maxWidth:]
			}
		}
	}
	if line != "" {
		res = append(res, line)
	}
	return res
}

func (t *Table) renderTopBorder() string {
	var sb strings.Builder
	sb.WriteString(t.getStyledChar(TopLeft))
	for i, w := range t.columnWidths {
		sb.WriteString(t.getStyledHLine(w + 2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(t.getStyledChar(TopT))
		}
	}
	sb.WriteString(t.getStyledChar(TopRight) + "\n")
	return sb.String()
}

func (t *Table) renderMiddleBorder() string {
	var sb strings.Builder
	sb.WriteString(t.getStyledChar(LeftT))
	for i, w := range t.columnWidths {
		sb.WriteString(t.getStyledHLine(w + 2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(t.getStyledChar(Cross))
		}
	}
	sb.WriteString(t.getStyledChar(RightT) + "\n")
	return sb.String()
}

func (t *Table) renderBottomBorder() string {
	var sb strings.Builder
	sb.WriteString(t.getStyledChar(BottomLeft))
	for i, w := range t.columnWidths {
		sb.WriteString(t.getStyledHLine(w + 2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(t.getStyledChar(BottomT))
		}
	}
	sb.WriteString(t.getStyledChar(BottomRight) + "\n")
	return sb.String()
}

// detectTerminalWidth gets the current terminal width or returns a default
func detectTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < minTerminalWidth {
		// Use minimum terminal width if we can't detect or if detected width is too small
		return minTerminalWidth
	}
	return width
}

// RapidFortTable creates a new Table with the given headers
func RapidFortTable(headers []string) *Table {
	// Auto-detect terminal width
	termWidth := detectTerminalWidth()

	table := &Table{
		Headers:            headers,
		Rows:               [][]string{},
		Descriptions:       make(map[int]string),
		columnWidths:       make([]int, len(headers)),
		alignments:         make([]string, len(headers)),
		consoleWidth:       termWidth,
		fillWidth:          false, // Change default to false - don't fill width unnecessarily
		dimBorder:          false,
		maxWidths:          make(map[int]int),
		highlightHeaders:   true,    // Always highlight headers by default
		highlightedHeaders: []int{}, // Initialize the highlighted headers slice
	}

	// Set default left alignment for all columns
	for i := range headers {
		table.alignments[i] = "left"
	}

	return table
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
		sb.WriteString(t.getStyledChar(VLine))
		for ci := range t.Headers {
			txt := ""
			if line < len(headerLines[ci]) {
				txt = headerLines[ci][line]
			}
			// Apply highlighting to header text if specified
			highlightedTxt := t.getHighlightedText(txt, ci)
			sb.WriteString(t.formatCellContent(highlightedTxt, ci))
			sb.WriteString(t.getStyledChar(VLine))
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
			sb.WriteString(t.getStyledChar(VLine))
			for ci := range row {
				txt := ""
				if line < len(rowLines[ci]) {
					txt = rowLines[ci][line]
				}
				sb.WriteString(t.formatCellContent(txt, ci))
				sb.WriteString(t.getStyledChar(VLine))
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
			sb.WriteString(t.getStyledChar(VLine))
			sb.WriteString(t.formatCellContent("", 0))
			sb.WriteString(t.getStyledChar(LeftT))

			// Draw border for all merged columns
			for i := 1; i < len(t.columnWidths); i++ {
				sb.WriteString(t.getStyledHLine(t.columnWidths[i] + 2))
				if i < len(t.columnWidths)-1 {
					// Use BottomT to "absorb" the column separator
					sb.WriteString(t.getStyledChar(BottomT))
				}
			}
			sb.WriteString(t.getStyledChar(RightT))
			sb.WriteString("\n")

			// Render notes section
			sb.WriteString(t.getStyledChar(VLine))
			sb.WriteString(t.formatCellContent("", 0)) // Empty first column
			sb.WriteString(t.getStyledChar(VLine))
			notesHeader := "[ RF Advisory ]"
			paddingSpace := mergedWidth - utf8.RuneCountInString(notesHeader)
			sb.WriteString(notesHeader)
			sb.WriteString(strings.Repeat(" ", paddingSpace))
			sb.WriteString(t.getStyledChar(VLine))
			sb.WriteString("\n")

			// Render bullet points
			for _, bp := range bulletPoints {
				bp = strings.TrimSpace(bp)
				if bp != "" {
					bulletLine := "   • " + bp
					wrappedLines := t.smartSplitByWords(bulletLine, mergedWidth-4)

					for _, wrappedLine := range wrappedLines {
						sb.WriteString(t.getStyledChar(VLine))
						sb.WriteString(t.formatCellContent("", 0)) // Empty first column
						sb.WriteString(t.getStyledChar(VLine))

						paddingSpace := mergedWidth - utf8.RuneCountInString(wrappedLine)
						sb.WriteString(wrappedLine)
						sb.WriteString(strings.Repeat(" ", paddingSpace))
						sb.WriteString(t.getStyledChar(VLine))
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
				sb.WriteString(t.getStyledChar(BottomLeft))
				sb.WriteString(t.getStyledHLine(t.columnWidths[0] + 2))
				sb.WriteString(t.getStyledChar(BottomT))
				sb.WriteString(t.getStyledHLine(mergedWidth))
				sb.WriteString(t.getStyledChar(BottomRight))
				sb.WriteString("\n")
			} else if nextHasDesc {
				sb.WriteString(t.renderDescToDataBorder())
			} else {
				// Next row is normal data - use Cross for normal separator
				sb.WriteString(t.getStyledChar(LeftT))
				sb.WriteString(t.getStyledHLine(t.columnWidths[0] + 2))
				sb.WriteString(t.getStyledChar(Cross))

				for i := 1; i < len(t.columnWidths); i++ {
					sb.WriteString(t.getStyledHLine(t.columnWidths[i] + 2))
					if i < len(t.columnWidths)-1 {
						sb.WriteString(t.getStyledChar(TopT))
					}
				}
				sb.WriteString(t.getStyledChar(RightT))
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

// SyncColumnWidths ensures all tables in the group have consistent column widths
func (g *TableGroup) SyncColumnWidths() {
	if len(g.tables) == 0 {
		return
	}

	// Initialize with the first table's column count
	firstTable := g.tables[0]
	colCount := len(firstTable.Headers)

	// Ensure first table has columnWidths initialized
	if len(firstTable.columnWidths) != colCount {
		firstTable.columnWidths = make([]int, colCount)
	}

	// Initialize the group's column widths
	g.columnWidths = make([]int, colCount)

	// Find the maximum width for each column across all tables
	for _, table := range g.tables {
		// Ensure each table has columnWidths initialized
		if len(table.columnWidths) != colCount {
			table.columnWidths = make([]int, colCount)
		}

		table.calculateInitialColumnWidths()

		for i := 0; i < colCount && i < len(table.columnWidths); i++ {
			// Apply maximum column width constraint
			width := table.columnWidths[i]
			if width > maxColumnWidth {
				width = maxColumnWidth
			}

			if width > g.columnWidths[i] {
				// Check if this column has a max width constraint
				if maxWidth, exists := table.maxWidths[i]; exists && width > maxWidth {
					g.columnWidths[i] = maxWidth
				} else {
					g.columnWidths[i] = width
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

// Render border from description to normal data row
func (t *Table) renderDescToDataBorder() string {
	var sb strings.Builder
	sb.WriteString(t.getStyledChar(LeftT))
	sb.WriteString(t.getStyledHLine(t.columnWidths[0] + 2))
	sb.WriteString(t.getStyledChar(Cross))
	for i := 1; i < len(t.columnWidths); i++ {
		sb.WriteString(t.getStyledHLine(t.columnWidths[i] + 2))
		if i < len(t.columnWidths)-1 {
			sb.WriteString(t.getStyledChar(TopT))
		}
	}
	sb.WriteString(t.getStyledChar(RightT) + "\n")
	return sb.String()
}
