package main

import (
	"fmt"
	"net/http"

	"github.com/gouniverse/dashboard/components"
	"github.com/gouniverse/hb"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/shadow-box", shadowBoxHandler)
	http.HandleFunc("/card", cardHandler)
	http.HandleFunc("/tabs", tabLayoutHandler)
	http.HandleFunc("/grid", gridLayoutHandler)
	http.HandleFunc("/combined", combinedComponentsHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Dashboard Components Examples</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
	</head>
	<body>
		<div class="container mt-5">
			<h1>Dashboard Components Examples</h1>
			<ul class="list-group mt-4">
				<li class="list-group-item"><a href="/shadow-box">Shadow Box Component</a></li>
				<li class="list-group-item"><a href="/card">Card Component</a></li>
				<li class="list-group-item"><a href="/tabs">Tab Layout Component</a></li>
				<li class="list-group-item"><a href="/grid">Grid Component</a></li>
				<li class="list-group-item"><a href="/combined">Combined Components</a></li>
			</ul>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func shadowBoxHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageWrapper("Shadow Box Component Example", ShadowBoxExample().ToHTML())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageWrapper("Card Component Example", CardExample().ToHTML())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func tabLayoutHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageWrapper("Tab Layout Component Example", TabLayoutExample().ToHTML())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func gridLayoutHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageWrapper("Grid Component Example", GridLayoutExample().ToHTML())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func combinedComponentsHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageWrapper("Combined Components Example", CombinedComponentsExample().ToHTML())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func getPageWrapper(title string, content string) string {
	return `
	<!DOCTYPE html>
	<html>
	<head>
		<title>` + title + `</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
	</head>
	<body>
		<div class="container mt-5">
			<h1>` + title + `</h1>
			<div class="mb-4">
				<a href="/" class="btn btn-outline-primary">Back to Examples</a>
			</div>
			<div>
				` + content + `
			</div>
		</div>
	</body>
	</html>
	`
}

// ShadowBoxExample demonstrates using the ShadowBoxComponent
func ShadowBoxExample() hb.TagInterface {
	// Create a shadow box with shadow, padding, and margin
	return components.NewShadowBoxComponent(components.ShadowBoxConfig{
		Content: "<p>This is content in a shadow box with shadow, padding, and margin.</p>",
		Padding: 15,
		Margin: 10,
	})
}

// CardExample demonstrates using the Card component
func CardExample() hb.TagInterface {
	// Create a card with a title and content
	return components.NewCard(components.CardConfig{
		Title:   "Card Title",
		Content: "<p>This is content inside a Bootstrap card component.</p>",
		Margin: 10,
	})
}

// TabLayoutExample demonstrates using the TabLayout component
func TabLayoutExample() hb.TagInterface {
	// Create a tab layout with multiple tabs
	return components.NewTabLayout(components.TabLayoutConfig{
		Items: []components.TabItem{
			{
				ID:      "tab1",
				Title:   "First Tab",
				Content: "<p>Content of the first tab.</p>",
				Active:  true,
			},
			{
				ID:      "tab2",
				Title:   "Second Tab",
				Content: "<p>Content of the second tab.</p>",
				Active:  false,
			},
			{
				ID:      "tab3",
				Title:   "Third Tab",
				Content: "<p>Content of the third tab.</p>",
				Active:  false,
			},
		},
		Margin: 10,
	})
}

// GridLayoutExample demonstrates using the GridLayout component
func GridLayoutExample() hb.TagInterface {
	// Create a grid layout with two rows
	row1 := []components.GridItem{
		{Content: "<div class='p-3 bg-light'>Column 1</div>", ColumnClass: "col-md-6"},
		{Content: "<div class='p-3 bg-light'>Column 2</div>", ColumnClass: "col-md-6"},
	}
	
	// Second row with three columns of different widths
	row2 := []components.GridItem{
		{Content: "<div class='p-3 bg-light'>Column 1</div>", ColumnClass: "col-md-3"},
		{Content: "<div class='p-3 bg-light'>Column 2</div>", ColumnClass: "col-md-6"},
		{Content: "<div class='p-3 bg-light'>Column 3</div>", ColumnClass: "col-md-3"},
	}
	
	return components.NewGridComponent(components.GridLayoutConfig{
		Rows: [][]components.GridItem{row1, row2},
		Gutters: 3,
		Margin: 10,
	})
}

// CombinedComponentsExample demonstrates combining multiple components
func CombinedComponentsExample() hb.TagInterface {
	// Create tabs
	tabs := components.NewTabLayout(components.TabLayoutConfig{
		Items: []components.TabItem{
			{
				ID:      "tab1",
				Title:   "First Tab",
				Content: "<p>Content of the first tab.</p>",
				Active:  true,
			},
			{
				ID:      "tab2",
				Title:   "Second Tab",
				Content: "<p>Content of the second tab.</p>",
				Active:  false,
			},
		},
	})
	
	// Put the tabs inside a card
	card := components.NewCard(components.CardConfig{
		Title:   "Card with Tabs",
		Content: tabs.ToHTML(),
	})
	
	// Put the card inside a shadow box
	return components.NewShadowBoxComponent(components.ShadowBoxConfig{
		Content: card.ToHTML(),
		Padding: 15,
		Margin: 10,
	})
}
