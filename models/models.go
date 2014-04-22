package models

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var (
	orp *gorp.DbMap
)

type database struct {
}

func InitDb() (dbmap *gorp.DbMap, err error) {

	db, err := sql.Open("sqlite3", "./db/database.sqlite")
	if err != nil {
		return nil, err
	}

	orp = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	if err = orp.CreateTablesIfNotExists(); err != nil {
		return
	}
	dbmap = orp
	return
}
