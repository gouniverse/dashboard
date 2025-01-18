package dashboard

// MenuItem is a menu item for the dashboard.
type MenuItem struct {
	Title    string
	URL      string
	Target   string
	Icon     string
	Sequence int
	Children []MenuItem
}
