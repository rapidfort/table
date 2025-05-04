package main

import (
	"fmt"

	"github.com/rapidfort/table" // Replace with your actual module path
)

func main() {
	// Create a table group to manage consistent column widths
	group := table.NewGroup()

	// ---- First Table: High Severity CVEs ----
	headers1 := []string{"#", "CVE ID", "Severity", "Package", "Installed", "Fixed In"}

	table1 := table.RapidFortTable(headers1)
	table1.SetDimBorder(true)

	// Set column alignments
	table1.SetAlignment(0, "center") // Center the # column
	table1.SetAlignment(2, "center") // Center the Severity column

	// Add rows to the first table
	table1.AddRow([]string{"1", "CVE-2023-48795", "HIGH", "openssh", "8.9p1-3ubuntu0.1", "8.9p1-3ubuntu0.3"})
	table1.AddDescription(0, "Terrapin attack\nVulnerability in SSH protocol allowing chosen-ciphertext attacks against ChaCha20-Poly1305 and AES-GCM")

	table1.AddRow([]string{"2", "CVE-2023-4863", "HIGH", "libwebp", "1.0.4-3", "1.0.4-3.1"})
	//table1.AddDescription(1, "Heap buffer overflow in WebP library\nOut-of-bounds memory access in VP8L decoding")

	// Add the first table to the group
	group.Add(table1)

	// ---- Second Table: Medium Severity CVEs ----
	headers2 := []string{"#", "CVE ID", "Severity", "Package", "Installed", "Fixed In"}
	table2 := table.RapidFortTable(headers2)

	// Set column alignments
	table2.SetAlignment(0, "center") // Center the # column
	table2.SetAlignment(2, "center") // Center the Severity column

	// Add rows to the second table
	table2.AddRow([]string{"1", "CVE-2024-10041", "MEDIUM", "lles", "1.7.0-1", "1.7.0-2, 1.7.0-2, 1.7.0-2"})
	//table2.AddDescription(0, " Fixed in upstream 1.6.0Known issue in PAM's session module Fixed in upstream 1.6.0Known issue in PAM's session module Fixed in upstream 1.6.0\nKnown issue in PAM's session module Fixed in upstream 1.6.0\nKnown issue in PAM's session module ")

	table2.AddRow([]string{"2", "CVE-2024-10963", "MEDIUM", "libc-bin", "2.39-0ubuntu8.4", "Pending"})
	//table2.AddDescription(1, "Patch cherry-picked\nCommit SHA: 940747f8")

	// Add the second table to the group
	group.Add(table2)

	// Synchronize column widths across both tables
	group.SyncColumnWidths()

	// Render and print both tables
	fmt.Println("=== HIGH SEVERITY VULNERABILITIES ===")
	fmt.Println(table1.Render())
	fmt.Println("\n=== MEDIUM SEVERITY VULNERABILITIES ===")
	fmt.Println(table2.Render())
}
