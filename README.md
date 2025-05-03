# RapidFort Table

A flexible Go library for creating beautiful ASCII tables with multiple border styles.

## Features

- Multiple border styles (Square, Rounded, ASCII, None)
- Intelligent text wrapping
- Column alignment control
- Dynamic width adjustment
- Description support
- Table grouping for consistent column widths

## Installation

```bash
go get RapidFort.dev/table
```

## Quick Example

```go
package main

import (
    "fmt"
    "RapidFort.dev/table"
)

func main() {
    // Create a table
    tbl := table.NewTable([]string{"Name", "Age", "City"})
    
    // Add rows
    tbl.AddRow([]string{"Alice", "30", "New York"})
    tbl.AddRow([]string{"Bob", "25", "San Francisco"})
    
    // Change border style (optional)
    tbl.SetBorderStyle(table.BorderStyleRounded)
    
    // Render and print
    fmt.Println(tbl.Render())
}
```

## Border Styles

- `BorderStyleSquare`: Sharp box-drawing characters (default)
- `BorderStyleRounded`: Rounded box-drawing characters
- `BorderStyleASCII`: Simple ASCII characters
- `BorderStyleNone`: No borders

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Specify your license here]