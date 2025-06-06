// Package table provides a utility for rendering ASCII tables with box-drawing characters
package table

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	// ANSI codes
	//DimStyleStart = "\x1b[38;5;242m"
	// faint + darker gray for very subdued borders
	DimStyleStart = "\x1b[2m\x1b[38;5;240m"

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
	Descriptions       map[int][]string // row index -> description
	DescriptionTitles  map[int][]string // row index -> title (optional)
	columnWidths       []int
	alignments         []string // "left", "right", "center" for each column
	consoleWidth       int      // Maximum width of the console
	fillWidth          bool
	maxWidths          map[int]int // Maximum width for specific columns
	dimBorder          bool        // New field
	supportANSI        bool        // Support for ANSI codes
	borderless         bool        // Flag to disable borders
	highlightHeaders   bool        // Always highlight headers
	highlightedHeaders []int       // Indices of headers to highlight
	rowCountEnabled    bool        // Flag to enable row count
	// Reference to the table group this table belongs to (if any)

	group *TableGroup
}

func (t *Table) EnableRowCount(enabled bool) *Table {
	t.rowCountEnabled = enabled
	return t
}

// ansiRegexp matches any CSI sequence (e.g. "\x1b[31m", "\x1b[0K", etc.)
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// stripANSI removes ALL ANSI escape sequences from s.
func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// extractWrappingANSI splits s into leading CSI prefix, trailing CSI suffix,
func extractWrappingANSI(s string) (prefix, suffix, core string) {
	// 1) Peel off all leading CSI sequences
	coreStart := 0
	for {
		// Find next ANSI match at the very start of the remaining string
		loc := ansiRegexp.FindStringIndex(s[coreStart:])
		if loc == nil || loc[0] != 0 {
			break
		}
		// Extend prefix to include that entire match
		matchEnd := coreStart + loc[1]
		prefix += s[coreStart:matchEnd]
		coreStart = matchEnd
		if coreStart >= len(s) {
			// we've consumed the whole string
			return prefix, suffix, ""
		}
	}

	// 2) Peel off all trailing CSI sequences
	coreEnd := len(s)
	for {
		// Find all matches in the substring up to coreEnd
		locs := ansiRegexp.FindAllStringIndex(s[:coreEnd], -1)
		if len(locs) == 0 {
			break
		}
		last := locs[len(locs)-1]
		// If the last match ends exactly at coreEnd, it's truly a suffix
		if last[1] == coreEnd {
			suffix = s[last[0]:coreEnd] + suffix
			coreEnd = last[0]
			if coreEnd <= coreStart {
				// nothing left in core
				return prefix, suffix, ""
			}
		} else {
			break
		}
	}

	// 3) What's in between is the printable core
	core = s[coreStart:coreEnd]
	return
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

// SetBorderless enables/disables drawing of any box‐drawing characters.
func (t *Table) SetBorderless(on bool) {
	t.borderless = on
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
	if t.borderless {
		return " "
	}
	if t.dimBorder && t.supportANSI {
		return DimStyleStart + char + DimStyleEnd
	}
	return char
}

// getStyledHLine returns a horizontal line string with optional dim styling
func (t *Table) getStyledHLine(width int) string {
	if t.borderless {
		return strings.Repeat(" ", width)
	}
	if t.dimBorder && t.supportANSI {
		return DimStyleStart + strings.Repeat(HLine, width) + DimStyleEnd
	}
	return strings.Repeat(HLine, width)
}

// getHighlightedText returns text with bold styling if it should be highlighted
func (t *Table) getHighlightedText(text string, headerIndex int) string {
	if !t.supportANSI {
		return text
	}

	if t.isHighlightedHeader(headerIndex) {
		return BoldStyleStart + text + BoldStyleEnd
	}
	return text
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
		if _, ok := t.Descriptions[rowIndex]; !ok {
			t.Descriptions[rowIndex] = []string{}
			t.DescriptionTitles[rowIndex] = []string{}
		}
		t.Descriptions[rowIndex] = append(t.Descriptions[rowIndex], description)
		t.DescriptionTitles[rowIndex] = append(t.DescriptionTitles[rowIndex], "")
	}
}

