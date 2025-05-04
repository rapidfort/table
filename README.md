# RapidFort Table

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/rapidfort/table.svg)](https://pkg.go.dev/github.com/rapidfort/table)

A flexible Go library for creating beautiful ASCII tables with box-drawing characters and advanced formatting options.

## Features

- 📚 Smart text wrapping and word splitting
- 🔄 Column alignment control (left, right, center)
- 📝 Description support with optional titles
- ✨ Header highlighting
- 🎨 Dim border styling for subtle tables
- 📊 Table grouping for consistent column widths
- 📐 Dynamic width adjustment and terminal detection
- 🔍 Support for borderless mode
- 💡 Fill width option for full terminal utilization

## Installation

```bash
go get github.com/rapidfort/table
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/rapidfort/table"
)

func main() {
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
```

This produces:
```
┌───────┬─────┬───────────────┐
│ Name  │ Age │ City          │
├───────┼─────┼───────────────┤
│ Alice │ 30  │ New York      │
│       ├─────┴───────────────┤
│       │   • Special         │
│       │     customer        │
│       │     discount        │
│       │     applied         │
├───────┼─────┬───────────────┤
│ Bob   │ 25  │ San Francisco │
└───────┴─────┴───────────────┘
```

## Advanced Usage

### Header Styling

```go
// Highlight specific headers
tbl.SetHighlightedHeaders([]int{0, 2}) // Highlights columns 0 and 2

// Always highlight all headers (bold style)
tbl.SetHeaderHighlighting(true)

// Add/remove highlighting dynamically
tbl.AddHighlightedHeader(1)
tbl.ClearHighlightedHeaders()
```

### Border and Style Control

```go
// Enable dim borders (light gray)
tbl.SetDimBorder(true)

// Remove borders completely
tbl.SetBorderless(true)

// Automatic terminal width detection (default)
tbl.SetConsoleWidth(80) // Manual override
```

### Column Formatting

```go
// Set column alignment
tbl.SetAlignment(0, "left")   // Column 0: left-aligned
tbl.SetAlignment(1, "right")  // Column 1: right-aligned
tbl.SetAlignment(2, "center") // Column 2: center-aligned

// Set maximum column widths
tbl.SetMaxWidth(1, 20)        // Column 1 max width: 20 chars

// Enable width filling
tbl.SetFillWidth(true)        // Expand to fill terminal width
```

### Table Groups

When multiple tables need consistent column widths:

```go
group := table.NewGroup()

// Create tables
table1 := table.RapidFortTable(headers)
table2 := table.RapidFortTable(headers)

// Add to group
group.Add(table1)
group.Add(table2)

// Sync column widths across all tables
group.SyncColumnWidths()

// Render tables with consistent widths
fmt.Println(table1.Render())
fmt.Println(table2.Render())
```

### Rich Descriptions

```go
// Add description with title
tbl.AddDescriptionWithTitle(rowIndex, "Advisory", "High demand item, consider restocking")

// Add simple description
tbl.AddDescription(rowIndex, "This item is currently on backorder")

// Multi-line descriptions are automatically wrapped
tbl.AddDescription(rowIndex, "This is a long description that will wrap across multiple lines based on the available space in the table")
```

## Examples

### Complete Feature Showcase

```go
// Create table with all features
tbl := table.RapidFortTable([]string{"ID", "Product", "Stock", "Status"})

// Configure styling
tbl.SetDimBorder(true)
tbl.SetHighlightedHeaders([]int{1, 3})
tbl.SetAlignment(2, "right")
tbl.SetMaxWidth(1, 15)

// Add data
tbl.AddRow([]string{"001", "Wireless Mouse", "125", "Active"})
tbl.AddRow([]string{"002", "USB-C Hub Pro", "25", "Low Stock"})
tbl.AddRow([]string{"003", "Mechanical Keyboard", "50", "Active"})

// Add descriptions
tbl.AddDescriptionWithTitle(1, "Alert", "Reorder soon to avoid stockout")
tbl.AddDescription(2, "Popular item, high demand expected")

fmt.Println(tbl.Render())
```

### Borderless Table

```go
tbl := table.RapidFortTable([]string{"Name", "Value"})
tbl.SetBorderless(true)
tbl.AddRow([]string{"Setting 1", "Enabled"})
tbl.AddRow([]string{"Setting 2", "Disabled"})
fmt.Println(tbl.Render())
```

## Documentation

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/rapidfort/table).

## Requirements

- Go 1.16 or higher
- Dependencies: `golang.org/x/term`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

MIT License

Copyright (c) 2025 RapidFort

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Acknowledgments

Built by the RapidFort team.