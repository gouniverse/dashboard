package dashboard

import (
	"net/http"
	"time"

	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// ThemeHandler checks for the supplied theme and sets the theme name in the session
func ThemeHandler(w http.ResponseWriter, r *http.Request) {
	themeName := utils.Req(r, "theme", "")
	redirect := utils.Req(r, "redirect", "/")

	themeName = themeNameVerifyAndFix(themeName)

	secureCookies := lo.Ternary(r.TLS == nil, false, true)
	cookie := http.Cookie{Name: THEME_COOKIE_KEY, Value: themeName, Path: "/", Secure: secureCookies, Expires: time.Now().Add(365 * 24 * time.Hour)}

	http.SetCookie(w, &cookie)

	r.AddCookie(&cookie)

	http.Redirect(w, r, redirect, http.StatusFound)
}
