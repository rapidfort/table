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

// New creates a new Table with the given headers
func RapidFortTable(headers []string) *Table {
	table := &Table{
		Headers:      headers,
		Rows:         [][]string{},
		Descriptions: make(map[int]string),
		columnWidths: make([]int, len(headers)),
		alignments:   make([]string, len(headers)),
		consoleWidth: detectWidth(), // Auto-detect console width by default
		maxWidths:    make(map[int]int),
	}
	for i := range headers {
		table.alignments[i] = "left"
	}
	return table
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

	// Draw a continuous line across the entire width with BottomT at column positions
	for i, w := range t.columnWidths {
		sb.WriteString(strings.Repeat(HLine, w+2))

		// Add separator if not the last column
		if i < len(t.columnWidths)-1 {
			sb.WriteString(BottomT) // Use BottomT (┴) instead of Cross (┼)
		}
	}

	sb.WriteString(RightT + "\n")
	return sb.String()
}

// renderDescToDataBorder renders the border between a description and the next data row
// Uses TopT (┬) for column separators instead of Cross (┼)
func (t *Table) renderDescToDataBorder() string {
	var sb strings.Builder
	sb.WriteString(LeftT)

	// Draw a continuous line across the entire width with TopT at column positions
	for i, w := range t.columnWidths {
		sb.WriteString(strings.Repeat(HLine, w+2))

		// Add separator if not the last column
		if i < len(t.columnWidths)-1 {
			sb.WriteString(TopT) // Use TopT (┬) instead of Cross (┼)
		}
	}

	sb.WriteString(RightT + "\n")
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
		prefix := " [ Notes ] "
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

// Render the table to a string
func (t *Table) Render() string {
	// If part of a group, the group will handle column width calculations
	if t.group == nil {
		t.calculateInitialColumnWidths()
		t.adjustColumnWidthsToFit()
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

		// If there's a description for this row
		if desc, ok := t.Descriptions[ri]; ok {
			// Add special separator for data-to-description transition
			sb.WriteString(t.renderDataToDescBorder())

			// Wrap and split the description text
			bulletPoints := strings.Split(desc, "\n")
			if len(bulletPoints) == 1 && strings.Contains(desc, ",") {
				// If no newlines but has commas, split by commas
				bulletPoints = strings.Split(desc, ",")
				for i, bp := range bulletPoints {
					bulletPoints[i] = strings.TrimSpace(bp)
				}
			}

			// Add description header
			sb.WriteString(t.renderDescriptionRow("", true))

			// Add each bullet point
			for _, bp := range bulletPoints {
				bp = strings.TrimSpace(bp)
				if bp != "" {
					// Wrap long bullet points
					wrappedLines := t.wrapDescriptionText(bp)
					for _, line := range wrappedLines {
						sb.WriteString(t.renderDescriptionRow(line, false))
					}
				}
			}

			// If this is the last row, add special bottom border for description
			if ri == len(t.Rows)-1 {
				sb.WriteString(t.renderBottomDescriptionBorder())
				return sb.String() // Return early since we've added the bottom border
			}

			// If not the last row, add special separator for description-to-data transition
			sb.WriteString(t.renderDescToDataBorder())
			continue
		}

		// If this is not the last row and it doesn't have a description, add normal middle border
		if ri < len(t.Rows)-1 {
			sb.WriteString(t.renderMiddleBorder())
		}
	}

	// Bottom border (only reached if the last row doesn't have a description)
	sb.WriteString(t.renderBottomBorder())

	return sb.String()
}

// Helper function to check if a row has a description
func (t *Table) hasDescription(rowIndex int) bool {
	_, hasDesc := t.Descriptions[rowIndex]
	return hasDesc
}

// Helper function to check if the next row exists and has a description
func (t *Table) nextRowHasDescription(rowIndex int) bool {
	if rowIndex+1 >= len(t.Rows) {
		return false
	}
	_, hasDesc := t.Descriptions[rowIndex+1]
	return hasDesc
}

// detectWidth returns the current terminal width (or 80 if it can't detect one)
func detectWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	return 80
}
