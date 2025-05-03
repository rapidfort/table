package main

import (
	"fmt"

	"github.com/rapidfort/table" // Replace with your actual module path
)

func main() {
	// Create a new table with headers matching your example
	headers := []string{"#", "CVE ID", "Severity", "Package", "Installed", "Fixed In"}
	t := table.RapidFortTable(headers)

	// Set column alignments (optional)
	t.SetAlignment(0, "center") // Center the # column
	t.SetAlignment(2, "center") // Center the Severity column

	// Add rows for each CVE
	t.AddRow([]string{"1", "CVE-2024-10041", "MEDIUM", "rf-libpam-modules", "1.7.0-1rfubu1", "RF fixed"})
	// Add description for the first row
	t.AddDescription(0, "Fixed in upstream 1.6.0\nKnown issue in PAM's session module add more content here so that we can see how it looks\n\nAdditional note: Ensure to monitor for updates on this issue")

	t.AddRow([]string{"2", "CVE-2024-10963", "MEDIUM", "rf-libpam-modules", "1.7.0-1rfubu1", "RF fixed"})
	// Add description for the second row
	t.AddDescription(1, "Patch cherry-picked\nCommit SHA: 940747f8")

	t.AddRow([]string{"3", "CVE-2016-20013", "LOW", "rf-libc-bin", "2.39-0ubuntu8.4", ""})
	// Add description for the third row with multiple lines
	t.AddDescription(2, "sha-crypt is O(n²); inherent to algorithm\nNot a glibc vulnerability — apps should limit password length\nAdditional note: Ensure to monitor for updates on this issue")

	// Render and print the table
	fmt.Println(t.Render())
}
