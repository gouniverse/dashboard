package dashboard

import (
	"strings"

	"github.com/gouniverse/hb"
)

func styles(stylesArray []string) string {
	requiredStyles := []string{
		// uncdn.BootstrapCss523(),
	}

	stylesArray = append(requiredStyles, stylesArray...)

	styles := ""
	for _, style := range stylesArray {
		if style == "<style></style>" {
			continue
		}

		if strings.HasPrefix(style, "http") || strings.HasPrefix(style, "//") {
			style = hb.StyleURL(style).ToHTML()
		} else if !strings.HasPrefix(style, "<style") {
			style = hb.Style(style).ToHTML()
		}

		styles += style
	}
	return styles
}
