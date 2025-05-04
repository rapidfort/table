package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rapidfort/table" // Replace with your actual module path
)

// Test case generator
type TestCase struct {
	title              string
	description        string
	useGroup           bool
	numTables          int
	dimBorder          bool
	fillWidth          bool
	maxColumnWidths    map[int]int
	hasDescriptions    bool
	highlightedHeaders []int
}

func main() {
    simple()
	// Initialize random source (modern way)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	testCases := []TestCase{
		{
			title:              "Test Case 1: Random Sized Columns",
			description:        "Testing columns with random sizes and content",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 2: Grouped Tables with Sync",
			description:        "Multiple tables with synced column widths",
			useGroup:           true,
			numTables:          3,
			dimBorder:          true,
			fillWidth:          false,
			maxColumnWidths:    map[int]int{1: 15, 3: 20},
			hasDescriptions:    false,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 3: Header Highlighting",
			description:        "Testing header highlighting with different patterns",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: []int{1, 3, 4}, // Highlight these headers
		},
		{
			title:              "Test Case 4: Mixed Highlighting and Styling",
			description:        "Testing headers with highlighting AND dim borders",
			useGroup:           false,
			numTables:          1,
			dimBorder:          true,
			fillWidth:          true,
			maxColumnWidths:    map[int]int{2: 10, 4: 15},
			hasDescriptions:    true,
			highlightedHeaders: []int{0, 2, 5}, // Highlight first, third, and last headers
		},
		{
			title:              "Test Case 5: Grouped Tables with Highlighting",
			description:        "Multiple tables with synced widths and consistent highlighting",
			useGroup:           true,
			numTables:          2,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    false,
			highlightedHeaders: []int{1, 3}, // Highlight name and value columns
		},
		{
			title:              "Test Case 6: All Headers Highlighted",
			description:        "Testing with all headers highlighted",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    false,
			highlightedHeaders: []int{0, 1, 2, 3, 4, 5}, // All headers
		},
		{
			title:              "Test Case 7: Dynamic Header Highlighting",
			description:        "Testing adding and removing highlights dynamically",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    false,
			highlightedHeaders: nil, // Will be set dynamically
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		runTestCase(tc, rng)
		fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
	}
}

func runTestCase(tc TestCase, rng *rand.Rand) {
	fmt.Printf("=== %s ===\n", tc.title)
	fmt.Printf("Description: %s\n\n", tc.description)

	var group *table.TableGroup
	if tc.useGroup {
		group = table.NewGroup()
	}

	switch tc.title {
	case "Test Case 1: Random Sized Columns":
		test1RandomSizedColumns(group, tc, rng)
	case "Test Case 2: Grouped Tables with Sync":
		test2GroupedTablesWithSync(group, tc, rng)
	case "Test Case 3: Header Highlighting":
		test3HeaderHighlighting(group, tc, rng)
	case "Test Case 4: Mixed Highlighting and Styling":
		test4MixedHighlightingAndStyling(group, tc, rng)
	case "Test Case 5: Grouped Tables with Highlighting":
		test5GroupedTablesWithHighlighting(group, tc, rng)
	case "Test Case 6: All Headers Highlighted":
		test6AllHeadersHighlighted(group, tc, rng)
	case "Test Case 7: Dynamic Header Highlighting":
		test7DynamicHeaderHighlighting(group, tc, rng)
	}
}

func test1RandomSizedColumns(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Short", "MediumLength", "A Very Long Header Name", "Tiny", "Another Column", "Last"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)
	tbl.SetFillWidth(tc.fillWidth)
	tbl.SetBorderless(false)

	// Add rows with random sized content
	for i := 0; i < 5; i++ {
		row := []string{
			generateRandomString(rng.Intn(10)+1, rng),
			generateRandomString(rng.Intn(30)+5, rng),
			generateRandomString(rng.Intn(50)+10, rng),
			generateRandomString(rng.Intn(5)+1, rng),
			generateRandomString(rng.Intn(25)+5, rng),
			generateRandomString(rng.Intn(15)+3, rng),
		}
		tbl.AddRow(row)

		if tc.hasDescriptions && rng.Float32() < 0.3 {
			tbl.AddDescription(i, generateRandomDescription(rng))
		}
	}

	fmt.Println(tbl.Render())
}

