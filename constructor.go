package dashboard

// NewDashboard creates a new Dashboard instance based on
// the given configuration.
//
// The function takes a Config struct as its parameter
// and returns a pointer to a Dashboard struct.
//
// Parameters:
// - config: A Config struct containing the configuration for the dashboard.
//
// Returns:
// - A pointer to a Dashboard struct.
func NewDashboard(config Config) *Dashboard {
	if config.MenuType == "" {
		config.MenuType = MENU_TYPE_OFFCANVAS // default
	}

	if config.Theme == "" && config.HTTPRequest != nil {
		config.Theme = ThemeNameRetrieveFromCookie(config.HTTPRequest)
	}

	config.Theme = themeNameVerifyAndFix(config.Theme)

	dashboard := &Dashboard{}
	dashboard.title = config.Title
	dashboard.content = config.Content
	dashboard.faviconURL = config.FaviconURL
	dashboard.logoImageURL = config.LogoImageURL
	dashboard.logoRawHtml = config.LogoRawHtml
	dashboard.logoRedirectURL = config.LogoRedirectURL
	dashboard.loginURL = config.LoginURL
	dashboard.registerURL = config.RegisterURL
	dashboard.navbarBackgroundColorMode = config.NavbarBackgroundColorMode
	dashboard.navbarBackgroundColor = config.NavbarBackgroundColor
	dashboard.navbarTextColor = config.NavbarTextColor
	dashboard.scripts = config.Scripts
	dashboard.scriptURLs = config.ScriptURLs
	dashboard.styles = config.Styles
	dashboard.styleURLs = config.StyleURLs
	dashboard.redirectUrl = config.RedirectUrl
	dashboard.redirectTime = config.RedirectTime
	dashboard.themeHandlerUrl = config.ThemeHandlerUrl
	dashboard.theme = config.Theme
	dashboard.themesRestrict = config.ThemesRestrict
	dashboard.uncdnHandlerEndpoint = config.UncdnHandlerEndpoint

	dashboard.menuItems = config.MenuItems
	dashboard.menuShowText = config.MenuShowText
	dashboard.menuType = config.MenuType

	dashboard.quickAccessMenu = config.QuickAccessMenu

	dashboard.user = config.User
	dashboard.userMenu = config.UserMenu

	return dashboard
}
