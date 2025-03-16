package components

import (
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
	"strconv"
)

// CardConfig holds configuration for the Card component
type CardConfig struct {
	Content     string
	Title       string
	HeaderClass string
	BodyClass   string
	CardClass   string
	Style       string
	Padding     int
	Margin      int
}

// NewCard creates a new card component and returns it as hb.TagInterface
func NewCard(config CardConfig) hb.TagInterface {
	// Create card
	card := hb.Div().
		Class("card").
		ClassIf(!lo.IsEmpty(config.CardClass), config.CardClass).
		StyleIf(!lo.IsEmpty(config.Style), config.Style).
		StyleIf(config.Padding > 0, "padding: "+strconv.Itoa(config.Padding)+"px;").
		StyleIf(config.Margin > 0, "margin: "+strconv.Itoa(config.Margin)+"px;")

	// Add header if title is provided
	if !lo.IsEmpty(config.Title) {
		headerClass := "card-header"
		if !lo.IsEmpty(config.HeaderClass) {
			headerClass = config.HeaderClass
		}
		
		header := hb.Div().
			Class(headerClass).
			HTML(config.Title)
			
		card.Child(header)
	}

	// Add body with content
	bodyClass := "card-body"
	if !lo.IsEmpty(config.BodyClass) {
		bodyClass = config.BodyClass
	}
	
	body := hb.Div().
		Class(bodyClass).
		HTML(config.Content)
		
	card.Child(body)

	return card
}
