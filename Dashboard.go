package dashboard

import (
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/uncdn"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

const MENU_TYPE_MODAL = "modal"
const MENU_TYPE_OFFCANVAS = "offcanvas"

type Dashboard struct {
	menu     []MenuItem
	user     User
	userMenu []MenuItem

	// Optional. Menu for quick access to various pages
	quickAccessMenu []MenuItem

	// Optional. The background color for the navbar, light or dark (default)
	navbarBackgroundColorMode string

	// Optional. The URL of the login page to use (if user is not provided)
	loginURL string

	// Optional. The URL of the register page to use (if user is not provided)
	registerURL string

	// Optional. The type of the main menu (see MENU_TYPE_* constants)
	menuType string

	// Optional. The web page title of the dashboard
	title string

	// Optional. The content of the dashboard
	content string

	// Optional. The URL of the favicon (base64 encoded can be used, default will be used otherwise)
	faviconURL      string
	logoImageURL    string
	logoRedirectURL string
	scripts         []string
	scriptURLs      []string
	styles          []string
	styleURLs       []string
	redirectUrl     string
	redirectTime    string

	// Optional. The URL of the theme switcher endpoint to use
	themeHandlerUrl string

	// Optional. The theme name to be activated on the dashboard (default will be used otherwise)
	theme string

	// Optional. The theme names to be visible in the theme switcher, the key is the theme, the value is the theme name (can be customized, default will be used otherwise)
	themesRestrict map[string]string

	// Optional. The URL of the UNCDN hadler endpoint to use
	uncdnHandlerEndpoint string
}

func (d *Dashboard) layout() string {
	content := d.content
	layout := hb.NewBorderLayout()
	layout.AddTop(hb.NewHTML(d.topNavigation()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	layout.AddCenter(hb.NewHTML(d.center(content)), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
	return layout.ToHTML()
}

func (d *Dashboard) SetUser(user User) *Dashboard {
	d.user = user
	return d
}

func (d *Dashboard) SetMenu(menuItems []MenuItem) *Dashboard {
	d.menu = menuItems
	return d
}

func (d *Dashboard) SetUserMenu(menuItems []MenuItem) *Dashboard {
	d.userMenu = menuItems
	return d
}

// ToHTML returns the HTML representation of the Dashboard.
//
// It does not take any parameters.
// It returns a string.
func (d *Dashboard) ToHTML() string {
	styleURLs := []string{
		// Icons
		cdn.BootstrapIconsCss_1_10_2(),
	}

	// Bootstrap Css
	if d.uncdnHandlerEndpoint != "" {
		styleURLs = append(styleURLs, uncdnThemeStyleURL(d.uncdnHandlerEndpoint, d.theme))
	} else {
		styleURLs = append(styleURLs, cdnThemeStyleUrl(d.theme))
	}

	scriptURLs := []string{}

	// Bootstrap JS
	if d.uncdnHandlerEndpoint != "" {
		scriptURLs = append(scriptURLs, uncdn.BootstrapJs523())
	} else {
		scriptURLs = append(scriptURLs, cdn.BootstrapJs_5_3_0())
	}

	faviconURL := d.faviconURL
	if faviconURL == "" {
		faviconURL = favicon()
	}

	webpage := hb.NewWebpage()
	webpage.SetTitle(d.title)

	// Required Style URLs
	webpage.AddStyleURLs(styleURLs)
	// User Style URLs
	webpage.AddStyleURLs(d.styleURLs)
	// Dashboard Styles
	webpage.AddStyle(d.dashboardStyle())
	// User Styles
	webpage.AddStyles(d.styles)

	// Required Script URLs
	webpage.AddScriptURLs(scriptURLs)
	// User Script URLs
	webpage.AddScriptURLs(d.scriptURLs)

	// Dashboard Scripts
	webpage.AddScript(d.dashboardScript())
	// User Scripts
	webpage.AddScripts(d.scripts)

	// webpage.AddScript(scripts(d.scripts))
	webpage.SetFavicon(faviconURL)
	if d.redirectUrl != "" && d.redirectTime != "" {
		webpage.Meta(hb.NewMeta().
			Attr("http-equiv", "refresh").
			Attr("content", d.redirectTime+"; url = "+d.redirectUrl))
	}

	menu := d.menuOffcanvas().ToHTML()
	if d.menuType == MENU_TYPE_MODAL {
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
			Data("bs-parent", "#DashboardMenu")
		for childIndex, childMenuItem := range children {
			childItem := buildSubmenuItem(childMenuItem, childIndex)
			ul.Child(childItem)
		}
		li.Child(ul)
	}

	return li
}

func (d *Dashboard) DashboardLayoutMenu() string {
	items := []*hb.Tag{}
	for index, menuItem := range d.menu {
		li := buildMenuItem(menuItem, index)
		items = append(items, li)
	}

	ul := hb.NewUL().
		ID("DashboardMenu").
		Class("navbar-nav justify-content-end flex-grow-1 pe-3").
		Children(items)

	return ul.ToHTML()
}

// topNavigation returns the HTML code for the top navigation toolbar in the Dashboard.
//
// No parameters.
// Returns a string.
func (d *Dashboard) topNavigation() string {
	isNavbarBackgroundDark := lo.Ternary(d.navbarBackgroundColorMode == "light", false, true)

	hasLogo := lo.Ternary(d.logoImageURL != "", true, false)
	logoRedirectURL := lo.Ternary(d.logoRedirectURL != "", d.logoRedirectURL, "#")

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
		Data("bs-target", "#ModalDashboardMenu").
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

	mainMenu := buttonOffcanvasToggle
	if d.menuType == MENU_TYPE_MODAL {
		mainMenu = buttonMenuToggle
	}

	logo := hb.NewHyperlink().
		Href(logoRedirectURL).
		Class("navbar-brand").
		Child(
			hb.NewImage().
				Src(d.logoImageURL).
				Style("max-height:35px;"),
		)

	toolbar := hb.NewNav().
		ID("Toolbar").
		Class("navbar "+navbarTheme).
		Style("z-index: 3;box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);transition: all .2s ease;padding-left: 20px;padding-right: 20px; display:block;").
		ChildIf(hasLogo, logo).
		Children([]*hb.Tag{
			mainMenu,
			// User Menu
			hb.If(!lo.IsEmpty(d.user) && (d.user.FirstName != "" || d.user.LastName != ""),
				hb.NewDiv().Class("float-end").
					Style("margin-left:10px;").
					Child(dropdownUser),
			),

			// Register Link
			hb.If(lo.IsEmpty(d.user) && d.registerURL != "",
				hb.NewHyperlink().
					HTML("Register").
					Href(d.registerURL).
					Class("btn "+buttonTheme+" float-end").
					Style("margin-left:10px;"),
			),

			// Login Link
			hb.If(lo.IsEmpty(d.user) && d.loginURL != "",
				hb.NewHyperlink().
					HTML("Login").
					Href(d.loginURL).
					Class("btn "+buttonTheme+" float-end").
					Style("margin-left:10px;"),
			),

			// Theme Switcher
			hb.If(d.themeHandlerUrl != "",
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

func (d *Dashboard) center(content string) string {
	contentHolder := hb.NewDiv().Class("shadow p-3 m-3").HTML(content)
	html := contentHolder.ToHTML()
	return html
}

func (d *Dashboard) menuOffcanvas() *hb.Tag {
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
					hb.NewHTML(d.DashboardLayoutMenu()),
				}),
		})

	return offcanvasMenu
}

func (d *Dashboard) menuModal() *hb.Tag {
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
		hb.NewHTML(d.DashboardLayoutMenu()),
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

// func (d *Dashboard) left() string {
// 	menu := d.DashboardLayoutMenu()

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

func (d *Dashboard) dashboardStyle() string {
	fullHeightSupport := `html, body{ height: 100%; }`

	css := fullHeightSupport
	return css
}

// scripts returns the JavaScript code for the Dashboard.
//
// No parameters.
// Returns a string.
func (d *Dashboard) dashboardScript() string {
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

// isThemeDark checks if the theme of the Dashboard is dark.
//
// It does so by checking if the Dashboard's theme name is contained
// in the list of dark themes.
//
// Returns a boolean indicating whether the theme is dark.
func (d *Dashboard) isThemeDark() bool {
	isDark := lo.Contains(lo.Keys(themesDark), d.theme)
	return isDark
}

// themeButton generates a dropdown menu with light and dark themes.
//
// It checks if the current theme is dark and creates dropdown items for both light and dark themes.
// The dropdown items are created dynamically based on the themesLight and themesDark maps.
// The function returns a *hb.Tag that represents the generated dropdown menu.
func (d *Dashboard) themeButton() *hb.Tag {
	isDark := d.isThemeDark()
	isNavbarBackgroundDark := lo.Ternary(d.navbarBackgroundColorMode == "light", false, true)

	buttonTheme := lo.
		If(isNavbarBackgroundDark, "btn-dark").
		Else("btn-light")

	// Light Themes
	lightDropdownItems := lo.Map(lo.Keys(themesLight), func(theme string, index int) *hb.Tag {
		name := themesLight[theme]
		active := lo.Ternary(d.theme == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.themeHandlerUrl, "?"), d.themeHandlerUrl+"&theme="+theme, d.themeHandlerUrl+"?theme="+theme)

		if len(d.themesRestrict) > 0 {
			if customName, exists := d.themesRestrict[theme]; exists {
				name = customName
			} else {
				return nil
			}
		}

		return hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				Child(hb.NewI().Class("bi bi-sun").Style("margin-right:5px;")).
				HTML(name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	// Dark Themes
	darkDropdownItems := lo.Map(lo.Keys(themesDark), func(theme string, index int) *hb.Tag {
		name := themesDark[theme]
		active := lo.Ternary(d.theme == theme, " active", "")
		url := lo.Ternary(strings.Contains(d.themeHandlerUrl, "?"), d.themeHandlerUrl+"&theme="+theme, d.themeHandlerUrl+"?theme="+theme)

		if len(d.themesRestrict) > 0 {
			if customName, exists := d.themesRestrict[theme]; exists {
				name = customName
			} else {
				return nil
			}
		}

		return hb.NewLI().Children([]*hb.Tag{
			hb.NewHyperlink().
				Class("dropdown-item"+active).
				Child(hb.NewI().Class("bi bi-moon-stars-fill").Style("margin-right:5px;")).
				HTML(name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	return hb.NewDiv().Class("dropdown").Children([]*hb.Tag{
		bs.Button().
			ID("buttonTheme").
			Class(buttonTheme+" dropdown-toggle").
			Data("bs-toggle", "dropdown").
			Children([]*hb.Tag{
				lo.Ternary(isDark, hb.NewI().Class("bi bi-sun"), hb.NewI().Class("bi bi-moon-stars-fill")),
			}),
		hb.NewUL().Class(buttonTheme+" dropdown-menu dropdown-menu-dark").
			Children(lightDropdownItems).
			ChildIf(
				len(lo.Filter(darkDropdownItems, func(item *hb.Tag, _ int) bool { return item != nil })) > 0 && len(lo.Filter(lightDropdownItems, func(item *hb.Tag, _ int) bool { return item != nil })) > 0,
				hb.NewLI().Children([]*hb.Tag{
					hb.NewHR().Class("dropdown-divider"),
				}),
			).
			Children(darkDropdownItems),
	})
}
