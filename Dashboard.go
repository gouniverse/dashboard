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

type dashboardTemplateParams struct {
	Title        string
	Content      string
	Scripts      []string
	ScriptURLs   []string
	Styles       []string
	StyleURLs    []string
	RedirectUrl  string
	RedirectTime string
}

type dashboard struct {
	menu                      []MenuItem
	user                      User
	userMenu                  []MenuItem
	quickAccessMenu           []MenuItem
	navbarBackgroundColorMode string
	MenuType                  string
	Title                     string
	Content                   string
	FaviconURL                string
	LogoURL                   string
	Scripts                   []string
	ScriptURLs                []string
	Styles                    []string
	StyleURLs                 []string
	RedirectUrl               string
	RedirectTime              string
	ThemeHandlerUrl           string
	ThemeName                 string
	UncdnHandlerEndpoint      string
}

func (d *dashboard) layout() string {
	content := d.Content
	layout := hb.NewBorderLayout()
	layout.AddTop(hb.NewHTML(d.topNavigation()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	layout.AddCenter(hb.NewHTML(d.center(content)), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	return layout.ToHTML()
}

func (d *dashboard) SetUser(user User) *dashboard {
	d.user = user
	return d
}

func (d *dashboard) SetMenu(menuItems []MenuItem) *dashboard {
	d.menu = menuItems
	return d
}

func (d *dashboard) SetUserMenu(menuItems []MenuItem) *dashboard {
	d.userMenu = menuItems
	return d
}

// ToHTML returns the HTML representation of the dashboard.
//
// It does not take any parameters.
// It returns a string.
func (d *dashboard) ToHTML() string {
	styleURLs := []string{
		// Icons
		cdn.BootstrapIconsCss_1_10_2(),
	}
	// Theme
	if d.UncdnHandlerEndpoint != "" {
		styleURLs = append(styleURLs, uncdnThemeStyleURL(d.UncdnHandlerEndpoint, d.ThemeName))
	} else {
		styleURLs = append(styleURLs, cdnThemeStyleUrl(d.ThemeName))
	}
	// Other Style URLs
	styleURLs = append(styleURLs, d.StyleURLs...)

	scriptURLs := []string{}

	scriptURLs = append(scriptURLs, d.ScriptURLs...)
	faviconURL := d.FaviconURL
	if faviconURL == "" {
		faviconURL = favicon()
	}

	webpage := hb.NewWebpage()
	webpage.SetTitle(d.Title)
	webpage.AddStyleURLs(styleURLs)
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
		link.Child(hb.NewSpan().
			Class("icon").
			Style("margin-right: 5px;").
			HTML(icon))
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
			Data("bs-parent", "#dashboardMenu")
		for childIndex, childMenuItem := range children {
			childItem := buildSubmenuItem(childMenuItem, childIndex)
			ul.Child(childItem)
		}
		li.Child(ul)
	}

	return li
}

func (d *dashboard) dashboardLayoutMenu() string {
	items := []*hb.Tag{}
	for index, menuItem := range d.menu {
		li := buildMenuItem(menuItem, index)
		items = append(items, li)
	}

	ul := hb.NewUL().
		ID("dashboardMenu").
		Class("navbar-nav justify-content-end flex-grow-1 pe-3").
		Children(items)

	return ul.ToHTML()
}

func (d *dashboard) topNavigation() string {
	// isThemeDark := d.isThemeDark()
	isNavbarBackgroundDark := lo.Ternary(d.navbarBackgroundColorMode == "light", false, true)

	navbarTheme := lo.
		If(isNavbarBackgroundDark, "bg-dark text-bg-dark").
		Else("bg-light text-bg-light")

	buttonTheme := lo.
		If(isNavbarBackgroundDark, "btn-dark").
		Else("btn-light")

	iconStyle := "margin-top:-4px;margin-right:5px;"
	dropdownUser := hb.NewDiv().Class("dropdown").
		Children([]*hb.Tag{
			hb.NewButton().
				ID("ButtonUser").
				Class("btn "+buttonTheme+" dropdown-toggle").
				Style("background:none;border:0px;").
				Type(hb.TYPE_BUTTON).
				Data("bs-toggle", "dropdown").
				Children([]*hb.Tag{
					icons.Icon("bi-person", 24, 24, "").Style(iconStyle),
					hb.NewSpan().HTML(d.user.FirstName + " " + d.user.LastName).Style("margin-right:10px;"),
				}),
			hb.NewUL().
				Class("dropdown-menu").
				Children(lo.Map(d.userMenu, func(item MenuItem, _ int) *hb.Tag {
					target := lo.Ternary(item.Target == "", "_self", item.Target)
					url := lo.Ternary(item.URL == "", "#", item.URL)

					return hb.NewLI().Children([]*hb.Tag{
						hb.If(item.Title == "",
							hb.NewHR().
								Class("dropdown-divider"),
						),

						hb.If(item.Title != "",
							hb.NewHyperlink().
								Class("dropdown-item").
								ChildIf(item.Icon != "", hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(item.Icon)).
								HTML(item.Title).
								Href(url).
								Target(target),
						),
					})
				})),
		})

	dropdownQuickAccess := hb.NewDiv().Class("dropdown").
		Children([]*hb.Tag{
			hb.NewButton().
				ID("ButtonUser").
				Class("btn "+buttonTheme+" dropdown-toggle").
				Style("background:none;border:0px;").
				Type(hb.TYPE_BUTTON).
				Data("bs-toggle", "dropdown").
				Children([]*hb.Tag{
					icons.Icon("bi-microsoft", 24, 24, "").Style("margin-top:-4px;margin-right:8px;"),
					hb.NewSpan().HTML("Quick Access").Style("margin-right:10px;"),
				}),
			hb.NewUL().
				Class("dropdown-menu").
				Children(lo.Map(d.quickAccessMenu, func(item MenuItem, _ int) *hb.Tag {
					target := lo.Ternary(item.Target == "", "_self", item.Target)
					url := lo.Ternary(item.URL == "", "#", item.URL)

					return hb.NewLI().Children([]*hb.Tag{
						hb.If(item.Title == "",
							hb.NewHR().
								Class("dropdown-divider"),
						),

						hb.If(item.Title != "",
							hb.NewHyperlink().
								Class("dropdown-item").
								ChildIf(item.Icon != "", hb.NewSpan().Class("icon").Style("margin-right: 5px;").HTML(item.Icon)).
								HTML(item.Title).
								Href(url).
								Target(target),
						),
					})
				})),
		})

	buttonMenuToggle := hb.NewButton().
		Class("btn "+buttonTheme).
		Style("background: none;").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModaldashboardMenu").
		Children([]*hb.Tag{
			icons.Icon("bi-list", 24, 24, "").Style(iconStyle),
			hb.NewSpan().HTML("Menu"),
		})

	buttonOffcanvasToggle := hb.NewButton().
		Class("btn "+buttonTheme).
		Style("background: none;").
		Data("bs-toggle", "offcanvas").
		Data("bs-target", "#OffcanvasMenu").
		Children([]*hb.Tag{
			icons.Icon("bi-list", 24, 24, "").Style(iconStyle),
			hb.NewSpan().HTML("Menu"),
		})

	menu := buttonOffcanvasToggle
	if d.MenuType == MENU_TYPE_MODAL {
		menu = buttonMenuToggle
	}

	toolbar := hb.NewNav().
		ID("Toolbar").
		Class("navbar " + navbarTheme).
		Style("z-index: 3;box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);transition: all .2s ease;padding-left: 20px;padding-right: 20px; display:block;").
		Children([]*hb.Tag{
			// Main Menu
			menu,
			// User Menu
			hb.If(d.user.FirstName != "" && d.user.LastName != "",
				hb.NewDiv().Class("float-end").
					Style("margin-left:10px;").
					Child(dropdownUser),
			),
			// Theme Switcher
			hb.If(d.ThemeHandlerUrl != "",
				hb.NewDiv().Class("float-end").
					Style("margin-left:10px;").
					Child(d.themeButton()),
			),
			// Quick Menu (if provided)
			hb.If(len(d.quickAccessMenu) > 0, hb.NewDiv().
				Class("float-end").
				Style("margin-left:10px;").
				Child(dropdownQuickAccess)),
		})

	return toolbar.ToHTML()
}

func (d *dashboard) center(content string) string {
	contentHolder := hb.NewDiv().Class("shadow p-3 m-3").HTML(content)
	html := contentHolder.ToHTML()
	return html
}

func (d *dashboard) menuOffcanvas() *hb.Tag {
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

func (d *dashboard) menuModal() *hb.Tag {
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
		ID("ModaldashboardMenu").
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

// func (d *dashboard) left() string {
// 	menu := d.dashboardLayoutMenu()

// 	var logo *hb.Tag
// 	logoURL := d.LogoURL
// 	if logoURL == "" {
// 		logoURL = utils.ImgPlaceholderURL(120, 80, "Logo")
// 		placeholderLogo := hb.NewImage().
// 			Src(logoURL).
// 			Style("width:100%;margin:0px 10px 0px 0px;")
// 		adminDiv := hb.NewDiv().
// 			HTML("ADMIN PANEL").
// 			Style("font-size:12px;text-align: center;")
// 		logo = hb.NewDiv().Class("Logo").Child(placeholderLogo).Child(adminDiv)
// 	} else {
// 		logo = hb.NewImage().Attr("src", logoURL).Style("width:100%;margin:0px 10px 0px 0px;")
// 	}

// 	sideMenu := hb.NewDiv().ID("SideMenu").Class("p-4").Style("height:100%;width:200px;").
// 		AddChildren([]*hb.Tag{
// 			logo,
// 			hb.NewDiv().Class("Menu").HTML(menu),
// 		})
// 	return sideMenu.ToHTML()
// }

func (d *dashboard) styles() string {
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

	// #ModaldashboardMenu .nav-item {
	// 	border: 1px solid #999;
	// 	background: #eee;
	// 	width: 100%;
	// 	margin: 10px 0px;
	// 	border-radius: 10px;
	// 	padding: 10px;
	// }

	// #ModaldashboardMenu .nav-item:hover {
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

// scripts returns the JavaScript code for the dashboard.
//
// No parameters.
// Returns a string.
func (d *dashboard) scripts() string {
	js := ``
	return js
}

// favicon returns the data URI for a website favicon.
//
// No parameters.
// Returns a string.
func favicon() string {
	favicon := "data:image/x-icon;base64,AAABAAEAEBAQAAAAAAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAzMzMAAAAmQBmZpkA////AJmZzAAzM5kAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAzMzMzMxQAA1YiIiIiUQADViIiIiZERANTIiIiJBVRRTMiIiJBNmJFNSIiJEMmZlRlIiJlYiJmVDUiImIiIiIUMzImMiIiJRFGUiZiImJkQEMzIiImFlEABDMyIiZiVAAEFTNiI2ZEAABBFTJmJRAAAABBRjQjQAAAAABFEBFACABwAAAAcAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAEAAMABAADAAQAA4AMAAPgDAAD+IwAA"
	return favicon
}

// isThemeDark checks if the theme of the dashboard is dark.
//
// It does so by checking if the dashboard's theme name is contained
// in the list of dark themes.
//
// Returns a boolean indicating whether the theme is dark.
func (d *dashboard) isThemeDark() bool {
	isDark := lo.Contains(lo.Keys(themesDark), d.ThemeName)
	return isDark
}

// themeButton generates a dropdown menu with light and dark themes.
//
// It checks if the current theme is dark and creates dropdown items for both light and dark themes.
// The dropdown items are created dynamically based on the themesLight and themesDark maps.
// The function returns a *hb.Tag that represents the generated dropdown menu.
func (d *dashboard) themeButton() *hb.Tag {
	isDark := d.isThemeDark()

	// Light Themes
	lightDropdownItems := lo.Map(lo.Keys(themesLight), func(theme string, index int) *hb.Tag {
		name := themesLight[theme]
		active := lo.Ternary(d.ThemeName == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.ThemeHandlerUrl, "?"), d.ThemeHandlerUrl+"&theme="+theme, d.ThemeHandlerUrl+"?theme="+theme)

		return hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				HTML("(Light) "+name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	// Dark Themes
	darkDropdownItems := lo.Map(lo.Keys(themesDark), func(theme string, index int) *hb.Tag {
		name := themesDark[theme]
		active := lo.Ternary(d.ThemeName == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.ThemeHandlerUrl, "?"), d.ThemeHandlerUrl+"&theme="+theme, d.ThemeHandlerUrl+"?theme="+theme)

		return hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				HTML("(Dark) "+name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	return hb.NewDiv().Class("dropdown").Children([]*hb.Tag{
		bs.Button().
			ID("buttonTheme").
			Class("dropdown-toggle").
			Data("bs-toggle", "dropdown").
			Children([]*hb.Tag{
				lo.Ternary(isDark, icons.Icon("bi-brightness-high-fill", 16, 16, "white"), icons.Icon("bi-brightness-high-fill", 16, 16, "black")),
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
