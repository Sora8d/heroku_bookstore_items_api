package controllers

import (
	"net/http"
)

var PingController pingControllersInterface = &pingController{}

type pingControllersInterface interface {
	Ping(http.ResponseWriter, *http.Request)
}

type pingController struct {
}

func (pc *pingController) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("pong"))
}
