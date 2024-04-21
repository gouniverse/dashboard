package dashboard

import "net/http"

type ThemeNameContextKey struct{}

var THEME_COOKIE_KEY = "theme"

type Config struct {
	Content     string
	FaviconURL  string
	HTTPRequest *http.Request
	LogoURL     string
	Menu        []MenuItem
	MenuType    string

	// Optional. The URL of the logo image
	LogoImageURL string

	// Optional. Raw HTML of the logo, if set will be used instead of logoImageURL
	LogoRawHtml string

	// Optional. The redirect URL of the logo image
	LogoRedirectURL string

	// Optional The background color for the navbar: light, dark (default),  primary, secondary, success, warning, info, danger
	NavbarBackgroundColorMode string

	// Optional. The background color for the navbar (default none)
	NavbarBackgroundColor string

	// Optional. The text color for the navbar (default light)
	NavbarTextColor string

	// Optional. The URL of the login page to use (if user is not provided)
	LoginURL string

	// Optional. The URL of the register page to use (if user is not provided)
	RegisterURL string

	// Optional. Menu for Quick Access
	QuickAccessMenu []MenuItem
	RedirectTime    string
	RedirectUrl     string
	Scripts         []string
	ScriptURLs      []string
	Styles          []string
	StyleURLs       []string

	// Optional. The theme to be activated on the dashboard (default will be used otherwise)
	Theme string

	// Optional. Sets the default theme (THEME_DEFAULT will be used otherwise)
	// ThemeDefault string

	// Optional. The URL of the theme switcher endpoint to use
	ThemeHandlerUrl string

	// Optional. The themes to be visible in the theme switcher, the key is the theme, the value is the name (can be customized, default will be used otherwise)
	ThemesRestrict map[string]string

	Title                string
	UncdnHandlerEndpoint string
	User                 User
	UserMenu             []MenuItem
}
