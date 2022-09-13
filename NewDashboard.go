package dashboard

func NewDashboard(config Config) Dashboard {
	dashboard := Dashboard{}
	dashboard.Title = config.Title
	dashboard.Content = config.Content
	dashboard.FaviconURL = config.FaviconURL
	dashboard.LogoURL = config.LogoURL
	dashboard.Scripts = config.Scripts
	dashboard.ScriptURLs = config.ScriptURLs
	dashboard.Styles = config.Styles
	dashboard.RedirectUrl = config.RedirectUrl
	dashboard.RedirectTime = config.RedirectTime
	dashboard.menu = config.Menu
	dashboard.useSmartMenu = config.UseSmartMenu
	dashboard.useMetisMenu = config.UseMetisMenu
	return dashboard
}