// AddDescriptionWithTitle adds a description with a title for a specific row
func (t *Table) AddDescriptionWithTitle(rowIndex int, title string, description string) {
	if rowIndex >= 0 && rowIndex < len(t.Rows) {
		if _, ok := t.Descriptions[rowIndex]; !ok {
			t.Descriptions[rowIndex] = []string{}
			t.DescriptionTitles[rowIndex] = []string{}
		}
		t.Descriptions[rowIndex] = append(t.Descriptions[rowIndex], description)
		t.DescriptionTitles[rowIndex] = append(t.DescriptionTitles[rowIndex], title)
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
		// headers have no ANSI, but let's strip anyway for consistency
		vis := stripANSI(header)
		if l := utf8.RuneCountInString(vis); l > t.columnWidths[i] {
			t.columnWidths[i] = l
		}
	}

	// Calculate minimum width needed for each cell
	for _, row := range t.Rows {
		for i, cell := range row {
			if i >= len(t.columnWidths) {
				continue
			}
			// strip out color codes before measuring
			vis := stripANSI(cell)
			if l := utf8.RuneCountInString(vis); l > t.columnWidths[i] {
				t.columnWidths[i] = l
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

}

// smartSplitCellContent splits a cell, preserving any ANSI prefix/suffix,
// and applies your original: comma-first, slash-second, then word-fallback.
func (t *Table) smartSplitCellContent(content string, colIndex int) []string {
	// 1) Peel off any ANSI wrapper
	prefix, suffix, core := extractWrappingANSI(content)

	// 2) Measure the visible length
	maxW := t.columnWidths[colIndex]
	visible := stripANSI(core)
	if utf8.RuneCountInString(visible) <= maxW {
		// nothing to wrap
		return []string{prefix + core + suffix}
	}

	// 3) Try your original split strategies on the **plain** core,
	//    then re-attach prefix/suffix to each piece.

	var parts []string
	if strings.Contains(core, ",") {
		parts = t.splitCommaSeparatedList(core, maxW)
	} else if strings.Contains(core, "/") || strings.Contains(core, ".") {
		parts = t.splitLongString(core, maxW)
	} else {
		parts = t.splitByWords(core, maxW)
	}

	// 4) Re-attach ANSI to every wrapped line
	out := make([]string, len(parts))
	for i, line := range parts {
		out[i] = prefix + line + suffix
	}
	return out
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

func RapidFortTable(headers []string) *Table {
	// Create a copy of the headers slice to avoid modifying the original
	headersCopy := make([]string, len(headers))
	copy(headersCopy, headers)

	return NewTable(headersCopy)
}

func NewTable(headers []string) *Table {
	// Auto-detect terminal width
	termWidth := detectTerminalWidth()

	table := &Table{
		Headers:            headers,
		Rows:               [][]string{},
		Descriptions:       make(map[int][]string),
		DescriptionTitles:  make(map[int][]string), // Initialize the new field
		columnWidths:       make([]int, len(headers)),
		alignments:         make([]string, len(headers)),
		consoleWidth:       termWidth,
		fillWidth:          false, // Change default to false - don't fill width unnecessarily
		dimBorder:          true,
		supportANSI:        term.IsTerminal(int(os.Stdout.Fd())),
		maxWidths:          make(map[int]int),
		highlightHeaders:   true,    // Always highlight headers by default
		highlightedHeaders: []int{}, // Initialize the highlighted headers slice
		rowCountEnabled:    false,
	}

	if !table.supportANSI {
		table.dimBorder = false
		table.highlightHeaders = false
	}

	// Set default left alignment for all columns
	for i := range headers {
		table.alignments[i] = "left"
	}

	return table
}

// smartSplitByWords splits text by words to fit within maxWidth
func (t *Table) smartSplitByWords(text string, maxWidth int) []string {
	// Strip ANSI for width calculation, but keep original for output
	textVisible := stripANSI(text)

	// If the text already fits, no need to split
	if utf8.RuneCountInString(textVisible) <= maxWidth {
		return []string{text}
	}

	// Extract ANSI prefix/suffix if any
	prefix, suffix, visibleContent := extractWrappingANSI(text)

	// Split the visible content by spaces
	words := strings.Fields(visibleContent)
	if len(words) == 0 {
		return []string{text}
	}

	var result []string
	currentLine := ""
	currentLineVisible := ""

	for _, word := range words {
		// Extract any ANSI codes in this word
		wordPrefix, wordSuffix, wordVisible := extractWrappingANSI(word)

		// Calculate visible length for the test line
		testLineVisible := currentLineVisible
		if testLineVisible != "" {
			testLineVisible += " "
		}
		testLineVisible += wordVisible

		if utf8.RuneCountInString(testLineVisible) <= maxWidth {
			// Word fits on current line
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += wordPrefix + wordVisible + wordSuffix
			currentLineVisible = testLineVisible
		} else {
			// Word doesn't fit, start new line
			if currentLine != "" {
				result = append(result, prefix+currentLine+suffix)
			}

			// If the word itself is too long, split it
			if utf8.RuneCountInString(wordVisible) > maxWidth {
				// Create chunks of the word that fit
				var chunks []string
				remaining := wordVisible
				for utf8.RuneCountInString(remaining) > 0 {
					chunkSize := maxWidth
					if utf8.RuneCountInString(remaining) <= chunkSize {
						chunks = append(chunks, remaining)
						break
					}

					chunk := remaining[:chunkSize]
					chunks = append(chunks, chunk)
					remaining = remaining[chunkSize:]
				}

				// Add chunks as separate lines
				for i, chunk := range chunks {
					if i == 0 {
						result = append(result, prefix+wordPrefix+chunk+wordSuffix+suffix)
					} else {
						result = append(result, prefix+chunk+suffix)
					}
				}

				currentLine = ""
				currentLineVisible = ""
			} else {
				currentLine = wordPrefix + wordVisible + wordSuffix
				currentLineVisible = wordVisible
			}
		}
	}

	// Add any remaining text
	if currentLine != "" {
		result = append(result, prefix+currentLine+suffix)
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

// Function to prepare the table with row counting
func (t *Table) prepareWithRowCount() *Table {
	if !t.rowCountEnabled {
		return t
	}

	// Create a new table with row counts
	newHeaders := append([]string{"#"}, t.Headers...)
	newTable := NewTable(newHeaders)

	// Copy properties from original table
	newTable.Descriptions = t.Descriptions
	newTable.DescriptionTitles = t.DescriptionTitles
	newTable.consoleWidth = t.consoleWidth
	newTable.fillWidth = t.fillWidth
	newTable.dimBorder = t.dimBorder
	newTable.supportANSI = t.supportANSI
	newTable.borderless = t.borderless
	newTable.highlightHeaders = t.highlightHeaders
	newTable.highlightedHeaders = t.highlightedHeaders
	newTable.group = t.group
	newTable.rowCountEnabled = false // Prevent infinite recursion

	// Copy alignments
	newTable.alignments = make([]string, len(newHeaders))
	newTable.alignments[0] = "right" // Right-align row numbers
	for i := 0; i < len(t.alignments); i++ {
		newTable.alignments[i+1] = t.alignments[i]
	}

	// Copy max widths
	for col, width := range t.maxWidths {
		newTable.maxWidths[col+1] = width
	}

	// Add rows with row numbers
	for i, row := range t.Rows {
		rowNum := fmt.Sprintf("%d", i+1)
		newTable.AddRow(append([]string{rowNum}, row...))
	}

	return newTable
}

func (t *Table) Render() string {
	// 1) If stdout isn't a real terminal, drop ALL ANSI and use minimal widths
	if t.rowCountEnabled {
		return t.prepareWithRowCount().Render()
	}

	if !t.supportANSI {
		// Disable ANSI-based decorations
		t.dimBorder = false
		t.highlightHeaders = false

		// Strip ANSI from headers
		for i, h := range t.Headers {
			t.Headers[i] = stripANSI(h)
		}
		// Strip ANSI from every table cell
		for ri, row := range t.Rows {
			for ci, cell := range row {
				t.Rows[ri][ci] = stripANSI(cell)
			}
		}
		// Strip ANSI from every description & title
		for ri, descs := range t.Descriptions {
			for i, desc := range descs {
				t.Descriptions[ri][i] = stripANSI(desc)
			}
			if titles, ok := t.DescriptionTitles[ri]; ok {
				for i, title := range titles {
					t.DescriptionTitles[ri][i] = stripANSI(title)
				}
			}
		}
		// Compute the absolute minimal column widths
		t.calculateInitialColumnWidths()

	} else {
		// 2) ANSI-capable (TTY) mode: use your existing optimal-width logic
		if t.group == nil {
			t.calculateOptimalColumnWidths(t.consoleWidth)
		} else {
			t.adjustColumnWidthsToFit()
		}
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
			highlighted := t.getHighlightedText(txt, ci)
			sb.WriteString(t.formatCellContent(highlighted, ci))
			sb.WriteString(t.getStyledChar(VLine))
		}
		sb.WriteString("\n")
	}

	// Header/Data separator
	sb.WriteString(t.renderMiddleBorder())

	// Rows + Descriptions
	for ri, row := range t.Rows {
		// Data row
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

		// Optional description blocks
		if descs, ok := t.Descriptions[ri]; ok && len(descs) > 0 {
			// Compute merged width of columns 2..n (used by all descriptions)
			mergedWidth := 0
			for i := 1; i < len(t.columnWidths); i++ {
				mergedWidth += t.columnWidths[i] + 2
				if i < len(t.columnWidths)-1 {
					mergedWidth += 1
				}
			}

			// Process each description for this row
			for di, desc := range descs {
				// Top border of description block (first one) or separator between descriptions
				if di == 0 {
					// First description - top border
					sb.WriteString(t.getStyledChar(VLine))
					sb.WriteString(t.formatCellContent("", 0))
					sb.WriteString(t.getStyledChar(LeftT))
					for i := 1; i < len(t.columnWidths); i++ {
						sb.WriteString(t.getStyledHLine(t.columnWidths[i] + 2))
						if i < len(t.columnWidths)-1 {
							sb.WriteString(t.getStyledChar(BottomT))
						}
					}
					sb.WriteString(t.getStyledChar(RightT) + "\n")
				} else {
					// Separator between descriptions
					sb.WriteString(t.getStyledChar(VLine))
					sb.WriteString(t.formatCellContent("", 0))
					sb.WriteString(t.getStyledChar(LeftT))
					// For inter-description separators, we don't want column divisions
					sb.WriteString(t.getStyledHLine(mergedWidth))
					sb.WriteString(t.getStyledChar(RightT) + "\n")
				}

				// Description title (if any)
				if titles, ok := t.DescriptionTitles[ri]; ok && di < len(titles) && titles[di] != "" {
					headerText := " [ " + BoldStyleStart + titles[di] + BoldStyleEnd + " ]"
					pad := mergedWidth - utf8.RuneCountInString(stripANSI(headerText))
					if pad < 0 {
						pad = 0
					}

					sb.WriteString(t.getStyledChar(VLine))
					sb.WriteString(t.formatCellContent("", 0))
					sb.WriteString(t.getStyledChar(VLine))
					sb.WriteString(headerText)
					sb.WriteString(strings.Repeat(" ", pad))
					sb.WriteString(t.getStyledChar(VLine) + "\n")
				}

				// Split into bullet points
				bps := strings.Split(desc, "\n")

				// Bullet lines
				for _, bp := range bps {
					bp = strings.TrimSpace(bp)
					if bp == "" {
						continue
					}
					prefix := " "
					textWidth := mergedWidth - utf8.RuneCountInString(prefix) - 2
					if textWidth < 0 {
						textWidth = 0
					}
					wrapped := t.smartSplitByWords(bp, textWidth)

					for i, wline := range wrapped {
						sb.WriteString(t.getStyledChar(VLine))
						sb.WriteString(t.formatCellContent("", 0))
						sb.WriteString(t.getStyledChar(VLine))

						var disp string
						if i == 0 {
							disp = prefix + wline
						} else {
							indent := strings.Repeat(" ", utf8.RuneCountInString(prefix))
							disp = indent + wline
						}
						pad := mergedWidth - utf8.RuneCountInString(stripANSI(disp))
						if pad < 0 {
							pad = 0
						}
						sb.WriteString(disp)
						sb.WriteString(strings.Repeat(" ", pad))
						sb.WriteString(t.getStyledChar(VLine) + "\n")
					}
				}
			}

			// Decide which border to draw next
			isLast := ri == len(t.Rows)-1
			nextHasDesc := !isLast && (func() bool {
				_, ok2 := t.Descriptions[ri+1]
				return ok2 && len(t.Descriptions[ri+1]) > 0
			})()

			if isLast {
				// Bottom border after last desc
				sb.WriteString(t.getStyledChar(BottomLeft))
				sb.WriteString(t.getStyledHLine(t.columnWidths[0] + 2))
				sb.WriteString(t.getStyledChar(BottomT))
				sb.WriteString(t.getStyledHLine(mergedWidth))
				sb.WriteString(t.getStyledChar(BottomRight) + "\n")
			} else if nextHasDesc {
				sb.WriteString(t.renderDescToDataBorder())
			} else {
				// Normal separator into next data row
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
			}
		} else if ri < len(t.Rows)-1 {
			// No description, normal middle border
			sb.WriteString(t.renderMiddleBorder())
		}
	}

	// Bottom border if last row had no description
	if descs, hasDesc := t.Descriptions[len(t.Rows)-1]; !hasDesc || len(descs) == 0 {
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

// renderDescToDataBorder draws the border after a description before the next data row,
// using ANSI-aware dimmed characters.
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
	sb.WriteString(t.getStyledChar(RightT))
	sb.WriteString("\n")
	return sb.String()
}
