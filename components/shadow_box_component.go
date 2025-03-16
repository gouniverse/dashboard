package components

import (
	"strconv"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

// ShadowBoxConfig holds configuration for the ShadowBoxComponent component
type ShadowBoxConfig struct {
	Content string
	Class   string
	Style   string
	Padding int
	Margin  int
}

// NewShadowBoxComponent creates a new shadow box component with shadow, padding, and margin
func NewShadowBoxComponent(config ShadowBoxConfig) hb.TagInterface {
	// Create the layout div
	return hb.Div().
		Class("shadow").
		ClassIf(!lo.IsEmpty(config.Class), config.Class).
		StyleIf(config.Padding > 0, "padding: "+strconv.Itoa(config.Padding)+"px;").
		StyleIf(config.Margin > 0, "margin: "+strconv.Itoa(config.Margin)+"px;").
		StyleIf(!lo.IsEmpty(config.Style), config.Style).
		HTML(config.Content)
}
