package table

import (
	"fmt"
	"strings"
	"testing"
)

// Basic functionality tests
func TestTableCreation(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	tbl := RapidFortTable(headers)

	if len(tbl.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(tbl.Headers))
	}

	if tbl.Headers[0] != "Name" || tbl.Headers[1] != "Age" || tbl.Headers[2] != "City" {
		t.Errorf("Headers do not match expected values")
	}
}

func TestTableAddRow(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	tbl := RapidFortTable(headers)

	tbl.AddRow([]string{"Alice", "30", "New York"})
	tbl.AddRow([]string{"Bob", "25"})

	if len(tbl.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(tbl.Rows))
	}

	// Check that missing columns are filled with empty string
	if len(tbl.Rows[1]) != 3 {
		t.Errorf("Expected row to be padded to 3 columns, got %d", len(tbl.Rows[1]))
	}

	if tbl.Rows[1][2] != "" {
		t.Errorf("Expected empty string for missing cell, got %s", tbl.Rows[1][2])
	}
}

// Description functionality tests
func TestTableAddDescription(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	tbl := RapidFortTable(headers)

	tbl.AddRow([]string{"Alice", "30", "New York"})
	tbl.AddDescription(0, "Sample description")

	desc, ok := tbl.Descriptions[0]
	if !ok {
		t.Errorf("Description not added correctly")
	}

	if desc != "Sample description" {
		t.Errorf("Description does not match expected value, got: %s", desc)
	}
}

func TestMultipleDescriptions(t *testing.T) {
	headers := []string{"Name", "Age"}
	tbl := RapidFortTable(headers)

	// Force console width to avoid negative padding calculation
	tbl.SetConsoleWidth(80)

	tbl.AddRow([]string{"Alice", "30"})
	tbl.AddRow([]string{"Bob", "25"})
	tbl.AddRow([]string{"Charlie", "35"})

	tbl.AddDescription(0, "First description")
	tbl.AddDescription(2, "Last description")

	rendered := tbl.Render()

	if !strings.Contains(rendered, "First description") {
		t.Errorf("First description not rendered correctly")
	}

	if !strings.Contains(rendered, "Last description") {
		t.Errorf("Last description not rendered correctly")
	}
}

// Alignment tests
func TestTableAlignment(t *testing.T) {
	headers := []string{"Left", "Right", "Center"}
	tbl := RapidFortTable(headers)

	// Force console width
	tbl.SetConsoleWidth(80)

	tbl.SetAlignment(0, "left")
	tbl.SetAlignment(1, "right")
	tbl.SetAlignment(2, "center")

	if tbl.alignments[0] != "left" {
		t.Errorf("Left alignment not set correctly")
	}

	if tbl.alignments[1] != "right" {
		t.Errorf("Right alignment not set correctly")
	}

	if tbl.alignments[2] != "center" {
		t.Errorf("Center alignment not set correctly")
	}

	tbl.AddRow([]string{"AAA", "BBB", "CCC"})
	rendered := tbl.Render()

	// Basic check to ensure rendering doesn't fail
	if !strings.Contains(rendered, "AAA") || !strings.Contains(rendered, "BBB") || !strings.Contains(rendered, "CCC") {
		t.Errorf("Table with alignments not rendered correctly")
	}
}

// Width constraints tests
func TestTableMaxWidth(t *testing.T) {
	headers := []string{"Column1", "Column2"}
	tbl := RapidFortTable(headers)

	// Force console width
	tbl.SetConsoleWidth(80)

	// Set max width constraint
	tbl.SetMaxWidth(0, 5)

	// Add row with content exceeding max width
	tbl.AddRow([]string{"ThisIsLong", "Normal"})

	// Calculate column widths
	tbl.calculateInitialColumnWidths()

	// Verify the constraint is applied
	if tbl.columnWidths[0] > 5 {
		t.Errorf("Max width constraint not applied, got width: %d", tbl.columnWidths[0])
	}
}

func TestFillWidth(t *testing.T) {
	headers := []string{"A", "B"}
	tbl := RapidFortTable(headers)

	// Force console width for testing
	tbl.SetConsoleWidth(50)
	tbl.SetFillWidth(true)

	tbl.AddRow([]string{"Short", "Content"})

	// Calculate and adjust column widths
	tbl.calculateInitialColumnWidths()
	tbl.adjustColumnWidthsToFit()

	// Calculate total used width
	totalWidth := 1 // Starting border
	for _, w := range tbl.columnWidths {
		totalWidth += w + 2 + 1 // column width + padding + separator
	}

	// Should be close to the console width (might not be exact due to constraints)
	if totalWidth < 40 {
		t.Errorf("Fill width not working correctly, table width: %d, expected near 50", totalWidth)
	}
}

