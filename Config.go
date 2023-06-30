package dashboard

import "net/http"

type ThemeNameContextKey struct{}

var THEME_COOKIE_KEY = "theme"

type Config struct {
	Content              string
	FaviconURL           string
	HTTPRequest          *http.Request
	LogoURL              string
	Menu                 []MenuItem
	MenuType             string
	RedirectTime         string
	RedirectUrl          string
	Scripts              []string
	ScriptURLs           []string
	Styles               []string
	StyleURLs            []string
	ThemeName            string
	ThemeHandlerUrl      string
	Title                string
	UncdnHandlerEndpoint string
	User                 User
	UserMenu             []MenuItem
}
