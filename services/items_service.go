package services

import (
	"github.com/Sora8d/bookstore_utils-go/rest_errors"
	"github.com/Sora8d/heroku_bookstore_items_api/domain/items"
	"github.com/Sora8d/heroku_bookstore_items_api/domain/queries"
)

var ItemsService itemsServiceInterface = &itemsService{}

type itemsServiceInterface interface {
	Create(item items.Item) (*items.Item, rest_errors.RestErr)
	Get(int64) (*items.Item, rest_errors.RestErr)
	Search(queries.PsQuery) ([]items.Item, rest_errors.RestErr)
}

type itemsService struct{}

func NewService() itemsServiceInterface {
	return &itemsService{}
}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {
	if err := item.Save(); err != nil {
		return nil, err
	}
	return &item, nil

}

func (s *itemsService) Get(id int64) (*items.Item, rest_errors.RestErr) {
	item := items.Item{Id: id}

	if err := item.Get(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Search(query queries.PsQuery) ([]items.Item, rest_errors.RestErr) {
	dao := items.Item{}
	return dao.Search(query)
}