func test2GroupedTablesWithSync(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"ID", "Name", "Category", "Value", "Status"}

	// First, create all tables without adding to group yet
	tables := make([]*table.Table, tc.numTables)

	for i := 0; i < tc.numTables; i++ {
		tbl := table.RapidFortTable(headers)
		tbl.SetDimBorder(tc.dimBorder)
		tbl.SetFillWidth(tc.fillWidth)

		// Set custom max widths if specified
		if tc.maxColumnWidths != nil {
			for col, width := range tc.maxColumnWidths {
				tbl.SetMaxWidth(col, width)
			}
		}

		// Add varying number of rows to each table
		numRows := rng.Intn(4) + 2
		for j := 0; j < numRows; j++ {
			row := []string{
				fmt.Sprintf("%d-%d", i, j),
				generateRandomString(rng.Intn(20)+5, rng),
				generateRandomString(rng.Intn(15)+3, rng),
				fmt.Sprintf("%.2f", rng.Float64()*1000),
				generateStatus(rng),
			}
			tbl.AddRow(row)
		}

		tables[i] = tbl

		fmt.Printf("Table %d (Initial):\n", i+1)
		fmt.Println(tbl.Render())
		fmt.Println()
	}

	// Now add all tables to group and sync
	if tc.useGroup {
		for _, tbl := range tables {
			group.Add(tbl)
		}

		fmt.Println("=== After Syncing Column Widths ===")
		group.SyncColumnWidths()

		// Re-render all tables with synced widths
		syncedTables := group.GetTables()
		for i, tbl := range syncedTables {
			fmt.Printf("Table %d (Synced):\n", i+1)
			fmt.Println(tbl.Render())
			fmt.Println()
		}
	}
}

func test3HeaderHighlighting(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"ID", "Product Name", "Category", "Stock", "Price", "Status"}

	tbl := table.RapidFortTable(headers)
	tbl.SetHighlightedHeaders(tc.highlightedHeaders)
	tbl.SetDimBorder(tc.dimBorder)
	tbl.SetFillWidth(tc.fillWidth)

	// Add sample data
	tbl.AddRow([]string{"1", "Laptop Pro X1", "Electronics", "25", "$1,299.99", "Active"})
	tbl.AddRow([]string{"2", "Wireless Mouse", "Accessories", "150", "$29.99", "Active"})
	tbl.AddRow([]string{"3", "USB-C Hub", "Accessories", "75", "$49.99", "Low Stock"})
	tbl.AddDescriptionWithTitle(0, "Advisory", "High demand item, consider restocking")
	tbl.AddDescriptionWithTitle(2, "Notes", "Running low on inventory, order more from supplier")

	// FIXED: Removed \n from end of Println
	fmt.Println("Headers [1], [3], and [4] are highlighted in bold:")
	fmt.Println()
	fmt.Println(tbl.Render())
}

func test4MixedHighlightingAndStyling(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Region", "Sales Q1", "Growth %", "Target", "Status", "Notes"}

	tbl := table.RapidFortTable(headers)
	tbl.SetHighlightedHeaders(tc.highlightedHeaders)
	tbl.SetDimBorder(tc.dimBorder)
	tbl.SetFillWidth(tc.fillWidth)

	// Apply max widths
	if tc.maxColumnWidths != nil {
		for col, width := range tc.maxColumnWidths {
			tbl.SetMaxWidth(col, width)
		}
	}

	// Add sample data
	tbl.AddRow([]string{"North", "$125,000", "+15%", "$120,000", "Exceeds", "Great performance"})
	tbl.AddRow([]string{"South", "$85,000", "-5%", "$90,000", "Below", "Needs improvement"})
	tbl.AddRow([]string{"East", "$200,000", "+25%", "$180,000", "Exceeds", "Outstanding results"})
	tbl.AddDescription(1, "Performance review scheduled for next week")

	// FIXED: Removed \n from end of Println
	fmt.Println("Testing mix of dim borders and highlighted headers [0], [2], and [5]:")
	fmt.Println()
	fmt.Println(tbl.Render())
}

