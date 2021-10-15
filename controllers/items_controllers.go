package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sora8d/bookstore_oauth-go/oauth"
	"github.com/Sora8d/bookstore_utils-go/rest_errors"
	"github.com/Sora8d/heroku_bookstore_items_api/domain/items"
	"github.com/Sora8d/heroku_bookstore_items_api/domain/queries"
	"github.com/Sora8d/heroku_bookstore_items_api/localerrors"
	"github.com/Sora8d/heroku_bookstore_items_api/services"
	"github.com/Sora8d/heroku_bookstore_items_api/utils/http_utils"
	"github.com/gorilla/mux"
)

var ItemsController itemsControllerInterface = &itemsController{}

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

type itemsController struct {
}

func (it *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		//TODO: Return error to the caller
		http_utils.RespondJson(w, err.Status(), err)
		return
	}
	sellerId := oauth.GetCallerId(r)
	if sellerId == 0 {
		respErr := localerrors.NewUnauthorizedError("invalid request body")
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestErr("invalid request body")
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}
	defer r.Body.Close()

	var itemRequest items.Item
	if err := json.Unmarshal(requestBody, &itemRequest); err != nil {
		respErr := rest_errors.NewBadRequestErr("invalid json body")
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}

	itemRequest.Seller = sellerId

	result, respErr := services.ItemsService.Create(itemRequest)
	if respErr != nil {
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}
	http_utils.RespondJson(w, http.StatusCreated, result)
}

func (it *itemsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId, err := strconv.ParseInt(strings.TrimSpace(vars["id"]), 10, 64)
	if err != nil {
		idErr := rest_errors.NewBadRequestErr("invalid id, must be a number")
		http_utils.RespondJson(w, idErr.Status(), idErr)
		return
	}

	item, idErr := services.ItemsService.Get(itemId)
	if idErr != nil {
		http_utils.RespondJson(w, idErr.Status(), idErr)
		return
	}
	http_utils.RespondJson(w, http.StatusOK, item)
}

func (it *itemsController) Search(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestErr("invalid json body")
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}
	defer r.Body.Close()

	var query queries.PsQuery
	if err := json.Unmarshal(bytes, &query); err != nil {
		respErr := rest_errors.NewBadRequestErr("invalid json body")
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}

	items, searchErr := services.ItemsService.Search(query)
	if searchErr != nil {
		http_utils.RespondJson(w, searchErr.Status(), searchErr)
		return
	}
	http_utils.RespondJson(w, http.StatusOK, items)
}
