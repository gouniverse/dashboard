package components

import (
	"strconv"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

// GridItem represents a single cell in the grid
type GridItem struct {
	Content     string
	ColumnClass string // e.g., "col-md-6", "col-lg-4", etc.
}

// GridLayoutConfig holds configuration for the GridLayout component
type GridLayoutConfig struct {
	Rows           [][]GridItem
	ContainerClass string
	RowClass       string
	Style          string
	Padding        int
	Margin         int
	Gutters        int // Bootstrap gutter spacing (g-0 to g-5)
}

// NewGridComponent creates a new grid layout component and returns it as hb.TagInterface
func NewGridComponent(config GridLayoutConfig) hb.TagInterface {
	// Create container
	container := hb.Div().
		ClassIf(!lo.IsEmpty(config.ContainerClass), config.ContainerClass).
		StyleIf(!lo.IsEmpty(config.Style), config.Style).
		StyleIf(config.Padding > 0, "padding: "+strconv.Itoa(config.Padding)+"px;").
		StyleIf(config.Margin > 0, "margin: "+strconv.Itoa(config.Margin)+"px;")
	
	// Process each row
	for _, rowItems := range config.Rows {
		// Set row class with gutters if specified
		rowClass := "row"
		if !lo.IsEmpty(config.RowClass) {
			rowClass = config.RowClass
		} else if config.Gutters > 0 {
			rowClass = "row g-" + strconv.Itoa(config.Gutters)
		}
		
		row := hb.Div().Class(rowClass)
		
		// Process each column in the row
		for _, item := range rowItems {
			colClass := "col"
			if !lo.IsEmpty(item.ColumnClass) {
				colClass = item.ColumnClass
			}
			
			col := hb.Div().
				Class(colClass).
				HTML(item.Content)
				
			row.Child(col)
		}
		
		container.Child(row)
	}
	
	return container
}
