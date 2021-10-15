package app

import (
	"net/http"

	"github.com/Sora8d/heroku_bookstore_items_api/controllers"
)

func Urlmaps() {
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	router.HandleFunc("/items", controllers.ItemsController.Create).Methods(http.MethodPost)
	router.HandleFunc("/items/{id}", controllers.ItemsController.Get).Methods(http.MethodGet)
	router.HandleFunc("/items/search", controllers.ItemsController.Search).Methods(http.MethodPost)
}
