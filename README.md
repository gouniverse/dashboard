# Dashboard <a href="https://gitpod.io/#https://github.com/gouniverse/dashboard" style="float:right;"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/gouniverse/dashboard/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/dashboard)](https://goreportcard.com/report/github.com/gouniverse/dashboard)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/dashboard)](https://pkg.go.dev/github.com/gouniverse/dashboard)

This is a project for quickly building a dashboard.

Out of the box it provides a general layout, main menu,
user dropdown menu, quick access menu, and theme switcher.

The content of the page itself is left blank. It is up to the developer
to customize the pages as needed by the application he is building.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Features

- Uses the latest [Bootstrap](https://getbootstrap.com) (v5.3.3)
- Uses the latest [Bootstrap Icons](https://icons.getbootstrap.com/) (v1.11.3)
- Preset sidebar menu
- Preset user dropdown menu
- Preset quick access menu
- Preset menu switcher
- Supported all [Bootswatch](https://bootswatch.com/) themes (26 total)
- Flexible component system for building UI elements
- Consistent functional API for all components

## Components

The dashboard includes a flexible component system that allows you to build custom UI elements. All components follow a consistent functional pattern that returns `hb.TagInterface` directly, making them easy to compose and combine.

### ShadowBoxComponent

The ShadowBoxComponent adds shadow, padding, and margin to your content, giving it a clean, separated look.

```golang
shadowBox := components.NewShadowBoxComponent(components.ShadowBoxConfig{
    Content: "<h1>Hello World</h1>",
    Padding: 15,
    Margin: 10,
})
```

### Card Component

The Card component creates a Bootstrap card with an optional title and customizable styling.

```golang
card := components.NewCard(components.CardConfig{
    Title:   "Card Title",
    Content: "<p>This is content inside a Bootstrap card component.</p>",
    Margin:  10,
})
```

### Tab Layout Component

The TabLayout component organizes content into tabs that users can switch between.

```golang
tabLayout := components.NewTabLayout(components.TabLayoutConfig{
    Items: []components.TabItem{
        {
            ID:      "tab1",
            Title:   "First Tab",
            Content: "<p>Content of the first tab.</p>",
            Active:  true,
        },
        {
            ID:      "tab2",
            Title:   "Second Tab",
            Content: "<p>Content of the second tab.</p>",
        },
    },
    Margin: 10,
})
```

### Grid Component

The Grid component creates a responsive Bootstrap grid system with rows and columns.

```golang
// Create a grid layout with two rows
row1 := []components.GridItem{
    {Content: "<div class='p-3 bg-light'>Column 1</div>", ColumnClass: "col-md-6"},
    {Content: "<div class='p-3 bg-light'>Column 2</div>", ColumnClass: "col-md-6"},
}

// Second row with three columns of different widths
row2 := []components.GridItem{
    {Content: "<div class='p-3 bg-light'>Column 1</div>", ColumnClass: "col-md-3"},
    {Content: "<div class='p-3 bg-light'>Column 2</div>", ColumnClass: "col-md-6"},
    {Content: "<div class='p-3 bg-light'>Column 3</div>", ColumnClass: "col-md-3"},
}

grid := components.NewGridComponent(components.GridLayoutConfig{
    Rows:    [][]components.GridItem{row1, row2},
    Gutters: 3,
    Margin:  10,
})
```

### Combining Components

Components can be easily combined by using the `.ToHTML()` method:

```golang
// Create tabs
tabs := components.NewTabLayout(components.TabLayoutConfig{
    Items: []components.TabItem{
        {
            ID:      "tab1",
            Title:   "First Tab",
            Content: "<p>Content of the first tab.</p>",
            Active:  true,
        },
        {
            ID:      "tab2",
            Title:   "Second Tab",
            Content: "<p>Content of the second tab.</p>",
        },
    },
})

// Put the tabs inside a card
card := components.NewCard(components.CardConfig{
    Title:   "Card with Tabs",
    Content: tabs.ToHTML(),
})

// Put the card inside a shadow box
shadowBox := components.NewShadowBoxComponent(components.ShadowBoxConfig{
    Content: card.ToHTML(),
    Padding: 15,
    Margin:  10,
})
```

## Running the Examples

The dashboard includes a runnable example that demonstrates all the components in action:

```bash
cd examples
go run component_examples.go
```

Then open your browser to http://localhost:8080 to see the examples.

## Example

- Adding to an HTTP handler

```golang
func dashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := dashboard.NewDashboard(dashboard.Config{
		Menu: []dashboard.MenuItem{
			{
				Title: "Home",
				URL:   "/",
			},
			{
				Title: "Logout",
				URL:   "/auth/logout",
			},
		},
	})

	w.Write([]byte(dashboard.ToHTML()))
}
```

- Adding to a layout function, to reuse on multiple places

```golang
func layout(r *http.Request, opts AdminDashboardOptions) string {
    authUser := helpers.GetAuthUser(r)

    logoImageURL = "YOUR_IMAGE_URL.png"
	logoRedirectURL = "/"

    dashboardMenuItems := []dashboard.MenuItem{
            {
                Title: "Home",
                URL:   links.NewAdminLinks().Home(map[string]string{}),
            },
            {
                Title: "Blog Manager",
                URL:   links.NewAdminLinks().Blog(map[string]string{}),
            },
            {
                Title: "Website Manager",
                URL:   links.NewAdminLinks().Cms(map[string]string{}),
            },
            {
                Title: "User Manager",
                URL:   links.NewAdminLinks().Users(map[string]string{}),
            },
        }

    dashboardUser := dashboard.User{
            FirstName: authUser.FirstName(),
            LastName:  authUser.LastName(),
    }

    dashboardQuickAccessMenuItems := []dashboard.MenuItem {
        {
            Title: "New post",
            URL: "/post-create",
        },
        {
            Title: "New page",
            URL: "/page-create",
        }
    }
        
    dashboardUserMenuItems := []dashboard.MenuItem {
        {
            Title: "Profile",
            URL: "/account/profile",
        },
        {
            Title: "Logout",
            URL: "/auth/logout",
        }
    }
        
    dashboard := dashboard.NewDashboard(dashboard.Config{
        HTTPRequest:                r,
        Content:                    opts.Content,
        Title:                      opts.Title,
        LogoImageURL                logoImageURL,
        LogoRedirectURL             logoRedirectURL,
        MenuItems:                  dashboardMenuItems,
        User:                       dashboardUser,
        UserMenuItems:              dashboardUserMenuItems,
        QuickAccessMenuItems:       dashboardQuickAccessMenuItems,
        Scripts:                    opts.Scripts,
        ScriptURLs:                 opts.ScriptURLs,
        Styles:                     opts.Styles,
        StyleURLs:                  opts.StyleURLs,
        // optional, defaults to dark
        // NavbarBackgroundColorMode: "light"
        // optional, defaults to the default Bootstrap theme
        // ThemeName:                 dashboard.THEME_MINTY,
        // ThemeHandlerUrl:      links.NewAdminLinks().Theme(map[string]string{"redirect": links.NewAdminLinks().Home(map[string]string{})}),   // Optional (Advanced)
        // UncdnHandlerEndpoint: links.NewAdminLinks().Uncdn(map[string]string{}),                                                              // Optional (Advanced)
    })

    return dashboard.ToHTML()
}
```

## Screenshots

- Main View

<img src="./screenshots/screenshot_main_view_20230712.png" />

- Main Menu

<img src="./screenshots/screenshot_main_menu_20230712.png" />

- Quick Access Menu

<img src="./screenshots/screenshot_quick_access_menu_20230712.png" />

- User Menu

<img src="./screenshots/screenshot_user_menu_20230712.png" />

## Development
For working on this package:
- Open in Gitpod (use the button provided)
- Run these commands sequentially
- Open the browser URL displayed in the terminal
```
task dev:init
task dev
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/gouniverse/dashboard.svg)](https://starchart.cc/gouniverse/dashboard)

## Noteworthy

- https://github.com/pro-dev-ph/bootstrap-simple-admin-template

- https://github.com/PlainAdmin/plain-free-bootstrap-admin-template

- https://github.com/tabler/tabler

- https://github.com/puikinsh/Adminator-admin-dashboard

- https://github.com/themesberg/volt-bootstrap-5-dashboard

- https://dribbble.com/shots/19114068-Dashboard

- https://demo.themefisher.com/focus/

## Similar Golang Projects

- https://github.com/oal/admin

- https://github.com/uadmin/uadmin

- https://github.com/GoAdminGroup/go-admin

- https://github.com/entkit/entkit
