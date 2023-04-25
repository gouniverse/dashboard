# Dashboard

This is a project for quickly building dashboards


## Example

- Adding to HTTP handled

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
	html := dashboard.ToHTML()
	w.Write([]byte(html))
}
```

- Adding to layout function, to reuse on multiple places
```golang
func layout(r *http.Request, opts AdminDashboardOptions) string {
    authUser := helpers.GetAuthUser(r)

    dashboardMenu := []dashboard.MenuItem{
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
        
    dashboard := dashboard.NewDashboard(dashboard.Config{
        HTTPRequest: r,
        Content:     opts.Content,
        Title:       opts.Title,
        Menu: ,
        User: dashboardUser,
        // ThemeHandlerUrl:      links.NewAdminLinks().Theme(map[string]string{"redirect": links.NewAdminLinks().Home(map[string]string{})}),
        // UncdnHandlerEndpoint: links.NewAdminLinks().Uncdn(map[string]string{}),
        Scripts:              opts.Scripts,
        ScriptURLs:           opts.ScriptURLs,
        Styles:               opts.Styles,
        StyleURLs:            opts.StyleURLs,
    })

    return dashboard.ToHTML()
}
```

## Noteworthy

- https://github.com/PlainAdmin/plain-free-bootstrap-admin-template

- https://github.com/tabler/tabler

- https://github.com/puikinsh/Adminator-admin-dashboard

- https://github.com/themesberg/volt-bootstrap-5-dashboard

- https://dribbble.com/shots/19114068-Dashboard

- https://demo.themefisher.com/focus/
