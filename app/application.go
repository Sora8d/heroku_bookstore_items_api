package app

import (
	"net/http"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/heroku_bookstore_items_api/config"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var port = config.Config["port"]

func StartApplication() {
	Urlmaps()
	srv := &http.Server{
		Handler: router,
		Addr:    config.Config["address"] + port,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("error starting app", err)
		panic(err)
	}
}
