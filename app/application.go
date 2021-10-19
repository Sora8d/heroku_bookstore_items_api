package app

import (
	"net/http"

	"github.com/Sora8d/heroku_bookstore_items_api/config"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var port = config.Config["port"]

func StartApplication() {
	Urlmaps()
	if port == "" {
		port = ":8080"
	}
	srv := &http.Server{
		Handler: router,
		//		Addr:    config.Config["address"] + port,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
