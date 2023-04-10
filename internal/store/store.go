package store

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	*sqlx.DB
}

func NewStorage() Storage {
	return Storage{sqlx.NewDb(getPostgres(), "psx")}
}

func getPostgres() *sql.DB {
	connStr := "dbname=e_university sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("cant parse config" + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("can`t ping db" + err.Error())
	}

	db.SetMaxOpenConns(10)

	return db
}

// Builder вернет squirrel SQL Builder обьект
func (s *Storage) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
