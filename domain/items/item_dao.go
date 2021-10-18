package items

import (
	"context"
	"errors"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/bookstore_utils-go/rest_errors"
	"github.com/Sora8d/heroku_bookstore_items_api/clients/postgresql"
	"github.com/Sora8d/heroku_bookstore_items_api/domain/queries"
	"github.com/jackc/pgx/v4"
)

const (
	queryGetItem                = "SELECT id, seller, title, price FROM item WHERE id=$1;"
	queryGetItemsDescription    = "SELECT d.plain_text, d.html FROM description d WHERE d.item_id=$1;"
	queryGetItemsPictures       = "SELECT p.id, p.url FROM picture p WHERE p.item_id=$1;"
	querySaveItem               = "INSERT INTO item(seller, title, price) VALUES ($1, $2, $3) RETURNING id;"
	querySaveItemsDescription   = "INSERT INTO description(item_id, plain_text, html) VALUES($1,$2,$3);"
	querySaveItemsPictures      = "INSERT INTO picture(item_id, url) VALUES($1, $2);"
	queryUpdateItem             = "UPDATE item SET seller=$1, title=$2, price=$3 WHERE id=$4;"
	queryUpdateItemsDescription = "UPDATE description SET plain_text= $1, html=$2 WHERE item_id=$3;"
	queryDeleteItem             = "DELETE FROM item WHERE id=$1;"
	queryDeleteItemDescription  = "DELETE FROM description WHERE item_id=$1"
	queryDeleteItemPictures     = "DELETE FROM picture WHERE item_id=$1"
	queryBuild                  = "SELECT item.id, item.seller, item.title FROM item INNER JOIN description d ON item.id = d.item_id WHERE %s;"

	picturetable = "picture"
	errornorows  = "no rows in result set"
)

var (
	picturecolumns = [2]string{"item_id", "url"}
	ctx            = context.Background()
)

