package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

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

	// Initialize random source (modern way)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	testCases := []TestCase{
		{
			title:              "Test Case 11: Simple Color Coded Descriptions",
			description:        "Testing basic color coded descriptions",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 12: Complex Color Combinations",
			description:        "Testing multiple color codes in descriptions",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 13: All Text Styling Combinations",
			description:        "Testing all available text styling options",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 14: Background Colors in Descriptions",
			description:        "Testing background colors and combinations",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 15: Color Coding with AddDescriptionWithTitle",
			description:        "Testing titled descriptions with colors",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 16: Mixed Color Cells and Descriptions",
			description:        "Testing color coordination between cells and descriptions",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 17: Long Colored Descriptions",
			description:        "Testing text wrapping with colored content",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 18: Rainbow Styled Descriptions",
			description:        "Testing gradient and rainbow effects",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 19: Nested Color Styling",
			description:        "Testing nested color codes and resets",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
		{
			title:              "Test Case 20: Color Accessibility",
			description:        "Testing high contrast and accessible color combinations",
			useGroup:           false,
			numTables:          1,
			dimBorder:          false,
			fillWidth:          false,
			maxColumnWidths:    nil,
			hasDescriptions:    true,
			highlightedHeaders: nil,
		},
	}

	// Run new test cases only
	for _, tc := range testCases {
		if strings.HasPrefix(tc.title, "Test Case 1") { // Only run test cases 11-20
			runTestCase(tc, rng)
			fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
		}
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
	case "Test Case 11: Simple Color Coded Descriptions":
		test11SimpleColorCodedDescriptions(group, tc, rng)
	case "Test Case 12: Complex Color Combinations":
		test12ComplexColorCombinations(group, tc, rng)
	case "Test Case 13: All Text Styling Combinations":
		test13AllTextStylingCombinations(group, tc, rng)
	case "Test Case 14: Background Colors in Descriptions":
		test14BackgroundColorsInDescriptions(group, tc, rng)
	case "Test Case 15: Color Coding with AddDescriptionWithTitle":
		test15ColorCodingWithTitle(group, tc, rng)
	case "Test Case 16: Mixed Color Cells and Descriptions":
		test16MixedColorCellsAndDescriptions(group, tc, rng)
	case "Test Case 17: Long Colored Descriptions":
		test17LongColoredDescriptions(group, tc, rng)
	case "Test Case 18: Rainbow Styled Descriptions":
		test18RainbowStyledDescriptions(group, tc, rng)
	case "Test Case 19: Nested Color Styling":
		test19NestedColorStyling(group, tc, rng)
	case "Test Case 20: Color Accessibility":
		test20ColorAccessibility(group, tc, rng)
	}
}

func test11SimpleColorCodedDescriptions(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"ID", "Task", "Status", "Priority", "Due Date"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"1", "Design new logo", "In Progress", "High", "2024-01-01"})
	tbl.AddRow([]string{"2", "Update documentation", "Completed", "Medium", "2023-12-15"})
	tbl.AddRow([]string{"3", "Fix login bug", "Pending", "Critical", "2023-12-20"})
	tbl.AddRow([]string{"4", "Write unit tests", "In Progress", "High", "2023-12-25"})

	// Add simple colored descriptions
	tbl.AddDescription(0, GREEN+"On track for delivery"+RESET)
	tbl.AddDescription(1, BLUE+"Task completed successfully"+RESET)
	tbl.AddDescription(2, RED+"Urgent attention required"+RESET)
	tbl.AddDescription(3, YELLOW+"Behind schedule - needs review"+RESET)

	fmt.Println("Simple color coded descriptions:")
	fmt.Println(tbl.Render())
}

