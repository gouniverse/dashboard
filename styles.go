package dashboard

import (
	"strings"

	"github.com/gouniverse/hb"
)

func styles(stylesArray []string) string {
	requiredStyles := []string{
		// templates.BootstrapCss513(),
		// templates.BootstrapIconsCss(),
		// guestStylesFromFile(),
		// guestStyles(),
	}

	stylesArray = append(requiredStyles, stylesArray...)

	styles := ""
	for _, style := range stylesArray {
		if style == "<style></style>" {
			continue
		}

		if strings.HasPrefix(style, "http") || strings.HasPrefix(style, "//") {
			style = hb.NewStyleURL(style).ToHTML()
		} else if !strings.HasPrefix(style, "<style") {
			style = hb.NewStyle(style).ToHTML()
		}

		styles += style
	}
	return styles
}
