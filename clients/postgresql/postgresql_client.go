package postgresql

import (
	"context"
	"fmt"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/Sora8d/heroku_bookstore_items_api/config"

	pgx "github.com/jackc/pgx/v4"
)

var (
	Client postGresInterface
)

type TxandClient interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}

type postGresInterface interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Execute(context.Context, string, ...interface{}) error
	Insert(context.Context, string, ...interface{}) pgx.Row
	Transaction() (pgx.Tx, error)
}

func init() {
	datasourceName := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		config.Config["items_postgres_username"],
		config.Config["items_postgres_password"],
		config.Config["items_postgres_schema"])
	newConn, err := pgx.Connect(context.Background(), datasourceName)
	if err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	if err = newConn.Ping(context.Background()); err != nil {
		logger.Error("Fatal error initializing db", err)
		panic(err)
	}
	Client = &postGresObject{conn: newConn}
	logger.Info("database succesfully configured")
}

type postGresObject struct {
	conn *pgx.Conn
}

func (pgc postGresObject) QueryRow(ctx context.Context, query string, id ...interface{}) pgx.Row {
	row := pgc.conn.QueryRow(ctx, query, id)
	return row
}

func (pgc postGresObject) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := pgc.conn.Query(ctx, query, args...)
	return rows, err
}
func (pgc postGresObject) Execute(ctx context.Context, query string, args ...interface{}) error {
	_, err := pgc.conn.Exec(ctx, query, args...)
	return err
}

func (pgc postGresObject) Insert(ctx context.Context, query string, args ...interface{}) pgx.Row {
	row := pgc.conn.QueryRow(ctx, query, args...)
	return row
}

func (pgc postGresObject) Transaction() (pgx.Tx, error) {
	return pgc.conn.Begin(context.Background())
}

//Allows to insert multiple rows fast
//func (pgc postGresObject) CopyForm(tableName string, columnNames []string, rowSrc [][]interface{})