// Cell content wrapping tests
func TestSmartSplitCellContent(t *testing.T) {
	headers := []string{"Package", "Description"}
	tbl := RapidFortTable(headers)

	// Set console width
	tbl.SetConsoleWidth(80)

	// Set column width and maxWidth
	tbl.SetMaxWidth(0, 10)
	tbl.columnWidths = []int{10, 15} // Explicitly set column widths for testing

	// Test package name splitting
	lines := tbl.smartSplitCellContent("github.com/user/long-package-name", 0)
	if len(lines) <= 1 {
		t.Errorf("Expected package path to be split into multiple lines, got %d lines", len(lines))
	}

	// Test comma list splitting
	lines = tbl.smartSplitCellContent("item1, item2, item3, a-very-long-item", 1)
	if len(lines) <= 1 {
		t.Errorf("Expected comma list to be split into multiple lines, got %d lines", len(lines))
	}

	// Test word splitting
	lines = tbl.smartSplitCellContent("This is a long sentence that should be wrapped", 1)
	if len(lines) <= 1 {
		t.Errorf("Expected long text to be split into multiple lines, got %d lines", len(lines))
	}
}

// Table group tests
func TestTableGroup(t *testing.T) {
	group := NewGroup()

	headers1 := []string{"Col1", "Col2"}
	headers2 := []string{"Col1", "Col2"}

	tbl1 := RapidFortTable(headers1)
	tbl2 := RapidFortTable(headers2)

	// Force console width
	tbl1.SetConsoleWidth(80)
	tbl2.SetConsoleWidth(80)

	// Add different width content to each table
	tbl1.AddRow([]string{"Short", "Data"})
	tbl2.AddRow([]string{"ThisIsVeryLongContent", "AlsoLong"})

	group.Add(tbl1)
	group.Add(tbl2)

	// Sync the column widths
	group.SyncColumnWidths()

	// Both tables should now have the same column widths
	if tbl1.columnWidths[0] != tbl2.columnWidths[0] {
		t.Errorf("Column widths not synchronized, got %d and %d",
			tbl1.columnWidths[0], tbl2.columnWidths[0])
	}

	// Make sure tables can be rendered without error
	_ = tbl1.Render()
	_ = tbl2.Render()
}

// Edge case tests
func TestEmptyTable(t *testing.T) {
	headers := []string{}
	tbl := RapidFortTable(headers)
	tbl.SetConsoleWidth(80)

	// Should not panic
	rendered := tbl.Render()
	if rendered == "" {
		t.Errorf("Empty table should render something, not an empty string")
	}
}

func TestTableWithEmptyCells(t *testing.T) {
	headers := []string{"A", "B", "C"}
	tbl := RapidFortTable(headers)
	tbl.SetConsoleWidth(80)

	tbl.AddRow([]string{"", "", ""})
	tbl.AddRow([]string{"Content", "", ""})

	// Should handle empty cells gracefully
	rendered := tbl.Render()
	if !strings.Contains(rendered, "Content") {
		t.Errorf("Table with empty cells not rendered correctly")
	}
}

// Full render test
func TestFullRender(t *testing.T) {
	headers := []string{"ID", "Name", "Status"}
	tbl := RapidFortTable(headers)
	tbl.SetConsoleWidth(80)

	tbl.AddRow([]string{"1", "Project Alpha", "Active"})
	tbl.AddRow([]string{"2", "Project Beta", "On Hold"})
	tbl.AddDescription(1, "This project is waiting for client approval.")

	rendered := tbl.Render()

	expectedElements := []string{
		"┌", "┐", // Top corners
		"└", "┘", // Bottom corners
		"│", "─", // Lines
		"┼", "┬", "┴", // Junctions
		"ID", "Name", "Status", // Headers
		"Project Alpha", "Active", // Row 1 data
		"Project Beta", "On Hold", // Row 2 data
		"[ Notes ]",                   // Description marker
		"waiting for client approval", // Description content
	}

	for _, element := range expectedElements {
		if !strings.Contains(rendered, element) {
			t.Errorf("Expected element %s not found in rendered table", element)
		}
	}
}

// Benchmark for performance-critical operations
func BenchmarkTableRender(b *testing.B) {
	headers := []string{"Col1", "Col2", "Col3", "Col4", "Col5"}
	tbl := RapidFortTable(headers)
	tbl.SetConsoleWidth(120)

	// Add 100 rows with varying content
	for i := 0; i < 100; i++ {
		tbl.AddRow([]string{
			fmt.Sprintf("Cell%d", i),
			"This is some longer content",
			fmt.Sprintf("Row%d/Item", i),
			"github.com/user/package/subpackage",
			"item1, item2, item3, item4",
		})

		// Add descriptions to some rows
		if i%5 == 0 {
			tbl.AddDescription(i, "This is a detailed description for row "+fmt.Sprintf("%d", i)+
				" that explains something important about this entry in the table.")
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tbl.Render()
	}
}
