package app

import (
	"fmt"
	"net/http"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/heroku_bookstore_items_api/config"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var address = fmt.Sprintf("%s:%s", config.Config["address"], config.Config["port"])

func StartApplication() {
	Urlmaps()
	logger.Info(fmt.Sprintf("starting app, trying connection in: %s", address))
	srv := &http.Server{
		Handler: router,
		Addr:    address,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("error starting app", err)
		panic(err)
	}
}
