package dashboard

type Config struct {
	Menu         []MenuItem
	User         User
	Title        string
	Content      string
	FaviconURL   string
	LogoURL      string
	Scripts      []string
	ScriptURLs   []string
	Styles       []string
	StyleURLs    []string
	RedirectUrl  string
	RedirectTime string
	UseSmartMenu bool
	UseMetisMenu bool
}
