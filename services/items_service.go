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
	Update(items.Item, bool) (*items.Item, rest_errors.RestErr)
	Delete(items.Item, bool) rest_errors.RestErr
}

type itemsService struct{}

func NewService() itemsServiceInterface {
	return &itemsService{}
}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {
	if err := item.Save(); err != nil {
		return nil, err
	}
	if err := item.Get(); err != nil {
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

func (s *itemsService) Update(item items.Item, permissions bool) (*items.Item, rest_errors.RestErr) {
	original_item := items.Item{Id: item.Id}
	if err := original_item.Get(); err != nil {
		return nil, err
	}
	if original_item.Seller != item.Seller && !permissions {
		return nil, rest_errors.NewUnauthorizedError("cant update item, bad credentials")
	}
	if item.Title != "" {
		original_item.Title = item.Title
	}
	if len(item.Pictures) != 0 {
		original_item.Pictures = item.Pictures
	}
	if item.Description.PlainText != "" {
		original_item.Description.PlainText = item.Description.PlainText
	}
	if item.Description.Html != "" {
		original_item.Description.Html = item.Description.Html
	}
	if item.Price != 0 {
		original_item.Price = item.Price
	}
	if err := original_item.Update(); err != nil {
		return nil, err
	}
	if err := item.Get(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Delete(item items.Item, permissions bool) rest_errors.RestErr {
	sellerId := item.Seller
	if err := item.Get(); err != nil {
		return err
	}
	if item.Seller != sellerId && !permissions {
		return rest_errors.NewUnauthorizedError("cant delete item, invalid credentials")
	}

	if err := item.Delete(); err != nil {
		return err
	}
	return nil
}

func (s *itemsService) Search(query queries.PsQuery) ([]items.Item, rest_errors.RestErr) {
	dao := items.Item{}
	return dao.Search(query)
}