func test5GroupedTablesWithHighlighting(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"ID", "Name", "Category", "Value", "Status"}

	tables := make([]*table.Table, tc.numTables)

	for i := 0; i < tc.numTables; i++ {
		tbl := table.RapidFortTable(headers)
		tbl.SetHighlightedHeaders(tc.highlightedHeaders)
		tbl.SetDimBorder(tc.dimBorder)
		tbl.SetFillWidth(tc.fillWidth)

		// Add sample data for each table
		for j := 0; j < 3; j++ {
			row := []string{
				fmt.Sprintf("%d-%d", i, j),
				fmt.Sprintf("Item %d", j+1),
				generateCategory(rng),
				fmt.Sprintf("$%.2f", rng.Float64()*1000),
				generateStatus(rng),
			}
			tbl.AddRow(row)
		}

		tables[i] = tbl
		group.Add(tbl)
	}

	// Sync widths
	group.SyncColumnWidths()

	// Display tables
	syncedTables := group.GetTables()
	for i, tbl := range syncedTables {
		fmt.Printf("Table %d (Headers [1] and [3] are highlighted):\n", i+1)
		fmt.Println(tbl.Render())
		fmt.Println()
	}
}

func test6AllHeadersHighlighted(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"A", "B", "C", "D", "E", "F"}

	tbl := table.RapidFortTable(headers)
	tbl.SetHighlightedHeaders(tc.highlightedHeaders)

	// Add simple data
	tbl.AddRow([]string{"1", "2", "3", "4", "5", "6"})
	tbl.AddRow([]string{"X", "Y", "Z", "A", "B", "C"})

	// FIXED: Removed \n from end of Println
	fmt.Println("All headers highlighted:")
	fmt.Println()
	fmt.Println(tbl.Render())
}

func test7DynamicHeaderHighlighting(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Test", "Name", "Result", "Score", "Status", "Time"}

	tbl := table.RapidFortTable(headers)

	// Add some data
	tbl.AddRow([]string{"Unit Test", "DataStore", "Pass", "100%", "Green", "0.5s"})
	tbl.AddRow([]string{"Integration", "API", "Pass", "95%", "Green", "1.2s"})
	tbl.AddRow([]string{"E2E", "UI Flow", "Fail", "60%", "Red", "5.0s"})

	fmt.Println("Step 1: No highlighting")
	fmt.Println(tbl.Render())
	fmt.Println()

	fmt.Println("Step 2: Add highlighting to first two headers")
	tbl.SetHighlightedHeaders([]int{0, 1})
	fmt.Println(tbl.Render())
	fmt.Println()

	fmt.Println("Step 3: Add more highlights")
	tbl.AddHighlightedHeader(2)
	tbl.AddHighlightedHeader(4)
	fmt.Println(tbl.Render())
	fmt.Println()

	fmt.Println("Step 4: Clear all highlights")
	tbl.ClearHighlightedHeaders()
	fmt.Println(tbl.Render())
	fmt.Println()

	fmt.Println("Step 5: Highlight only status columns")
	tbl.SetHighlightedHeaders([]int{4})
	fmt.Println(tbl.Render())
}

// Helper functions for generating test data

func generateRandomString(length int, rng *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomDescription(rng *rand.Rand) string {
	descriptions := []string{
		"This is a sample description for testing line wrapping behavior",
		"Another description with different length content",
		"Short desc",
		"A very long description that should definitely wrap across multiple lines in the table cell",
		"Description with special characters: @#$%^&*()",
		"Multiple lines With line breaks Testing formatting Multiple lines With line breaks Testing formatting Multiple lines With line breaks Testing formatting\n Multiple lines With line breaks Testing formatting",
	}
	return descriptions[rng.Intn(len(descriptions))]
}

func generateStatus(rng *rand.Rand) string {
	statuses := []string{"Active", "Pending", "Completed", "Failed", "In Progress", "Cancelled"}
	return statuses[rng.Intn(len(statuses))]
}

func generateCategory(rng *rand.Rand) string {
	categories := []string{"Electronics", "Books", "Clothing", "Food", "Tools", "Other"}
	return categories[rng.Intn(len(categories))]
}

func simple() {
    // Create a table
    tbl := table.RapidFortTable([]string{"Name", "Age", "City"})
    
    // Add rows
    tbl.AddRow([]string{"Alice", "30", "New York"})
    tbl.AddRow([]string{"Bob", "25", "San Francisco"})
    
    // Add description to a row
    tbl.AddDescription(0, "Special customer discount applied")
    
    // Render and print
    fmt.Println(tbl.Render())
}
