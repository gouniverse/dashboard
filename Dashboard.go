package dashboard

import (
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

const MENU_TYPE_MODAL = "modal"
const MENU_TYPE_OFFCANVAS = "offcanvas"

type Dashboard struct {
	menu                 []MenuItem
	user                 User
	MenuType             string
	Title                string
	Content              string
	FaviconURL           string
	LogoURL              string
	Scripts              []string
	ScriptURLs           []string
	Styles               []string
	StyleURLs            []string
	RedirectUrl          string
	RedirectTime         string
	ThemeHandlerUrl      string
	ThemeName            string
	UncdnHandlerEndpoint string
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
	layout.AddTop(hb.NewHTML(d.top()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	// layout.AddLeft(hb.NewHTML(d.left()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	// layout.AddBottom(hb.NewHTML("BOTTOM"), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
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

	styleURLs = append(styleURLs, d.StyleURLs...)

	scriptURLs := []string{}

	scriptURLs = append(scriptURLs, d.ScriptURLs...)
	faviconURL := d.FaviconURL
	if faviconURL == "" {
		faviconURL = favicon()
	}

	webpage := hb.NewWebpage()
	webpage.SetTitle(d.Title)
	webpage.AddStyleURLs(d.themeStyleURLs(d.ThemeName))
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

	menu := d.menuOffcanvas().ToHTML()
	if d.MenuType == MENU_TYPE_MODAL {
		menu += d.menuModal().ToHTML()
	}

	webpage.AddChild(hb.NewHTML(d.layout() + menu))

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

	link := hb.NewHyperlink().Class("nav-link px-0")
	if icon != "" {
		link.Child(hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(icon))
	} else {
		link.Child(hb.NewHTML(`
		    <svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" fill="currentColor" class="bi bi-caret-right-fill" viewBox="0 0 16 16">
		        <path d="m12.14 8.753-5.482 4.796c-.646.566-1.658.106-1.658-.753V3.204a1 1 0 0 1 1.659-.753l5.48 4.796a1 1 0 0 1 0 1.506z"/>
		    </svg>
		`))
	}
	link.Child(hb.NewSpan().Class("d-inline").HTML(title))
	link.Href(url)
	if hasChildren {
		link.Data("bs-toggle", "collapse")
	}

	return hb.NewLI().
		Class("w-100").
		Child(link)
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

	link := hb.NewHyperlink().Class("nav-link align-middle px-0")
	if icon != "" {
		link.Child(hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(icon))
	}
	link.HTML(title)
	link.Attr("href", url)
	if hasChildren {
		link.Data("bs-toggle", "collapse")
	}
	if hasChildren {
		html := `<b class="caret">
			<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" fill="currentColor" class="bi bi-caret-down-fill" viewBox="0 0 16 16">
			<path d="M7.247 11.14 2.451 5.658C1.885 5.013 2.345 4 3.204 4h9.592a1 1 0 0 1 .753 1.659l-4.796 5.48a1 1 0 0 1-1.506 0z"/>
			</svg>
		</b>`
		link.Child(hb.NewHTML(html))
	}

	li := hb.NewLI().Class("nav-item").Child(link)

	if hasChildren {
		ul := hb.NewUL().
			ID(submenuId).
			Class("collapse hide nav flex-column ms-1").
			Data("bs-parent", "#DashboardMenu")
		for childIndex, childMenuItem := range children {
			childItem := buildSubmenuItem(childMenuItem, childIndex)
			ul.Child(childItem)
		}
		li.Child(ul)
	}

	return li
}

func (d Dashboard) dashboardLayoutMenu() string {
	items := []*hb.Tag{}
	for index, menuItem := range d.menu {
		li := buildMenuItem(menuItem, index)
		items = append(items, li)
	}

	ul := hb.NewUL().
		ID("DashboardMenu").
		// Class("nav nav-pills flex-column mb-sm-auto mb-0 align-items-start").
		Class("navbar-nav justify-content-end flex-grow-1 pe-3").
		Children(items)

	return ul.ToHTML()
}

func (d Dashboard) top() string {
	dropdownUser := hb.NewDiv().Class("dropdown").
		Children([]*hb.Tag{
			hb.NewButton().
				ID("ButtonUser").
				Class("btn btn-secondary dropdown-toggle").
				Attr("type", "button").
				Data("bs-toggle", "dropdown").
				HTML(d.user.FirstName + " " + d.user.LastName),
			hb.NewUL().Class("dropdown-menu").
				Children([]*hb.Tag{
					hb.NewLI().Children([]*hb.Tag{
						hb.NewHyperlink().
							Class("dropdown-item").
							HTML("Logout").
							Href("/auth/logout"),
					}),
				}),
		})

	buttonMenuToggle := hb.NewButton().Class("btn btn-secondary").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModalDashboardMenu").
		Children([]*hb.Tag{
			icons.Icon("bi-list", 16, 16, "").Style("margin-top:-4px;margin-right:5px;"),
			hb.NewSpan().HTML("Menu"),
		})

	buttonOffcanvasToggle := hb.NewButton().
		Class("btn btn-secondary"). // outline-dark
		Data("bs-toggle", "offcanvas").
		Data("bs-target", "#OffcanvasMenu").
		Children([]*hb.Tag{
			icons.Icon("bi-list", 16, 16, "").Style("margin-top:-4px;margin-right:5px;"),
			hb.NewSpan().HTML("Menu"),
		})

	menu := buttonOffcanvasToggle
	if d.MenuType == MENU_TYPE_MODAL {
		menu = buttonMenuToggle
	}

	toolbar := hb.NewNav().
		ID("Toolbar").
		Class("navbar navbar-dark bg-dark").
		Style("background-color: #fff;z-index: 3;box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);transition: all .2s ease;padding-left: 20px;padding-right: 20px;").
		Children([]*hb.Tag{
			menu,
			hb.NewDiv().Children([]*hb.Tag{
				hb.NewDiv().Class("float-end").
					Style("margin-left:10px;").
					Child(dropdownUser),
				hb.NewDiv().Class("float-end").
					Style("margin-left:10px;").
					Child(d.themeButton(d.ThemeName)),
			}),
		})

	return toolbar.ToHTML()
}

func (d Dashboard) center(content string) string {
	contentHolder := hb.NewDiv().Class("shadow p-3 m-3").HTML(content)
	html := contentHolder.ToHTML()
	return html
}

func (d Dashboard) menuOffcanvas() *hb.Tag {
	offcanvasMenu := hb.NewDiv().
		ID("OffcanvasMenu").
		Class("offcanvas offcanvas-start text-bg-dark").
		Attr("tabindex", "-1").
		Children([]*hb.Tag{
			hb.NewDiv().Class("offcanvas-header").
				Children([]*hb.Tag{
					hb.NewHeading5().
						Class("offcanvas-title").
						HTML("Menu"),
					hb.NewButton().
						Class("btn-close btn-close-white").
						Attr("type", "button").
						Data("bs-dismiss", "offcanvas").
						Attr("aria-label", "Close"),
				}),
			hb.NewDiv().Class("offcanvas-body").
				Children([]*hb.Tag{
					hb.NewHTML(d.dashboardLayoutMenu()),
				}),
		})

	return offcanvasMenu
}

func (d Dashboard) menuModal() *hb.Tag {
	modalHeader := hb.NewDiv().Class("modal-header").
		Children([]*hb.Tag{
			hb.NewHeading5().HTML("Menu").Class("modal-title"),
			hb.NewButton().Attrs(map[string]string{
				"type":            "button",
				"class":           "btn-close",
				"data-bs-dismiss": "modal",
				"aria-label":      "Close",
			}),
		})

	modalBody := hb.NewDiv().Class("modal-body").Children([]*hb.Tag{
		hb.NewHTML(d.dashboardLayoutMenu()),
	})

	modalFooter := hb.NewDiv().Class("modal-footer").Children([]*hb.Tag{
		hb.NewButton().
			HTML("Close").
			Class("btn btn-secondary w-100").
			Data("bs-dismiss", "modal"),
	})

	modal := hb.NewDiv().
		ID("ModalDashboardMenu").
		Class("modal fade").
		Children([]*hb.Tag{
			hb.NewDiv().Class("modal-dialog modal-lg").
				Children([]*hb.Tag{
					hb.NewDiv().Class("modal-content").
						Children([]*hb.Tag{
							modalHeader,
							modalBody,
							modalFooter,
						}),
				}),
		})

	return modal
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

	menu = d.dashboardLayoutMenu()

	var logo *hb.Tag
	logoURL := d.LogoURL
	if logoURL == "" {
		logoURL = utils.ImgPlaceholderURL(120, 80, "Logo")
		placeholderLogo := hb.NewImage().Attr("src", logoURL).Style("width:100%;margin:0px 10px 0px 0px;")
		adminDiv := hb.NewDiv().HTML("ADMIN PANEL").Style("font-size:12px;text-align: center;")
		logo = hb.NewDiv().Class("Logo").Child(placeholderLogo).Child(adminDiv)
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

func (d Dashboard) styles() string {
	// @media (min-width: 1200px) {
	// 	.span12, .container {
	// 		width: 1170px;
	// 	}
	// }
	// #SideMenu{
	// 	background: #343957;
	// }
	// #SideMenu a {
	// 	color: #fff;
	// }
	// #SideMenu div.Logo {
	// 	border:2px solid #999;
	// 	color:#666;
	// 	background: #eee;
	// }
	// #SideMenu div.Menu {
	// 	margin: 30px 0px 30px 0px;
	// 	padding: 10px 10px 10px 10px;
	// 	background: #444;
	// 	background-image: linear-gradient(to right, #444 , #555, #444);
	// 	border-radius: 5px;
	// }
	// #Toolbar{
	// 	background: #fff;
	// }

	// #ModalDashboardMenu .nav-item {
	// 	border: 1px solid #999;
	// 	background: #eee;
	// 	width: 100%;
	// 	margin: 10px 0px;
	// 	border-radius: 10px;
	// 	padding: 10px;
	// }

	// #ModalDashboardMenu .nav-item:hover {
	// 	background: cornsilk;
	// }

	// .well {
	// 	min-height: 20px;
	// 	padding: 19px;
	// 	margin-bottom: 20px;
	// 	background-color: #fafafa;
	// 	border: 1px solid #e8e8e8;
	// 	border-radius: 0;
	// 	box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
	// }
	// html, body{
	// 	height: 100%;
	// 	background: #eee;
	// }
	css := `
html, body{
	height: 100%;
}
	`
	return css
}

func (d Dashboard) scripts() string {
	js := ``
	return js
}

func favicon() string {
	favicon := "data:image/x-icon;base64,AAABAAEAEBAQAAAAAAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAzMzMAAAAmQBmZpkA////AJmZzAAzM5kAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAzMzMzMxQAA1YiIiIiUQADViIiIiZERANTIiIiJBVRRTMiIiJBNmJFNSIiJEMmZlRlIiJlYiJmVDUiImIiIiIUMzImMiIiJRFGUiZiImJkQEMzIiImFlEABDMyIiZiVAAEFTNiI2ZEAABBFTJmJRAAAABBRjQjQAAAAABFEBFACABwAAAAcAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAEAAMABAADAAQAA4AMAAPgDAAD+IwAA"
	return favicon
}

func (d Dashboard) themeButton(themeName string) *hb.Tag {
	isDark := lo.Contains(lo.Keys(themesDark), themeName)
	lightDropdownItems := []*hb.Tag{}
	lo.ForEach(lo.Keys(themesLight), func(theme string, index int) {
		name := themesLight[theme]
		active := lo.Ternary(themeName == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.ThemeHandlerUrl, "?"), d.ThemeHandlerUrl+"&theme="+theme, d.ThemeHandlerUrl+"?theme="+theme)

		lightDropdownItems = append(lightDropdownItems, hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				HTML("(Light) "+name).
				Href(url).
				Attr("ref", "nofollow"),
		}))
	})

	darkDropdownItems := []*hb.Tag{}
	lo.ForEach(lo.Keys(themesDark), func(theme string, index int) {
		name := themesDark[theme]
		active := lo.Ternary(themeName == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.ThemeHandlerUrl, "?"), d.ThemeHandlerUrl+"&theme="+theme, d.ThemeHandlerUrl+"?theme="+theme)

		darkDropdownItems = append(darkDropdownItems, hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				HTML("(Dark) "+name).
				Href(url).
				Attr("ref", "nofollow"),
		}))
	})
	return hb.NewDiv().Class("dropdown").Children([]*hb.Tag{
		bs.Button().
			ID("buttonTheme").
			Class("btn-secondary dropdown-toggle").
			Data("bs-toggle", "dropdown").
			Children([]*hb.Tag{
				lo.Ternary(isDark, icons.Icon("bi-brightness-high-fill", 16, 16, "#fff"), icons.Icon("bi-brightness-high-fill", 16, 16, "#333")),
			}),
		hb.NewUL().Class("dropdown-menu dropdown-menu-dark").
			Children(lightDropdownItems).
			Children([]*hb.Tag{
				hb.NewLI().Children([]*hb.Tag{
					hb.NewHR().Class("dropdown-divider"),
				}),
			}).
			Children(darkDropdownItems),
	})
}

