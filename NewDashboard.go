package dashboard

func NewDashboard(config Config) Dashboard {
	if config.MenuType == "" {
		config.MenuType = MENU_TYPE_OFFCANVAS // default
	}

	if config.ThemeName == "" && config.HTTPRequest != nil {
		config.ThemeName = ThemeNameRetrieveFromCookie(config.HTTPRequest)
	}

	config.ThemeName = themeNameVerifyAndFix(config.ThemeName)

	dashboard := Dashboard{}
	dashboard.Title = config.Title
	dashboard.Content = config.Content
	dashboard.FaviconURL = config.FaviconURL
	dashboard.LogoURL = config.LogoURL
	dashboard.Scripts = config.Scripts
	dashboard.ScriptURLs = config.ScriptURLs
	dashboard.Styles = config.Styles
	dashboard.StyleURLs = config.StyleURLs
	dashboard.MenuType = config.MenuType
	dashboard.RedirectUrl = config.RedirectUrl
	dashboard.RedirectTime = config.RedirectTime
	dashboard.ThemeHandlerUrl = config.ThemeHandlerUrl
	dashboard.ThemeName = config.ThemeName
	dashboard.UncdnHandlerEndpoint = config.UncdnHandlerEndpoint
	dashboard.menu = config.Menu
	dashboard.user = config.User
	dashboard.userMenu = config.UserMenu
	return dashboard
}
