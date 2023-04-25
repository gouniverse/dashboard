package dashboard

type MenuItem struct {
	Title    string
	URL      string
	Target   string
	Icon     string
	Sequence int
	Children []MenuItem
}
