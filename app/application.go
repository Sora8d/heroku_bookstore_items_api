package app

import (
	"net/http"

	"github.com/Sora8d/heroku_bookstore_items_api/config"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func StartApplication() {
	Urlmaps()

	srv := &http.Server{
		Handler: router,
		Addr:    config.Config["address"],
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
