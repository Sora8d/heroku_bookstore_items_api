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
	"github.com/Sora8d/heroku_bookstore_items_api/services"
	"github.com/Sora8d/heroku_bookstore_items_api/utils/http_utils"
	"github.com/gorilla/mux"
)

func init() {
	oauth.OauthRestClient.SetClient("http://127.0.0.1:8081")
}

var ItemsController itemsControllerInterface = &itemsController{}

const (
	headerXCallerId = "X-User-Id"
	headerXAdmin    = "X-Admin"
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
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
		respErr := rest_errors.NewUnauthorizedError("invalid request body")
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
	itemId, idErr := getId(r)
	if idErr != nil {
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

func (it *itemsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		//TODO: Return error to the caller
		http_utils.RespondJson(w, err.Status(), err)
		return
	}

	var permissions bool = false
	if r.Header.Get(headerXAdmin) == "true" {
		permissions = true
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

	itemId, idErr := getId(r)
	if idErr != nil {
		http_utils.RespondJson(w, idErr.Status(), idErr)
		return
	}
	itemRequest.Id = itemId

	itemRequest.Seller = oauth.GetCallerId(r)

	result, respErr := services.ItemsService.Update(itemRequest, permissions)
	if respErr != nil {
		http_utils.RespondJson(w, respErr.Status(), respErr)
		return
	}
	http_utils.RespondJson(w, http.StatusCreated, result)
}

func (it *itemsController) Delete(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		http_utils.RespondJson(w, err.Status(), err)
		return
	}
	itemId, idErr := getId(r)
	if idErr != nil {
		http_utils.RespondJson(w, idErr.Status(), idErr)
		return
	}
	itemSeller := oauth.GetCallerId(r)

	var permissions bool = false
	if r.Header.Get(headerXAdmin) == "true" {
		permissions = true
	}

	requestItem := items.Item{Id: itemId, Seller: itemSeller}

	if err := services.ItemsService.Delete(requestItem, permissions); err != nil {
		http_utils.RespondJson(w, err.Status(), err)
		return
	}

	http_utils.RespondJson(w, http.StatusOK, map[string]string{"stauts": "deleted"})

}

func getId(r *http.Request) (int64, rest_errors.RestErr) {
	vars := mux.Vars(r)
	itemId, err := strconv.ParseInt(strings.TrimSpace(vars["id"]), 10, 64)
	if err != nil {
		idErr := rest_errors.NewBadRequestErr("invalid id, must be a number")
		return 0, idErr
	}
	return itemId, nil
}
