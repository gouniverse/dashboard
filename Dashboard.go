package dashboard

import (
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/uncdn"
	"github.com/gouniverse/utils"
)

type Dashboard struct {
	menu         []MenuItem
	user         User
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
	useMetisMenu bool // TODO
	useSmartMenu bool
}

type MenuItem struct {
	Title    string
	URL      string
	Target   string
	Icon     string
	Sequence int
	Children []MenuItem
}

type User struct {
	FirstName string
	LastName  string
}

func (d Dashboard) layout() string {
	content := d.Content
	layout := hb.NewBorderLayout()
	top := top()
	if top != "" {
		layout.AddTop(hb.NewHTML(top), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	}
	layout.AddLeft(hb.NewHTML(d.left()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	//layout.AddBottom(hb.NewHTML("BOTTOM"), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	// layout.AddRight(hb.NewHTML("RIGHT"), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	layout.AddCenter(hb.NewHTML(d.center(content)), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	return layout.ToHTML()
}

type DashboardTemplateParams struct {
	Title        string
	Content      string
	Scripts      []string
	ScriptURLs   []string
	Styles       []string
	StyleURLs    []string
	RedirectUrl  string
	RedirectTime string
}

func (d Dashboard) SetUser(user User) Dashboard {
	d.user = user
	return d
}

func (d Dashboard) SetMenu(menuItems []MenuItem) Dashboard {
	d.menu = menuItems
	return d
}

func (d Dashboard) ToHTML() string {
	styleURLs := []string{
		cdn.BootstrapIconsCss191(),
	}

	additionalStyles := []string{}

	if d.useSmartMenu {
		// smartMenuStyleURLs := []string{
		// 	"https://cdnjs.cloudflare.com/ajax/libs/jquery.smartmenus/1.1.0/css/sm-core-css.css",
		// 	"https://cdnjs.cloudflare.com/ajax/libs/jquery.smartmenus/1.1.0/css/sm-simple/sm-simple.min.css",
		// 	"https://cdnjs.cloudflare.com/ajax/libs/jquery.smartmenus/1.1.0/css/sm-blue/sm-blue.min.css",
		// 	"https://cdnjs.cloudflare.com/ajax/libs/jquery.smartmenus/1.1.0/addons/bootstrap-4/jquery.smartmenus.bootstrap-4.css",
		// }
		// styleURLs = append(styleURLs, smartMenuStyleURLs...)
		smartMenuStyleURLs := []string{
			uncdn.JquerySmartMenusCss110(),
			uncdn.JquerySmartMenusCssSimpleTheme110(),
			uncdn.JquerySmartMenusCssBlueTheme110(),
			uncdn.JquerySmartMenusCssBootstrap4AddOn110(),
		}
		additionalStyles = append(additionalStyles, smartMenuStyleURLs...)
	}

	styleURLs = append(styleURLs, d.StyleURLs...)

	scriptURLs := []string{}

	if d.useSmartMenu {
		smartMenuScriptURLs := []string{
			"https://cdnjs.cloudflare.com/ajax/libs/jquery.smartmenus/1.1.0/jquery.smartmenus.min.js",
		}
		scriptURLs = append(scriptURLs, smartMenuScriptURLs...)
	}

	scriptURLs = append(scriptURLs, d.ScriptURLs...)
	faviconURL := d.FaviconURL
	if faviconURL == "" {
		faviconURL = favicon()
	}

	webpage := hb.NewWebpage()
	webpage.SetTitle(d.Title)
	webpage.AddStyleURLs(styleURLs)
	webpage.AddStyle("html,body{width:100%; height:100%;}")
	webpage.AddStyle(styles(append(d.Styles, additionalStyles...)))
	webpage.AddStyle(d.styles())
	webpage.AddScriptURLs(scriptURLs)
	webpage.AddScript(scripts(d.Scripts))
	webpage.AddScript(d.scripts())
	webpage.SetFavicon(faviconURL)
	if d.RedirectUrl != "" && d.RedirectTime != "" {
		webpage.Head.AddChild(hb.NewMeta().Attr("http-equiv", "refresh").Attr("content", d.RedirectTime+"; url = "+d.RedirectUrl))
	}

	webpage.AddChild(hb.NewHTML(d.layout()))

	return webpage.ToHTML()
}

func buildSubmenuItem(menuItem MenuItem, index int) *hb.Tag {
	title := menuItem.Title
	if title == "" {
		title = "n/a"
	}
	url := menuItem.URL
	if url == "" {
		url = "#"
	}
	icon := menuItem.Icon
	target := menuItem.Target
	if target == "" {
		target = "_self"
	}

	children := menuItem.Children
	hasChildren := len(children) > 0
	// menuId := "menu_" + utils.ToString(index)
	submenuId := "submenu_" + utils.ToString(index)
	if hasChildren {
		url = "#" + submenuId
	}

	a := hb.NewHyperlink().Class("nav-link px-0")
	if icon != "" {
		a.Child(hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(icon))
	} else {
		a.Child(hb.NewHTML(`
		    <svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" fill="currentColor" class="bi bi-caret-right-fill" viewBox="0 0 16 16">
		        <path d="m12.14 8.753-5.482 4.796c-.646.566-1.658.106-1.658-.753V3.204a1 1 0 0 1 1.659-.753l5.48 4.796a1 1 0 0 1 0 1.506z"/>
		    </svg>
		`))
	}
	a.Child(hb.NewSpan().Class("d-inline").HTML(title))
	a.Attr("href", url)
	if hasChildren {
		a.Attr("data-bs-toggle", "collapse")
	}
	li := hb.NewLI().Class("w-100").Child(a)
	return li
}

func buildMenuItem(menuItem MenuItem, index int) *hb.Tag {
	title := menuItem.Title
	if title == "" {
		title = "n/a"
	}
	url := menuItem.URL
	if url == "" {
		url = "#"
	}
	icon := menuItem.Icon
	children := menuItem.Children
	hasChildren := len(children) > 0
	submenuId := "submenu_" + utils.ToString(index)
	if hasChildren {
		url = "#" + submenuId
	}

	a := hb.NewHyperlink().Class("nav-link align-middle px-0")
	if icon != "" {
		a.AddChild(hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(icon))
	}
	a.HTML(title)
	a.Attr("href", url)
	if hasChildren {
		a.Attr("data-bs-toggle", "collapse")
	}
	if hasChildren {
		html := `<b class="caret">
			<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" fill="currentColor" class="bi bi-caret-down-fill" viewBox="0 0 16 16">
			<path d="M7.247 11.14 2.451 5.658C1.885 5.013 2.345 4 3.204 4h9.592a1 1 0 0 1 .753 1.659l-4.796 5.48a1 1 0 0 1-1.506 0z"/>
			</svg>
		</b>`
		a.Child(hb.NewHTML(html))
	}

	li := hb.NewLI().Attr("class", "nav-item").AddChild(a)

	if hasChildren {
		ul := hb.NewUL().
			ID(submenuId).
			Class("collapse hide nav flex-column ms-1").
			Attr("data-bs-parent", "#menu")
		for childIndex, childMenuItem := range children {
			childItem := buildSubmenuItem(childMenuItem, childIndex)
			ul.AddChild(childItem)
		}
		li.AddChild(ul)
	}

	return li
}

// func (d Dashboard) smartMenuBuild([]MenuItem) *hb.Tag {
// 	items := []*hb.Tag{}

// 	for index, menuItem := range d.menu {
// 		li := buildMenuItem(menuItem, index)
// 		items = append(items, li)
// 	}

// 	ul := hb.NewUL().ID("menu").AddChildren(items)

// 	return ul
// }

func (d Dashboard) smartMenu() string {
	html := `
	<nav role="navigation">
	<!-- Sample menu definition -->
	<ul id="main-menu" class="sm sm-blue sm-vertical">
	  <li><a href="http://www.smartmenus.org/">Home</a></li>
	  <li><a href="#">Long sub 1</a>
		<ul>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		</ul>
	  </li>
	  <li><a href="#">Long sub 2</a>
		<ul>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		</ul>
	  </li>
	  <li><a href="#">Sub 3</a>
		<ul>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">Dummy item</a></li>
		  <li><a href="#">more...</a>
			<ul>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">A pretty long text to test the default subMenusMaxWidth:20em setting for the sub menus</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			  <li><a href="#">Dummy item</a></li>
			</ul>
		  </li>
		</ul>
	  </li>
	</ul>
  </nav>
	`
	return html
}
func (d Dashboard) dashboardLayoutMenu() string {
	items := []*hb.Tag{}
	for index, menuItem := range d.menu {
		li := buildMenuItem(menuItem, index)
		items = append(items, li)
	}

	ul := hb.NewUL().
		ID("menu").
		Class("nav nav-pills flex-column mb-sm-auto mb-0 align-items-start").
		Children(items)

	return ul.ToHTML()
}

func (d Dashboard) center(content string) string {
	dropdownUser := hb.NewDiv().Class("dropdown").
		Child(hb.NewButton().Class("btn btn-secondary dropdown-toggle").
			Attr("type", "button").
			//Attr("id", "dropdownMenuButton1").
			Attr("data-bs-toggle", "dropdown").
			Attr("style", "background:#00A65A;").
			//Attr("aria-expanded", "false").
			HTML(d.user.FirstName + " " + d.user.LastName)).
		Child(hb.NewUL().Class("dropdown-menu").
			Child(hb.NewLI().
				//Attr("aria-labelledby", "dropdownMenuButton1").
				Child(hb.NewHyperlink().Class("dropdown-item").HTML("Logout").Attr("href", "/auth/logout"))))
	buttonMenu := hb.NewButton().Class("btn btn-outline-dark").
		Style("position:relative;width:30px;height:30px;border-radius:25px;padding:0px;").
		OnClick("toggleSideMenu()").
		AddChild(icons.Icon("bi-list", 16, 16, "").Style("margin-top:-4px;"))
	toolbar := hb.NewDiv().ID("Toolbar").Class("p-4").Style("background-color: #fff;z-index: 3;box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);transition: all .2s ease;").
		AddChild(buttonMenu).
		AddChild(hb.NewDiv().Class("float-end").AddChild(dropdownUser))

	contentHolder := hb.NewDiv().Class("shadow bg-white p-3 m-3").HTML(content)
	html := toolbar.ToHTML() + contentHolder.ToHTML()
	return html
}
func (d Dashboard) left() string {

	// 	personDropdownUseIfneeded := `
	// 	<div class="dropdown pb-4">
	// 	<a href="#" class="d-flex align-items-center text-white text-decoration-none dropdown-toggle" id="dropdownUser1" data-bs-toggle="dropdown" aria-expanded="false">
	// 		<img src="https://github.com/mdo.png" alt="hugenerd" width="30" height="30" class="rounded-circle">
	// 		<span class="d-none d-sm-inline mx-1">loser</span>
	// 	</a>
	// 	<ul class="dropdown-menu dropdown-menu-dark text-small shadow">
	// 		<li><a class="dropdown-item" href="#">New project...</a></li>
	// 		<li><a class="dropdown-item" href="#">Settings</a></li>
	// 		<li><a class="dropdown-item" href="#">Profile</a></li>
	// 		<li>
	// 			<hr class="dropdown-divider">
	// 		</li>
	// 		<li><a class="dropdown-item" href="#">Sign out</a></li>
	// 	</ul>
	// </div>
	// 	`

	menu := ""

	if d.useSmartMenu {
		menu = d.smartMenu()
	} else {
		menu = d.dashboardLayoutMenu()
	}

	var logo *hb.Tag
	logoURL := d.LogoURL
	if logoURL == "" {
		logoURL = utils.ImgPlaceholderURL(120, 80, "Logo")
		placeholderLogo := hb.NewImage().Attr("src", logoURL).Style("width:100%;margin:0px 10px 0px 0px;")
		adminDiv := hb.NewDiv().HTML("ADMIN PANEL").Style("font-size:12px;text-align: center;")
		logo = hb.NewDiv().Class("Logo").AddChild(placeholderLogo).AddChild(adminDiv)
	} else {
		logo = hb.NewImage().Attr("src", logoURL).Style("width:100%;margin:0px 10px 0px 0px;")
	}

	sideMenu := hb.NewDiv().ID("SideMenu").Class("p-4").Style("height:100%;width:200px;").
		AddChildren([]*hb.Tag{
			logo,
			hb.NewDiv().Class("Menu").HTML(menu),
		})
	return sideMenu.ToHTML()
}

// func bottom() {

// }

func top() string {
	html := ``
	return html
}

// func right() {

// }

func (d Dashboard) styles() string {
	css := `
html, body{
	height: 100%;
	background: #eee;
}
@media (min-width: 1200px) {
	.span12, .container {
		width: 1170px;
	}
}
#SideMenu{
	background: #343957;
}
#SideMenu a {
	color: #fff;
}
#SideMenu div.Logo {
	border:2px solid #999;
	color:#666;
	background: #eee;
}
#SideMenu div.Menu {
	margin: 30px 0px 30px 0px;
	padding: 10px 10px 10px 10px;
	background: #444;
	background-image: linear-gradient(to right, #444 , #555, #444);
	border-radius: 5px;
}
#Toolbar{
	background: #fff;
}
.well {
	min-height: 20px;
	padding: 19px;
	margin-bottom: 20px;
	background-color: #fafafa;
	border: 1px solid #e8e8e8;
	border-radius: 0;
	box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
}
	`
	return css
}

func (d Dashboard) scripts() string {

	toggleSidemenu := `function toggleSideMenu() { $("#SideMenu").toggle(); }`

	js := `
/**
* One of xs, sm, md, lg, xl. xxl
* @returns {String}
*/
function checkMode() {
	var width = $(window).width();
	if (width < 576) {
		return 'xs';
	} else if (width < 768) {
		return 'sm';
	} else if (width < 992) {
		return 'md';
	} else if (width < 1200) {
		return 'lg';
	} else if (width < 1400) {
		return 'xl';
	} else {
		return 'xxl';
	}
}
function onResized() {
	var mode = checkMode();
	if (mode == "xs") {
		$("#SideMenu").hide();
	}
}
$(window).on('resize', function () {
	onResized()
});

$(() => {
	onResized();
});
	`

	js += toggleSidemenu
	if d.useSmartMenu {
		js += `$(function () { $('#main-menu').smartmenus(); });`
	}

	return js
}

func favicon() string {
	favicon := "data:image/x-icon;base64,AAABAAEAEBAQAAAAAAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAzMzMAAAAmQBmZpkA////AJmZzAAzM5kAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAzMzMzMxQAA1YiIiIiUQADViIiIiZERANTIiIiJBVRRTMiIiJBNmJFNSIiJEMmZlRlIiJlYiJmVDUiImIiIiIUMzImMiIiJRFGUiZiImJkQEMzIiImFlEABDMyIiZiVAAEFTNiI2ZEAABBFTJmJRAAAABBRjQjQAAAAABFEBFACABwAAAAcAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAEAAMABAADAAQAA4AMAAPgDAAD+IwAA"
	return favicon
}
