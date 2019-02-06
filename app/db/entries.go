package db

import (
	"time"

	"github.com/ttacon/glog"
)

// entriesFields are the db columns used for entries
const entriesFields = `id, device, time, sgv, delta, direction, type, filtered, unfiltered, rssi, noise`

// entriesFieldMap contains a mapping representing aliased fields, and its
// keys represent the only allowed filters
var entriesFieldMap = map[string]interface{}{
	"id":         nil,
	"device":     nil,
	"time":       nil,
	"sgv":        nil,
	"delta":      nil,
	"direction":  nil,
	"type":       nil,
	"filtered":   nil,
	"unfiltered": nil,
	"rssi":       nil,
	"noise":      nil,
	"_id":        "id",
	"created_at": "time",
	"sysTime":    "time",
}

// GetAllEntries returns all entries in the database
func (db *Db) GetAllEntries() []Entry {
	var out []Entry
	db.dbMap.Select(&out, `SELECT `+entriesFields+` FROM entries`)

	return out
}

// GetEntries returns the limit most recent entries
func (db *Db) GetEntries(limit int) []Entry {
	var out []Entry
	_, err := db.dbMap.Select(&out, `SELECT `+entriesFields+` FROM entries ORDER BY time DESC LIMIT :limit`,
		map[string]interface{}{
			"limit": limit,
		})
	glog.FatalIf(err)
	return out
}

// GetEntriesWithFind returns limit treatments with the given entries find operators
func (db *Db) GetEntriesWithFind(finds FindArguments, limit int) []Entry {
	var out []Entry
	query, args := finds.BuildQueryArgs(`SELECT `+entriesFields+` FROM entries`, limit, entriesFieldMap)

	_, err := db.dbMap.Select(&out, query, args)
	glog.FatalIf(err)
	return out
}

// GetEntryWithID returns a single entry with the given ID
func (db *Db) GetEntryWithID(id int) *Entry {
	var out *Entry
	err := db.dbMap.SelectOne(&out, `SELECT `+entriesFields+` FROM entries WHERE id = :id`,
		map[string]interface{}{
			"id": id,
		})
	glog.FatalIf(err)
	return out
}

// GetEntryWithTimeAndType returns a single entry with the given time and type, if exists
// Used for upsert operations
func (db *Db) GetEntryWithTimeAndType(time time.Time, typ string) *Entry {
	var out []Entry
	_, err := db.dbMap.Select(&out, `SELECT `+entriesFields+` FROM entries WHERE time = :time AND type = :type LIMIT 1`,
		map[string]interface{}{
			"time": time,
			"type": typ,
		})
	glog.FatalIf(err)
	if len(out) > 0 {
		return &out[0]
	}
	return nil
}

// UpsertEntry either inserts the entry or, if an entry already exists with the same
// Time and Type, replaces it entirely with the given entry
func (db *Db) UpsertEntry(entry Entry) (int64, error) {
	exists := db.GetEntryWithTimeAndType(entry.Time, entry.Type)
	glog.Infoln("upsert exists: ", exists)
	if exists == nil {
		return db.insertEntry(&entry)
	} else {
		entry.ID = exists.ID
		_, err := db.dbMap.Update(&entry)
		return entry.ID, err
	}
}

func (db *Db) insertEntry(entry *Entry) (int64, error) {
	err := db.dbMap.Insert(entry)
	return entry.ID, err
}
