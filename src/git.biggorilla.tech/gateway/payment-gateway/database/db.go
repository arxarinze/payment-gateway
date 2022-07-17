package database

import (
	"database/sql"

	"golang.org/x/net/context"
)

type Database interface {
	ConnectDatabase(psqlInfo string) sql.DB
}

type database struct {
}

func NewDatabase(ctx context.Context) Database {
	return &database{}
}

func (r *database) ConnectDatabase(psqlInfo string) sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return *db
}