func test12ComplexColorCombinations(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"System", "Status", "Health", "Response Time", "Alerts"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"Database", "Running", "Healthy", "5ms", "0"})
	tbl.AddRow([]string{"API Server", "Running", "Warning", "150ms", "2"})
	tbl.AddRow([]string{"Cache Layer", "Error", "Critical", "N/A", "15"})
	tbl.AddRow([]string{"Load Balancer", "Running", "Healthy", "1ms", "0"})

	// Add complex colored descriptions
	tbl.AddDescription(0, GREEN+"Database performance "+BOLD+"excellent"+RESET+CYAN+" | CPU: 45% | Memory: 2.3GB"+RESET)
	tbl.AddDescription(1, YELLOW+"API response time "+UNDERLINE+"elevated"+RESET+" | Suggest review of slow endpoints")
	tbl.AddDescription(2, RED+"CRITICAL: "+RESET+BG_RED+"Service failure detected"+RESET+RED+" | Immediate action required"+RESET)
	tbl.AddDescription(3, GREEN+"Load balancer "+ITALIC+"optimal"+RESET+" | Traffic distribution: "+CYAN+"25% balanced"+RESET)

	fmt.Println("Complex color combinations in descriptions:")
	fmt.Println(tbl.Render())
}

func test13AllTextStylingCombinations(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Style", "Example", "Notes", "Impact", "Usage"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows with different styling
	tbl.AddRow([]string{"Bold", BOLD + "Important Text" + RESET, "High emphasis", "High", "Headers"})
	tbl.AddRow([]string{"Italic", ITALIC + "Emphasized Text" + RESET, "Medium emphasis", "Medium", "Notes"})
	tbl.AddRow([]string{"Underline", UNDERLINE + "Underlined Text" + RESET, "Low emphasis", "Low", "Links"})
	tbl.AddRow([]string{"Dim", DIM + "Dimmed Text" + RESET, "Deemphasize", "Low", "Secondary info"})
	tbl.AddRow([]string{"Combinations", BOLD + ITALIC + "Bold Italic" + RESET, "Combined styles", "High", "Headings"})

	// Add descriptions showing all styling options
	tbl.AddDescription(0, BOLD+"Bold style sample: "+RESET+"Use for important information")
	tbl.AddDescription(1, ITALIC+"Italic style sample: "+RESET+"Use for emphasis")
	tbl.AddDescription(2, UNDERLINE+"Underline style sample: "+RESET+"Rarely used, consider alternatives")
	tbl.AddDescription(3, DIM+"Dim style sample: "+RESET+"Good for metadata and secondary info")
	tbl.AddDescription(4, BOLD+ITALIC+UNDERLINE+"Combined styles: "+RESET+"Powerful, use sparingly")

	fmt.Println("All text styling combinations:")
	fmt.Println(tbl.Render())
}

func test14BackgroundColorsInDescriptions(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Alert Type", "Severity", "Count", "Source", "Action"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"Error", "High", "15", "Database", "Restart"})
	tbl.AddRow([]string{"Warning", "Medium", "8", "API", "Monitor"})
	tbl.AddRow([]string{"Info", "Low", "5", "Log", "Review"})
	tbl.AddRow([]string{"Critical", "Maximum", "3", "System", "Immediate"})

	// Add descriptions with background colors
	tbl.AddDescription(0, BG_RED+RED+" High priority alert - requires immediate attention "+RESET)
	tbl.AddDescription(1, BG_YELLOW+YELLOW+" Warning level - review within 24 hours "+RESET)
	tbl.AddDescription(2, BG_BLUE+CYAN+" Information only - no action required "+RESET)
	tbl.AddDescription(3, BG_RED+BRIGHT_WHITE+BOLD+" CRITICAL ALERT - ESCALATE IMMEDIATELY "+RESET)

	fmt.Println("Background colors in descriptions:")
	fmt.Println(tbl.Render())
}

