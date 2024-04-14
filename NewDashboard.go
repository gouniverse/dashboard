package dashboard

// NewDashboard creates a new Dashboard instance based on the given configuration.
//
// The function takes a Config struct as its parameter and returns a pointer to a Dashboard struct.
func NewDashboard(config Config) *Dashboard {
	if config.MenuType == "" {
		config.MenuType = MENU_TYPE_OFFCANVAS // default
	}

	if config.Theme == "" && config.HTTPRequest != nil {
		config.Theme = ThemeNameRetrieveFromCookie(config.HTTPRequest)
	}

	if config.NavbarBackgroundColorMode == "" {
		config.NavbarBackgroundColorMode = "dark"
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
	dashboard.scripts = config.Scripts
	dashboard.scriptURLs = config.ScriptURLs
	dashboard.styles = config.Styles
	dashboard.styleURLs = config.StyleURLs
	dashboard.menuType = config.MenuType
	dashboard.redirectUrl = config.RedirectUrl
	dashboard.redirectTime = config.RedirectTime
	dashboard.themeHandlerUrl = config.ThemeHandlerUrl
	dashboard.theme = config.Theme
	dashboard.themesRestrict = config.ThemesRestrict
	dashboard.uncdnHandlerEndpoint = config.UncdnHandlerEndpoint
	dashboard.menu = config.Menu
	dashboard.quickAccessMenu = config.QuickAccessMenu
	dashboard.user = config.User
	dashboard.userMenu = config.UserMenu
	return dashboard
}
