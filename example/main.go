package main

import (
	"fmt"

	"github.com/rapidfort/table" // Replace with your actual module path
)

// Color constants for ANSI styling
const (
	RED       = "\x1b[31m"
	GREEN     = "\x1b[32m"
	YELLOW    = "\x1b[33m"
	BLUE      = "\x1b[34m"
	MAGENTA   = "\x1b[35m"
	CYAN      = "\x1b[36m"
	GRAY      = "\x1b[90m"
	RESET     = "\x1b[0m"
	BOLD      = "\x1b[1m"
	DIM       = "\x1b[2m"
	ITALIC    = "\x1b[3m"
	UNDERLINE = "\x1b[4m"
	BLINK     = "\x1b[5m"
	REVERSE   = "\x1b[7m"
	HIDDEN    = "\x1b[8m"
	STRIKE    = "\x1b[9m"
)

// Background colors
const (
	BG_BLACK   = "\x1b[40m"
	BG_RED     = "\x1b[41m"
	BG_GREEN   = "\x1b[42m"
	BG_YELLOW  = "\x1b[43m"
	BG_BLUE    = "\x1b[44m"
	BG_MAGENTA = "\x1b[45m"
	BG_CYAN    = "\x1b[46m"
	BG_WHITE   = "\x1b[47m"
)

// Bright colors
const (
	BRIGHT_RED     = "\x1b[91m"
	BRIGHT_GREEN   = "\x1b[92m"
	BRIGHT_YELLOW  = "\x1b[93m"
	BRIGHT_BLUE    = "\x1b[94m"
	BRIGHT_MAGENTA = "\x1b[95m"
	BRIGHT_CYAN    = "\x1b[96m"
	BRIGHT_WHITE   = "\x1b[97m"
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
	simpleGroup()

}

func simple() {
	// Create a table
	tbl := table.RapidFortTable([]string{"Name", "Age", "City"})

	// Add rows
	tbl.AddRow([]string{"Alice", "30", "New York"})
	tbl.AddRow([]string{"Bob", "25", "San Francisco"})

	// Add description to a row
	tbl.AddDescriptionWithTitle(0, "A", "Special customer discount applied")
	tbl.AddDescriptionWithTitle(0, "B", "Loyalty program: Gold member")

	// Render and print
	fmt.Println(tbl.Render())
}

func simpleGroup() {
	// pick your favourite ANSI colours
	red := "\x1b[31m"
	green := "\x1b[32m"
	blue := "\x1b[34m"
	bold := "\x1b[1m"
	dim := "\x1b[2m"
	reset := "\x1b[0m"

	headers := []string{
		bold + "ID" + reset,
		bold + "Message" + reset,
		bold + "Status" + reset,
	}

	// Table A
	tA := table.RapidFortTable(headers)
	//tA.SetDimBorder(true)           // dim the box-drawing lines
	//tA.SetHeaderHighlighting(false) // we've hard-coloured headers already

	// colour individual cells however you like:
	tA.AddRow([]string{
		red + "1" + reset,
		green + "All systems go" + reset,
		blue + "OK" + reset,
	})

	longDesc := dim + " This is a very long advisory note that “wraps” across multiple lines " +
		"inside the table, dimmed so it doesn’t compete with your main data." + reset
	tA.AddDescriptionWithTitle(
		0,
		bold+blue+"ADVISORY"+reset, // coloured description title
		longDesc,
	)
	// add a second advisory note
	tA.AddDescription(0, green+"Additional advisory: maintenance scheduled"+reset)

	// Table B
	tB := table.RapidFortTable(headers)
	//tB.SetDimBorder(true)

	tB.AddRow([]string{
		//red + "1" + reset,
		"1",
		green + "Partial outage" + reset,
		red + "FAIL" + reset,
	})
	tB.AddDescriptionWithTitle(
		0,
		bold+red+"ERROR"+reset,
		dim+"Immediate attention required!"+reset,
	)
	// add a follow‑up error note
	tB.AddDescription(0, red+"Follow-up: support notified"+reset)

	tB.AddRow([]string{
		red + "2" + reset,
		green + "Partial outage" + reset,
		red + "FAIL" + reset,
	})
	tB.AddDescriptionWithTitle(
		1,
		bold+red+"ERROR"+reset,
		dim+"Immediate attention required!"+reset,
	)
	// add a repeat error note
	tB.AddDescription(1, red+"Urgent: escalate to on-call"+reset)

	// Put them in a group so columns line up
	grp := table.NewGroup()
	grp.Add(tA)
	grp.Add(tB)
	grp.SyncColumnWidths()

	// print coloured titles
	fmt.Println(bold + "=== Table A: System Health ===" + reset)
	fmt.Println(tA.Render())

	fmt.Println(bold + "=== Table B: Error States ===" + reset)
	fmt.Println(tB.Render())
}