func (d Dashboard) themeStyleURLs(theme string) []string {
	themeStyle := lo.
		If(theme == "cerulean", "bs523"+theme).
		ElseIf(theme == "cosmo", "bs523"+theme).
		ElseIf(theme == "cyborg", "bs523"+theme).
		ElseIf(theme == "darkly", "bs523"+theme).
		ElseIf(theme == "flatly", "bs523"+theme).
		ElseIf(theme == "journal", "bs523"+theme).
		ElseIf(theme == "litera", "bs523"+theme).
		ElseIf(theme == "lumen", "bs523"+theme).
		ElseIf(theme == "lux", "bs523"+theme).
		ElseIf(theme == "materia", "bs523"+theme).
		ElseIf(theme == "minty", "bs523"+theme).
		ElseIf(theme == "morph", "bs523"+theme).
		ElseIf(theme == "pulse", "bs523"+theme).
		ElseIf(theme == "quartz", "bs523"+theme).
		ElseIf(theme == "sandstone", "bs523"+theme).
		ElseIf(theme == "simplex", "bs523"+theme).
		ElseIf(theme == "sketchy", "bs523"+theme).
		ElseIf(theme == "slate", "bs523"+theme).
		ElseIf(theme == "solar", "bs523"+theme).
		ElseIf(theme == "spacelab", "bs523"+theme).
		ElseIf(theme == "superhero", "bs523"+theme).
		ElseIf(theme == "united", "bs523"+theme).
		ElseIf(theme == "vapor", "bs523"+theme).
		ElseIf(theme == "yeti", "bs523"+theme).
		ElseIf(theme == "zephyr", "bs523"+theme).
		Else("bs523")

	themeURL := d.UncdnHandlerEndpoint + "/" + themeStyle + ".css"

	urls := []string{
		themeURL,
		// cdn.BootstrapIconsCss191(),
	}
	return urls
}
