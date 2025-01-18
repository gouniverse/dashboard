package dashboard

import (
	"strings"

	"github.com/samber/lo"
)

func uncdnThemeStyleURL(uncdnHandlerEndpoint string, theme string) string {
	themeStyle := lo.
		If(theme == "cerulean", "bs523"+theme).
		ElseIf(theme == "cosmo", "bs523"+theme).
		ElseIf(theme == "cyborg", "bs523"+theme).
		ElseIf(theme == "darkly", "bs523"+theme).
		ElseIf(theme == "flatly", "bs523"+theme).
		ElseIf(theme == "journal", "bs523"+theme).
		ElseIf(theme == "litera", "bs523"+theme).
		ElseIf(theme == "lumen", "bs523"+theme).
		ElseIf(theme == "lux", "bs523"+theme).
		ElseIf(theme == "materia", "bs523"+theme).
		ElseIf(theme == "minty", "bs523"+theme).
		ElseIf(theme == "morph", "bs523"+theme).
		ElseIf(theme == "pulse", "bs523"+theme).
		ElseIf(theme == "quartz", "bs523"+theme).
		ElseIf(theme == "sandstone", "bs523"+theme).
		ElseIf(theme == "simplex", "bs523"+theme).
		ElseIf(theme == "sketchy", "bs523"+theme).
		ElseIf(theme == "slate", "bs523"+theme).
		ElseIf(theme == "solar", "bs523"+theme).
		ElseIf(theme == "spacelab", "bs523"+theme).
		ElseIf(theme == "superhero", "bs523"+theme).
		ElseIf(theme == "united", "bs523"+theme).
		ElseIf(theme == "vapor", "bs523"+theme).
		ElseIf(theme == "yeti", "bs523"+theme).
		ElseIf(theme == "zephyr", "bs523"+theme).
		Else("bs523")

	return strings.TrimSuffix(uncdnHandlerEndpoint, "/") + "/" + themeStyle + ".css"
}