func test15ColorCodingWithTitle(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Service", "Status", "Version", "Last Check", "Health"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"Auth Service", "Active", "v1.2.3", "2min ago", "OK"})
	tbl.AddRow([]string{"Payment API", "Active", "v2.0.1", "1min ago", "Warning"})
	tbl.AddRow([]string{"Email Service", "Down", "v1.5.2", "30min ago", "Error"})
	tbl.AddRow([]string{"Notification", "Active", "v3.1.0", "5min ago", "OK"})

	// Add descriptions with colored titles
	tbl.AddDescriptionWithTitle(0, GREEN+"System Status"+RESET, "All authentication services running normally. Response time under 100ms")
	tbl.AddDescriptionWithTitle(1, YELLOW+"Performance Alert"+RESET, "Payment processing experiencing "+YELLOW+"increased latency"+RESET+". Current avg: 1.2s")
	tbl.AddDescriptionWithTitle(2, RED+"OUTAGE ALERT"+RESET, "Email service down due to "+RED+BOLD+"database connection failure"+RESET+". Team investigating")
	tbl.AddDescriptionWithTitle(3, CYAN+"Info"+RESET, "Notification service updated successfully. All "+GREEN+"green"+RESET)

	fmt.Println("Colored titles with descriptions:")
	fmt.Println(tbl.Render())
}

func test16MixedColorCellsAndDescriptions(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Component", "Status", "Load", "Memory", "Alerts"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows with colored content
	tbl.AddRow([]string{
		"Frontend",
		GREEN + "Active" + RESET,
		"12%",
		"1.2GB",
		"0",
	})
	tbl.AddRow([]string{
		"API Gateway",
		YELLOW + "Degraded" + RESET,
		"" + RED + "85%" + RESET,
		"" + YELLOW + "3.8GB" + RESET,
		"" + YELLOW + "2" + RESET,
	})
	tbl.AddRow([]string{
		"Database",
		RED + "Critical" + RESET,
		"" + RED + "95%" + RESET,
		"" + RED + "7.5GB" + RESET,
		"" + RED + "12" + RESET,
	})

	// Add coordinated colored descriptions
	tbl.AddDescription(0, GREEN+"✓ All systems operational "+RESET+"| Last check: "+CYAN+"just now"+RESET)
	tbl.AddDescription(1, YELLOW+"⚠ Service degradation "+RESET+"| Cause: "+YELLOW+"High traffic"+RESET+" | "+CYAN+"Monitoring"+RESET)
	tbl.AddDescription(2, BG_RED+BRIGHT_WHITE+" CRITICAL SERVICE DOWN "+RESET+RED+" | Action: Failover to backup | ETA: 30min"+RESET)

	fmt.Println("Coordinated cell and description colors:")
	fmt.Println(tbl.Render())
}

func test17LongColoredDescriptions(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Log ID", "Timestamp", "Level", "Source", "Type"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"LOG-001", "2023-12-16 14:30:15", "INFO", "auth", "Login"})
	tbl.AddRow([]string{"LOG-002", "2023-12-16 14:30:22", "ERROR", "payment", "Transaction"})
	tbl.AddRow([]string{"LOG-003", "2023-12-16 14:30:45", "WARN", "database", "Connection"})

	// Add long colored descriptions with proper text wrapping
	tbl.AddDescription(0, BLUE+"User authentication successful. "+RESET+"IP: "+CYAN+"192.168.1.100"+RESET+", Device: "+GREEN+"Chrome/Windows"+RESET+", Location: "+YELLOW+"Mountain View, CA"+RESET+". Previous login: "+DIM+"3 days ago"+RESET)

	tbl.AddDescription(1, RED+"PAYMENT FAILED: "+RESET+"Transaction ID "+YELLOW+"TX-98765"+RESET+" declined by "+RED+"Visa"+RESET+" (card ending "+DIM+"4242"+RESET+"). Reason: "+RED+BOLD+"Insufficient funds"+RESET+". Customer notified via "+CYAN+"email"+RESET+" and "+CYAN+"SMS"+RESET)

	tbl.AddDescription(2, YELLOW+"Database connection pool warning: "+RESET+"Active connections: "+RED+"85/100"+RESET+" ("+YELLOW+"85%"+RESET+" capacity). Average response time increased to "+YELLOW+"250ms"+RESET+" (normal: "+GREEN+"50ms"+RESET+"). Recommendation: "+CYAN+"Add 2 more instances"+RESET)

	fmt.Println("Long colored descriptions with text wrapping:")
	fmt.Println(tbl.Render())
}

