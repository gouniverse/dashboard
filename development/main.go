package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/dashboard"
	"github.com/gouniverse/utils"
)

func main() {
	log.Println("1. Initializing environment variables...")
	utils.EnvInitialize()

	// log.Println("3. Initializing Dashboard...")

	log.Println("2. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", dashboard1)
	mux.HandleFunc("/dashboard-2", dashboard2)

	srv := &http.Server{
		Handler: mux,
		Addr:    utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func dashboard1(w http.ResponseWriter, r *http.Request) {
	dashboardUser := dashboard.User{
		FirstName: "John",
		LastName:  "Doe",
	}
	dashboard := dashboard.NewDashboard(dashboard.Config{
		Menu: []dashboard.MenuItem{
			{
				Title: "Home",
				URL:   "/",
				Icon:  `<i class="bi bi-house"></i>`,
			},
			{
				Title: "Website Manager",
				URL:   "/website-manager",
				Icon:  `<i class="bi bi-globe"></i>`,
			},
			{
				Title: "Blog Manager",
				URL:   "/blog-manager",
				Icon:  `<i class="bi bi-newspaper"></i>`,
			},
			{
				Title: "User Manager",
				URL:   "/user-manager",
				Icon:  `<i class="bi bi-people"></i>`,
			},
		},
		QuickAccessMenu: []dashboard.MenuItem{
			{
				Title: "Add new post",
				URL:   "/post-create",
				Icon:  `<i class="bi bi-plus-circle"></i>`,
			},
			{
				Title: "Add new user",
				URL:   "/user-create",
				Icon:  `<i class="bi bi-plus-circle"></i>`,
			},
		},
		User: dashboardUser,
		UserMenu: []dashboard.MenuItem{
			{
				Title: "Profile",
				URL:   "/account/profile",
				Icon:  `<i class="bi bi-pencil"></i>`,
			},
			{
				Title: "",
			},
			{
				Title: "Logout",
				URL:   "/logout",
				Icon:  `<i class="bi bi-box-arrow-right"></i>`,
			},
		},
		ThemeName:                 dashboard.THEME_MINTY,
		ThemeHandlerUrl:           "/theme-switcher",
		NavbarBackgroundColorMode: "light",
	}).ToHTML()

	w.Write([]byte(dashboard))
}

func dashboard2(w http.ResponseWriter, r *http.Request) {
	dashboardUser := dashboard.User{
		FirstName: "John",
		LastName:  "Doe",
	}
	dashboard := dashboard.NewDashboard(dashboard.Config{
		Menu: []dashboard.MenuItem{
			{
				Title: "Dashboard 1",
				URL:   "/",
			},
			{
				Title: "Dashboard 2",
				URL:   "/dashboard-2",
			},
		},
		User:                      dashboardUser,
		ThemeName:                 dashboard.THEME_MINTY,
		NavbarBackgroundColorMode: "light",
	})
	html := dashboard.ToHTML()
	w.Write([]byte(html))
}
