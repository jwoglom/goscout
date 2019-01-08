package db

import (
	"github.com/ttacon/glog"
)

// treatmentFields are the db columns used for treatments
const treatmentFields = `id, enteredBy, carbs, insulin, glucose, notes, eventType, time`

// treatmentFieldMap contains a mapping representing aliased fields, and its
// keys represent the only allowed filters
var treatmentFieldMap = map[string]interface{}{
	"id":         nil,
	"enteredBy":  nil,
	"carbs":      nil,
	"insulin":    nil,
	"glucose":    nil,
	"notes":      nil,
	"eventType":  nil,
	"time":       nil,
	"_id":        "id",
	"created_at": "time",
}

// GetAllTreatments returns all treatments in the database
func (db *Db) GetAllTreatments() []Treatment {
	var out []Treatment
	db.dbMap.Select(&out, `SELECT `+treatmentFields+` FROM treatments`)

	return out
}

// GetTreatments returns the limit most recent treatments
func (db *Db) GetTreatments(limit int) []Treatment {
	var out []Treatment
	_, err := db.dbMap.Select(&out, `SELECT `+treatmentFields+` FROM treatments ORDER BY time DESC LIMIT :limit`,
		map[string]interface{}{
			"limit": limit,
		})
	glog.FatalIf(err)
	return out
}

// GetTreatmentsWithFind returns limit treatments with the given treatment find operators
func (db *Db) GetTreatmentsWithFind(finds FindArguments, limit int) []Treatment {
	var out []Treatment
	query, args := finds.BuildQueryArgs(`SELECT `+treatmentFields+` FROM treatments`, limit, treatmentFieldMap)

	_, err := db.dbMap.Select(&out, query, args)
	glog.FatalIf(err)
	return out
}

// GetTreatmentWithID returns a single treatment with the given ID
func (db *Db) GetTreatmentWithID(id int) *Treatment {
	var out *Treatment
	err := db.dbMap.SelectOne(&out, `SELECT `+treatmentFields+` FROM treatments WHERE id = :id`,
		map[string]interface{}{
			"id": id,
		})
	glog.FatalIf(err)
	return out
}
