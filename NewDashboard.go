package dashboard

func NewDashboard(config Config) Dashboard {
	dashboard := Dashboard{}
	dashboard.menu = config.Menu
	dashboard.useSmartMenu = config.UseSmartMenu
	dashboard.useMetisMenu = config.UseMetisMenu
	return dashboard
}
