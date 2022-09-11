package dashboard

func NewDashboard(config Config) Dashboard {
	dashboard := Dashboard{}
	dashboard.Title = config.Title
	dashboard.Content = config.Content
	dashboard.Scripts = config.Scripts
	dashboard.ScriptURLs = config.ScriptURLs
	dashboard.Styles = config.Styles
	dashboard.StyleURLs = config.StyleURLs
	dashboard.menu = config.Menu
	dashboard.useSmartMenu = config.UseSmartMenu
	dashboard.useMetisMenu = config.UseMetisMenu
	return dashboard
}
