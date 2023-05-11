package store

import (
	"context"
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

func (db *Storage) Selectx(ctx context.Context, dest interface{}, sqlizer sq.Sqlizer) error {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	return db.SelectContext(ctx, dest, stmt, args...)
}

func (db *Storage) Getx(ctx context.Context, dest interface{}, sqlizer sq.Sqlizer) error {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	return db.GetContext(ctx, dest, stmt, args...)
}

func (db *Storage) Exec(ctx context.Context, sqlizer sq.Sqlizer) error {
	stmt, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}
