package components

import (
	"fmt"
	"strconv"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

// TabItem represents a single tab in the TabLayout
type TabItem struct {
	ID      string
	Title   string
	Content string
	Active  bool
}

// TabLayoutConfig holds configuration for the TabLayout component
type TabLayoutConfig struct {
	Items        []TabItem
	Class        string
	Style        string
	NavClass     string
	ContentClass string
	Padding      int
	Margin       int
}

// NewTabLayout creates a new tab layout component and returns it as hb.TagInterface
func NewTabLayout(config TabLayoutConfig) hb.TagInterface {
	// Create container
	container := hb.Div().
		ClassIf(!lo.IsEmpty(config.Class), config.Class).
		StyleIf(!lo.IsEmpty(config.Style), config.Style).
		StyleIf(config.Padding > 0, "padding: "+strconv.Itoa(config.Padding)+"px;").
		StyleIf(config.Margin > 0, "margin: "+strconv.Itoa(config.Margin)+"px;")
	
	// Create tab navigation
	navClass := "nav nav-tabs"
	if !lo.IsEmpty(config.NavClass) {
		navClass = config.NavClass
	}
	
	nav := hb.NewTag("ul").
		Class(navClass).
		Attr("role", "tablist")
	
	// Create tab content container
	contentClass := "tab-content mt-2"
	if !lo.IsEmpty(config.ContentClass) {
		contentClass = config.ContentClass
	}
	
	contentContainer := hb.Div().Class(contentClass)
	
	// Process each tab
	for i, item := range config.Items {
		// Ensure ID is valid
		id := item.ID
		if lo.IsEmpty(id) {
			id = fmt.Sprintf("tab-%d", i)
		}
		
		// Create nav item
		navItem := hb.NewTag("li").
			Class("nav-item").
			Attr("role", "presentation")
		
		// Create nav link
		navLink := hb.Button().
			Class("nav-link").
			ClassIf(item.Active, "active").
			ID(id + "-tab").
			Data("bs-toggle", "tab").
			Data("bs-target", "#" + id).
			Type("button").
			Attr("role", "tab").
			Attr("aria-controls", id).
			Attr("aria-selected", fmt.Sprintf("%t", item.Active)).
			HTML(item.Title)
		
		navItem.Child(navLink)
		nav.Child(navItem)
		
		// Create tab content
		tabPane := hb.Div().
			Class("tab-pane fade").
			ClassIf(item.Active, "show active").
			ID(id).
			Attr("role", "tabpanel").
			Attr("aria-labelledby", id + "-tab").
			HTML(item.Content)
		
		contentContainer.Child(tabPane)
	}
	
	container.Child(nav)
	container.Child(contentContainer)
	
	return container
}