func (i *Item) Save() rest_errors.RestErr {
	ctx := context.Background()

	tx, err := postgresql.Client.Transaction()
	if err != nil {
		logger.Error("There was an error in the save function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, querySaveItem, i.Seller, i.Title, i.Price)
	var resultId int64
	if err := row.Scan(&resultId); err != nil {
		logger.Error("There was an error in the save function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	_, err = tx.Exec(ctx, querySaveItemsDescription, resultId, i.Description.PlainText, i.Description.Html)
	if err != nil {
		logger.Error("There was an error in the save function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	if err = savePictures(i, resultId, tx); err != nil {

		logger.Error("There was an error in the save function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))

	}
	/*
		pictures_row := i.GetPictures(resultId)
		_, err = tx.CopyFrom(ctx, pgx.Identifier{picturetable}, picturecolumns[:], pgx.CopyFromRows(pictures_row))
		if err != nil {
			logger.Error("There was an error in the save function in items dao", err)
			return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
		}
	*/
	tx.Commit(ctx)

	i.Id = resultId
	return nil
}

func (i *Item) Get() rest_errors.RestErr {

	tx, err := postgresql.Client.Transaction()
	if err != nil {
		logger.Error("There was an error in the get function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
	}
	defer tx.Rollback(ctx)

	itemRow := tx.QueryRow(ctx, queryGetItem, i.Id)
	if err := itemRow.Scan(&i.Id, &i.Seller, &i.Title, &i.Price); err != nil {
		if err.Error() == errornorows {
			return rest_errors.NewNotFoundError("No item with given id")
		}
		logger.Error("There was an error in the get function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
	}
	/*
		descRow := tx.QueryRow(ctx, queryGetItemsDescription, i.Id)
		if err := descRow.Scan(&i.Description.PlainText, &i.Description.Html); err != nil {
			return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
		}
	*/
	if err := getDescription(i, tx); err != nil {
		logger.Error("There was an error in the get function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
	}
	/*
		picRows, err := tx.Query(ctx, queryGetItemsPictures, i.Id)
		if err != nil {
			return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
		}
		defer picRows.Close()
		var pics []Picture
		for picRows.Next() {
			var currentPic Picture
			err = picRows.Scan(&currentPic.Id, &currentPic.Url)
			if err != nil {
				return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
			}
		}
		i.Pictures = pics
	*/
	if err := getPictures(i, tx); err != nil {
		logger.Error("There was an error in the get function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
	}
	tx.Commit(ctx)
	return nil
}

func (i *Item) Search(query queries.PsQuery) ([]Item, rest_errors.RestErr) {
	var items []Item
	q, values := query.Build(queryBuild)
	rows, err := postgresql.Client.Query(ctx, q, values...)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("There was an  error building the query", errors.New("database error"))
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.Seller, &item.Title, &item.Price); err != nil {
			return nil, rest_errors.NewInternalServerError("There was an error parsin the search results", errors.New("database error"))
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, rest_errors.NewNotFoundError("No items with specified criteria")
	}
	for index, item := range items {
		tx, err := postgresql.Client.Transaction()
		if err != nil {
			logger.Error("There was an error in the search function in items dao", err)
			return nil, rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
		}
		if err := getDescription(&item, tx); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
		}
		if err := getPictures(&item, tx); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to get item", errors.New("database error"))
		}
		tx.Commit(ctx)
		items[index] = item
		tx.Rollback(ctx)
	}
	return items, nil
}

func (i *Item) Delete() rest_errors.RestErr {
	tx, err := postgresql.Client.Transaction()
	if err != nil {
		logger.Error("There was an error in the update function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to update item", errors.New("database error"))
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, queryDeleteItemDescription, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}
	_, err = tx.Exec(ctx, queryDeleteItemPictures, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}
	_, err = tx.Exec(ctx, queryDeleteItem, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}
	tx.Commit(ctx)
	return nil
}

func (i *Item) Update() rest_errors.RestErr {
	tx, err := postgresql.Client.Transaction()
	if err != nil {
		logger.Error("There was an error in the update function in items dao", err)
		return rest_errors.NewInternalServerError("error when trying to update item", errors.New("database error"))
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, queryUpdateItem, i.Seller, i.Title, i.Price, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}

	_, err = tx.Exec(ctx, queryUpdateItemsDescription, i.Description.PlainText, i.Description.Html, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}

	_, err = tx.Exec(ctx, queryDeleteItemPictures, i.Id)
	if err != nil {
		logger.Error("Catching to check if not found is considered an err", err)
		return rest_errors.NewBadRequestErr("There was an error trying to update this item in the database")
	}

	if err = savePictures(i, i.Id, tx); err != nil {
		return rest_errors.NewBadRequestErr("There was an error updating the pictures of the item")
	}

	tx.Commit(ctx)
	return nil
}

func getDescription(i *Item, client postgresql.TxandClient) error {
	descRow := client.QueryRow(ctx, queryGetItemsDescription, i.Id)
	if err := descRow.Scan(&i.Description.PlainText, &i.Description.Html); err != nil {
		if err.Error() == errornorows {
			return nil
		}
		return err
	}
	return nil
}

func getPictures(i *Item, client postgresql.TxandClient) error {
	picRows, err := client.Query(ctx, queryGetItemsPictures, i.Id)
	if err != nil {
		return err
	}
	defer picRows.Close()
	var pics []Picture
	for picRows.Next() {
		var currentPic Picture
		err = picRows.Scan(&currentPic.Id, &currentPic.Url)
		if err != nil {
			return err
		}
		pics = append(pics, currentPic)
	}
	i.Pictures = pics
	return nil
}

func savePictures(i *Item, id int64, client postgresql.TxandClient) error {
	pictures_row := i.GetPictures(id)
	_, err := client.CopyFrom(ctx, pgx.Identifier{picturetable}, picturecolumns[:], pgx.CopyFromRows(pictures_row))
	if err != nil {
		return err
	}
	return nil
}
