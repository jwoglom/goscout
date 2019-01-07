package db

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ttacon/glog"
)

// Db provides an abstraction for database methods
type Db struct {
	dbMap *gorp.DbMap
}

// NewDb creates a new Db object, which connects to the database and initializes tables
func NewDb() *Db {
	db := &Db{
		dbMap: newDbMap(),
	}

	glog.Infoln("adding tables")
	db.addTables()
	db.dbMap.CreateTablesIfNotExists()

	return db
}

func newDbMap() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	glog.FatalIf(err)

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.SqliteDialect{},
	}

	return dbmap

}
