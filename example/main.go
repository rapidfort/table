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
	title           string
	description     string
	useGroup        bool
	numTables       int
	dimBorder       bool
	fillWidth       bool
	maxColumnWidths map[int]int
	hasDescriptions bool
}

func main() {
	// Initialize random source (modern way)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	testCases := []TestCase{
		{
			title:           "Test Case 1: Random Sized Columns",
			description:     "Testing columns with random sizes and content",
			useGroup:        false,
			numTables:       1,
			dimBorder:       false,
			fillWidth:       false,
			maxColumnWidths: nil,
			hasDescriptions: true,
		},
		{
			title:           "Test Case 2: Grouped Tables with Sync",
			description:     "Multiple tables with synced column widths",
			useGroup:        true,
			numTables:       3,
			dimBorder:       true,
			fillWidth:       false,
			maxColumnWidths: map[int]int{1: 15, 3: 20},
			hasDescriptions: false,
		},
	}

	// Run only the first two test cases for debugging
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
	}
}

func test1RandomSizedColumns(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Short", "MediumLength", "A Very Long Header Name", "Tiny", "Another Column", "Last"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)
	tbl.SetFillWidth(tc.fillWidth)

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
		"Multiple lines\nWith line breaks\nTesting formatting",
	}
	return descriptions[rng.Intn(len(descriptions))]
}

func generateStatus(rng *rand.Rand) string {
	statuses := []string{"Active", "Pending", "Completed", "Failed", "In Progress", "Cancelled"}
	return statuses[rng.Intn(len(statuses))]
}