func test18RainbowStyledDescriptions(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Type", "Status", "Progress", "Performance", "Rating"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"Deploy", "Complete", "100%", "Excellent", "5/5"})
	tbl.AddRow([]string{"Build", "In Progress", "75%", "Good", "4/5"})
	tbl.AddRow([]string{"Test", "Failed", "45%", "Poor", "2/5"})

	// Rainbow gradient descriptions
	rainbowText := func(text string) string {
		colors := []string{RED, YELLOW, GREEN, CYAN, BLUE, MAGENTA}
		result := ""
		for i, char := range text {
			result += colors[i%len(colors)] + string(char) + RESET
		}
		return result
	}

	tbl.AddDescription(0, rainbowText("Deployment completed successfully across all regions"))
	tbl.AddDescription(1, YELLOW+"Build "+GREEN+"progressing "+CYAN+"smoothly"+BLUE+" with "+MAGENTA+"no "+RED+"issues"+RESET)
	tbl.AddDescription(2, RED+"Test failures detected in: "+YELLOW+"auth module"+GREEN+", API "+CYAN+"endpoints"+BLUE+", database "+MAGENTA+"connections"+RESET)

	fmt.Println("Rainbow styled descriptions:")
	fmt.Println(tbl.Render())
}

func test19NestedColorStyling(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Module", "Coverage", "Quality", "Tests", "Status"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows
	tbl.AddRow([]string{"User Auth", "85%", "A", "245/280", "Stable"})
	tbl.AddRow([]string{"Payment", "92%", "A+", "312/340", "Stable"})
	tbl.AddRow([]string{"Reporting", "65%", "B-", "130/200", "Needs Work"})

	// Add descriptions with nested color styling
	tbl.AddDescription(0, CYAN+"Test Coverage: "+RESET+"85% "+GREEN+"(Good)"+RESET+" | Remaining tests: "+YELLOW+"35"+RESET+" focusing on "+BLUE+"edge cases"+RESET)

	tbl.AddDescription(1, GREEN+"Excellent coverage: "+RESET+"92% "+BOLD+GREEN+"(Outstanding)"+RESET+" | "+CYAN+"All critical paths covered"+RESET+" | "+DIM+"Only minor UI tests pending"+RESET)

	tbl.AddDescription(2, YELLOW+"Coverage insufficient: "+RESET+"65% "+RED+"(Below threshold)"+RESET+" | Priority areas: "+BOLD+RED+"API authentication"+RESET+RED+", data validation, error handling"+RESET)

	fmt.Println("Nested color styling in descriptions:")
	fmt.Println(tbl.Render())
}

func test20ColorAccessibility(group *table.TableGroup, tc TestCase, rng *rand.Rand) {
	headers := []string{"Element", "Contrast", "Accessibility", "WCAG", "Notes"}

	tbl := table.RapidFortTable(headers)
	tbl.SetDimBorder(tc.dimBorder)

	// Add rows showing accessibility-focused colors
	tbl.AddRow([]string{"Error", "7.1:1", "AAA", "Pass", "High contrast"})
	tbl.AddRow([]string{"Warning", "4.8:1", "AA", "Pass", "Good contrast"})
	tbl.AddRow([]string{"Success", "8.2:1", "AAA", "Pass", "Excellent"})
	tbl.AddRow([]string{"Info", "3.5:1", "AA", "Pass", "Minimal contrast"})

	// Add descriptions that follow accessibility guidelines
	tbl.AddDescription(0, RED+BG_WHITE+" Error: High visibility with maximum contrast "+RESET+" (7.1:1 ratio)")
	tbl.AddDescription(1, YELLOW+BG_BLACK+" Warning: Good readability with sufficient contrast "+RESET+" (4.8:1 ratio)")
	tbl.AddDescription(2, GREEN+BG_BLACK+" Success: Optimized for both light and dark modes "+RESET+" (8.2:1 ratio)")
	tbl.AddDescription(3, BLUE+BG_WHITE+" Information: Meets WCAG AA standards "+RESET+" (3.5:1 ratio)")

	fmt.Println("Accessibility-focused color combinations:")
	fmt.Println(tbl.Render())
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
