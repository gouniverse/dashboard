package dashboard

import "github.com/samber/lo"

// themeNameVerifyAndFix verifies the theme name against the list of theme names
// in the case not found returns default
func themeNameVerifyAndFix(themeName string) string {
	return lo.
		If(lo.Contains(lo.Keys(themesLight), themeName), themeName).
		ElseIf(lo.Contains(lo.Keys(themesDark), themeName), themeName).
		Else("bootstrap")
}
