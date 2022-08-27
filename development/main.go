package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/utils"
)

func main() {
	log.Println("1. Initializing environment variables...")
	utils.EnvInitialize()

	log.Println("3. Initializing CMS...")
	// myCms, err := cms.NewCms(cms.Config{
	// 	DbInstance:          db,
	// 	BlocksEnable:        true,
	// 	CacheAutomigrate:    true,
	// 	CacheEnable:         true,
	// 	EntitiesAutomigrate: true,
	// 	LogsAutomigrate:     true,
	// 	LogsEnable:          true,
	// 	MenusEnable:         true,
	// 	PagesEnable:         true,
	// 	SettingsAutomigrate: true,
	// 	SettingsEnable:      true,
	// 	SessionAutomigrate:  true,
	// 	SessionEnable:       true,
	// 	TemplatesEnable:     true,
	// 	Prefix:              "cms2",
	// 	CustomEntityList:    entityList(),
	// })

	// if err != nil {
	// 	log.Panicln(err.Error())
	// }

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	// mux.HandleFunc("/", myCms.Router)
	// mux.HandleFunc("/cms", myCms.Router)

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
