package models

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	"github.com/jijeshmohan/mtgrid/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	orp *gorp.DbMap
)

func InitDb(c *config.Config) (dbmap *gorp.DbMap, err error) {
	var db *sql.DB

	if c.DbType == "sqlite3" {
		db, err = sql.Open("sqlite3", c.DbPath)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unsupported database type")
	}

	orp = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	addDeviceTable()

	if err = orp.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	dbmap = orp
	orp.TraceOn("SQL", log.New(os.Stdout, "SQL:", log.Lmicroseconds))
	return
}
