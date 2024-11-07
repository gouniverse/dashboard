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

	// Optional. The background color for the navbar (default dark)
	navbarBackgroundColor string

	// Optional. The text color for the navbar (default light)
	navbarTextColor string

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
	faviconURL string

	// Optional. The URL of the logo (base64 encoded can be used, default will be used otherwise)
	logoImageURL string
	// Optional. Raw HTML of the logo, if set will be used instead of logoImageURL
	logoRawHtml     string
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

	// Optional. The URL of the UNCDN handler endpoint to use
	uncdnHandlerEndpoint string
}

func (d *Dashboard) layout() string {
	content := d.content
	layout := hb.NewBorderLayout()
	layout.AddTop(hb.Raw(d.topNavigation()), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_MIDDLE)
	layout.AddCenter(hb.Raw(d.center(content)), hb.BORDER_LAYOUT_ALIGN_LEFT, hb.BORDER_LAYOUT_ALIGN_TOP)
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
		cdn.BootstrapIconsCss_1_11_3(),
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
		scriptURLs = append(scriptURLs, cdn.BootstrapJs_5_3_3())
	}

	faviconURL := d.faviconURL
	if faviconURL == "" {
		faviconURL = favicon()
	}

	webpage := hb.Webpage()
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
		webpage.Meta(hb.Meta().
			Attr("http-equiv", "refresh").
			Attr("content", d.redirectTime+"; url = "+d.redirectUrl))
	}

	menu := d.menuOffcanvas().ToHTML()

	if d.menuType == MENU_TYPE_MODAL {
		menu += d.menuModal().ToHTML()
	}

	webpage.AddChild(hb.Raw(d.layout() + menu))

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
	submenuId := "submenu_" + utils.ToString(index)
	if hasChildren {
		url = "#" + submenuId
	}

	link := hb.Hyperlink().Class("nav-link px-0")

	if icon != "" {
		link.Child(hb.Span().
			Class("icon").
			Style("margin-right: 5px;").
			HTML(icon))
	} else {
		link.Child(hb.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" fill="currentColor" class="bi bi-caret-right-fill" viewBox="0 0 16 16">
	<path d="m12.14 8.753-5.482 4.796c-.646.566-1.658.106-1.658-.753V3.204a1 1 0 0 1 1.659-.753l5.48 4.796a1 1 0 0 1 0 1.506z"/>
</svg>`))
	}
	link.Child(hb.Span().Class("d-inline").HTML(title))
	link.Href(url)
	if hasChildren {
		link.Data("bs-toggle", "collapse")
	}

	return hb.LI().
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

	link := hb.Hyperlink().Class("nav-link align-middle px-0")
	if icon != "" {
		link.Child(hb.Span().Class("icon").Style("margin-right: 5px;").HTML(icon))
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
		link.Child(hb.Raw(html))
	}

	li := hb.LI().Class("nav-item").Child(link)

	if hasChildren {
		ul := hb.UL().
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
	items := []hb.TagInterface{}
	for index, menuItem := range d.menu {
		li := buildMenuItem(menuItem, index)
		items = append(items, li)
	}

	ul := hb.UL().
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
	hasNavbarBackgroundColor := lo.Ternary(d.navbarBackgroundColor == "", false, true)
	hasNavbarTextColor := lo.Ternary(d.navbarTextColor == "", false, true)

	hasLogoImage := lo.Ternary(d.logoImageURL != "", true, false)
	hasLogoRawHTML := lo.Ternary(d.logoRawHtml != "", true, false)
	hasLogo := hasLogoImage || hasLogoRawHTML
	logoRedirectURL := lo.Ternary(d.logoRedirectURL != "", d.logoRedirectURL, "#")

	navbarThemeBackgroundClass := d.navbarBackgroundThemeClass()

	iconStyle := "margin-top:-4px;margin-right:5px;"

	dropdownUser := d.navbarDropdownUser(iconStyle)
	dropdownQuickAccess := d.navbarDropdownQuickAccess(iconStyle)
	dropdownThemeSwitch := d.navbarDropdownThemeSwitch()

	buttonTheme := d.navbarButtonThemeClass()

	buttonMenuToggle := hb.Button().
		Class("btn "+buttonTheme).
		Style("background: none; border:none;").
		StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModalDashboardMenu").
		Children([]hb.TagInterface{
			icons.Icon("bi-list", 24, 24, "").Style(iconStyle),
			hb.Span().HTML("Menu"),
		})

	buttonOffcanvasToggle := hb.Button().
		Class("btn "+buttonTheme).
		Style("background: none; border:none;").
		StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
		Data("bs-toggle", "offcanvas").
		Data("bs-target", "#OffcanvasMenu").
		Children([]hb.TagInterface{
			icons.Icon("bi-list", 24, 24, "").Style(iconStyle),
			hb.Span().HTML("Menu"),
		})

	mainMenu := buttonOffcanvasToggle
	if d.menuType == MENU_TYPE_MODAL {
		mainMenu = buttonMenuToggle
	}

	logo := lo.
		If(hasLogoRawHTML, hb.Raw(d.logoRawHtml)).
		ElseIf(hasLogoImage, hb.Image(d.logoImageURL).Style("max-height:35px;")).
		Else(nil)

	logoLink := hb.Hyperlink().
		Href(logoRedirectURL).
		Class("navbar-brand").
		Child(logo)

	loginLink := hb.Hyperlink().
		Text("Login").
		Href(d.loginURL).
		//Class("btn "+buttonTheme+" float-end").
		Class("btn btn-outline-info float-end").
		StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
		Style("margin-left:10px;")

	registerLink := hb.Hyperlink().
		Text("Register").
		Href(d.registerURL).
		Class("btn "+buttonTheme+" float-end").
		StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
		Style("margin-left:10px;  border:none;")

	toolbar := hb.Nav().
		ID("Toolbar").
		Class("navbar").
		ClassIf(d.navbarHasBackgroundThemeClass(), navbarThemeBackgroundClass).
		Style("z-index: 3;box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);transition: all .2s ease;padding-left: 20px;padding-right: 20px; display:block;").
		StyleIf(hasNavbarBackgroundColor, `background-color: `+d.navbarBackgroundColor+`;`).
		StyleIf(hasNavbarTextColor, `color: `+d.navbarTextColor+`;`).
		ChildIf(hasLogo, logoLink).
		Children([]hb.TagInterface{
			mainMenu,

			// User Menu
			hb.If(!lo.IsEmpty(d.user) && (d.user.FirstName != "" || d.user.LastName != ""),
				hb.Div().Class("float-end").
					Style("margin-left:10px;").
					Child(dropdownUser),
			),

			// Register Link
			hb.If(lo.IsEmpty(d.user) && d.registerURL != "",
				registerLink,
			),

			// Login Link
			hb.If(lo.IsEmpty(d.user) && d.loginURL != "",
				loginLink,
			),

			// Theme Switcher
			hb.If(d.themeHandlerUrl != "",
				hb.Div().Class("float-end").
					Style("margin-left:10px;").
					Child(dropdownThemeSwitch),
			),

			// Quick Menu (if provided)
			hb.If(len(d.quickAccessMenu) > 0, hb.Div().
				Class("float-end").
				Style("margin-left:10px;").
				Child(dropdownQuickAccess)),
		})

	return toolbar.ToHTML()
}

func (d *Dashboard) center(content string) string {
	contentHolder := hb.Div().Class("shadow p-3 m-3").HTML(content)
	html := contentHolder.ToHTML()
	return html
}

func (d *Dashboard) menuOffcanvas() *hb.Tag {
	backgroundClass := d.navbarBackgroundThemeClass()

	offcanvasMenu := hb.Div().
		ID("OffcanvasMenu").
		Class("offcanvas offcanvas-start").
		Class(backgroundClass).
		ClassIfElse(backgroundClass == "bg-light", "text-bg-light", "text-bg-dark").
		Attr("tabindex", "-1").
		Children([]hb.TagInterface{
			hb.Div().Class("offcanvas-header").
				Children([]hb.TagInterface{
					hb.Heading5().
						Class("offcanvas-title").
						Text("Menu"),
					hb.Button().
						Class("btn-close btn-close-white").
						ClassIf(backgroundClass == "bg-light", "text-bg-light").
						Type(hb.TYPE_BUTTON).
						Data("bs-dismiss", "offcanvas").
						Attr("aria-label", "Close"),
				}),
			hb.Div().Class("offcanvas-body").
				Children([]hb.TagInterface{
					hb.Raw(d.DashboardLayoutMenu()),
				}),
		})

	return offcanvasMenu
}

func (d *Dashboard) menuModal() *hb.Tag {
	modalHeader := hb.Div().Class("modal-header").
		Children([]hb.TagInterface{
			hb.Heading5().HTML("Menu").Class("modal-title"),
			hb.Button().Attrs(map[string]string{
				"type":            "button",
				"class":           "btn-close",
				"data-bs-dismiss": "modal",
				"aria-label":      "Close",
			}),
		})

	modalBody := hb.Div().Class("modal-body").Children([]hb.TagInterface{
		hb.Raw(d.DashboardLayoutMenu()),
	})

	modalFooter := hb.Div().Class("modal-footer").Children([]hb.TagInterface{
		hb.Button().
			HTML("Close").
			Class("btn btn-secondary w-100").
			Data("bs-dismiss", "modal"),
	})

	modal := hb.Div().
		ID("ModalDashboardMenu").
		Class("modal fade").
		Children([]hb.TagInterface{
			hb.Div().Class("modal-dialog modal-lg").
				Children([]hb.TagInterface{
					hb.Div().Class("modal-content").
						Children([]hb.TagInterface{
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
// 		placeholderLogo := hb.Image().
// 			Src(logoURL).
// 			Style("width:100%;margin:0px 10px 0px 0px;")
// 		adminDiv := hb.Div().
// 			HTML("ADMIN PANEL").
// 			Style("font-size:12px;text-align: center;")
// 		logo = hb.Div().Class("Logo").Child(placeholderLogo).Child(adminDiv)
// 	} else {
// 		logo = hb.Image().Attr("src", logoURL).Style("width:100%;margin:0px 10px 0px 0px;")
// 	}

// 	sideMenu := hb.Div().ID("SideMenu").Class("p-4").Style("height:100%;width:200px;").
// 		AddChildren([]hb.TagInterface{
// 			logo,
// 			hb.Div().Class("Menu").HTML(menu),
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

func (d *Dashboard) navbarHasBackgroundThemeClass() bool {
	hasNavbarBackgroundColor := lo.Ternary(d.navbarBackgroundColor == "", false, true)
	hasNavbarBackgroundTheme := lo.Ternary(!hasNavbarBackgroundColor && d.navbarBackgroundColorMode != "", true, false)
	return hasNavbarBackgroundTheme
}

func (d *Dashboard) navbarBackgroundThemeClass() string {
	navbarThemeBackgroundClass := lo.
		If(d.navbarHasBackgroundThemeClass(), "bg-"+d.navbarBackgroundColorMode).
		// ElseIf(!hasNavbarBackgroundTheme && !hasNavbarBackgroundColor, "bg-dark").
		Else("")

	return navbarThemeBackgroundClass
}

func (d *Dashboard) navbarButtonThemeClass() string {
	buttonTheme := lo.
		If(d.navbarHasBackgroundThemeClass(), "btn-"+d.navbarBackgroundColorMode).
		Else("")
	return buttonTheme
}

func (d *Dashboard) navbarDropdownQuickAccess(iconStyle string) *hb.Tag {
	hasNavbarTextColor := lo.Ternary(d.navbarTextColor == "", false, true)
	buttonTheme := d.navbarButtonThemeClass()

	button := hb.Button().
		ID("ButtonQuickAccess").
		Class("btn "+buttonTheme+" dropdown-toggle").
		Style("background:none;border:0px;").
		StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
		Type(hb.TYPE_BUTTON).
		Data("bs-toggle", "dropdown").
		Children([]hb.TagInterface{
			icons.Icon("bi-microsoft", 24, 24, "").
				Style(iconStyle).
				Style("margin-top:-4px;margin-right:8px;"),
			hb.Span().Text("Quick Access").Style("margin-right:10px;"),
		})

	dropdownQuickAccess := hb.Div().
		Class("dropdown").
		Style(`margin:0px;`).
		Child(button).
		Child(hb.UL().
			Class("dropdown-menu").
			Children(lo.Map(d.quickAccessMenu, func(item MenuItem, _ int) hb.TagInterface {
				target := lo.Ternary(item.Target == "", "_self", item.Target)
				url := lo.Ternary(item.URL == "", "#", item.URL)

				return hb.LI().Children([]hb.TagInterface{
					hb.If(item.Title == "",
						hb.HR().
							Class("dropdown-divider"),
					),

					hb.If(item.Title != "",
						hb.Hyperlink().
							Class("dropdown-item").
							ChildIf(item.Icon != "", hb.Span().Class("icon").Style("margin-right: 5px;").HTML(item.Icon)).
							Text(item.Title).
							Href(url).
							Target(target),
					),
				})
			})))

	return dropdownQuickAccess
}

// themeButton generates a dropdown menu with light and dark themes.
//
// It checks if the current theme is dark and creates dropdown items for both light and dark themes.
// The dropdown items are created dynamically based on the themesLight and themesDark maps.
// The function returns a *hb.Tag that represents the generated dropdown menu.
func (d *Dashboard) navbarDropdownThemeSwitch() *hb.Tag {
	hasNavbarTextColor := lo.Ternary(d.navbarTextColor == "", false, true)
	buttonTheme := d.navbarButtonThemeClass()

	isDark := d.isThemeDark()

	// Light Themes
	lightDropdownItems := lo.Map(lo.Keys(themesLight), func(theme string, index int) hb.TagInterface {
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

		return hb.LI().Children([]hb.TagInterface{
			hb.Hyperlink().
				Class("dropdown-item"+active).
				Child(hb.I().Class("bi bi-sun").Style("margin-right:5px;")).
				HTML(name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	// Dark Themes
	darkDropdownItems := lo.Map(lo.Keys(themesDark), func(theme string, index int) hb.TagInterface {
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

		return hb.LI().Children([]hb.TagInterface{
			hb.Hyperlink().
				Class("dropdown-item"+active).
				Child(hb.I().Class("bi bi-moon-stars-fill").Style("margin-right:5px;")).
				HTML(name).
				Href(url).
				Attr("ref", "nofollow"),
		})
	})

	button := bs.Button().
		ID("buttonTheme").
		Class(buttonTheme+" dropdown-toggle").
		Style("background:none;border:0px;").
		StyleIf(hasNavbarTextColor, "color:"+d.navbarTextColor).
		Data("bs-toggle", "dropdown").
		Children([]hb.TagInterface{
			lo.Ternary(isDark, hb.I().Class("bi bi-sun"), hb.I().Class("bi bi-moon-stars-fill")),
		})

	return hb.Div().
		Class("dropdown").
		Style(`margin:0px;`).
		Child(button).
		Child(hb.UL().
			Class(buttonTheme+" dropdown-menu dropdown-menu-dark").
			Children(lightDropdownItems).
			ChildIf(
				len(lo.Filter(darkDropdownItems, func(item hb.TagInterface, _ int) bool { return item != nil })) > 0 && len(lo.Filter(lightDropdownItems, func(item hb.TagInterface, _ int) bool { return item != nil })) > 0,
				hb.LI().Children([]hb.TagInterface{
					hb.HR().Class("dropdown-divider"),
				}),
			).
			Children(darkDropdownItems))
}

func (d *Dashboard) navbarDropdownUser(iconStyle string) *hb.Tag {
	hasNavbarTextColor := lo.Ternary(d.navbarTextColor == "", false, true)
	buttonTheme := d.navbarButtonThemeClass()

	dropdownUser := hb.Div().
		Class("dropdown").
		Children([]hb.TagInterface{
			hb.Button().
				ID("ButtonUser").
				Class("btn "+buttonTheme+" dropdown-toggle").
				Style("background:none;border:0px;").
				StyleIf(hasNavbarTextColor, "color: "+d.navbarTextColor+";").
				Type(hb.TYPE_BUTTON).
				Data("bs-toggle", "dropdown").
				Children([]hb.TagInterface{
					icons.Icon("bi-person", 24, 24, "").Style(iconStyle),
					hb.Span().
						Text(d.user.FirstName + " " + d.user.LastName).
						Style("margin-right:10px;"),
				}),
			hb.UL().
				Class("dropdown-menu dropdown-menu-dark").
				Class(buttonTheme).
				Children(lo.Map(d.userMenu, func(item MenuItem, _ int) hb.TagInterface {
					target := lo.Ternary(item.Target == "", "_self", item.Target)
					url := lo.Ternary(item.URL == "", "#", item.URL)

					return hb.LI().Children([]hb.TagInterface{
						hb.If(item.Title == "",
							hb.HR().
								Class("dropdown-divider"),
						),

						hb.If(item.Title != "",
							hb.Hyperlink().
								Class("dropdown-item").
								ChildIf(item.Icon != "", hb.Span().Class("icon").Style("margin-right: 5px;").HTML(item.Icon)).
								Text(item.Title).
								Href(url).
								Target(target),
						),
					})
				})),
		})

	return dropdownUser
}
