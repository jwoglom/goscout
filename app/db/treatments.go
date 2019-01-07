package db

import (
	"fmt"
	"strings"

	"github.com/ttacon/glog"
)

// Fields is the db columns used for treatments
const treatmentFields = `id, enteredBy, carbs, insulin, glucose, notes, eventType, time`

// GetAll returns all treatments
func (db *Db) GetAllTreatments() []Treatment {
	var out []Treatment
	db.dbMap.Select(&out, `SELECT `+treatmentFields+` FROM treatments`)

	return out
}

// GetTreatments returns limit treatments
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
	query, args := finds.BuildQueryArgs(`SELECT `+treatmentFields+` FROM treatments`, limit)
	glog.Infoln("getTreatmentsWithFind: ", query, args)
	// FIXME: SQL injection
	for k, v := range args {
		var vstr string
		if vn, ok := v.(string); ok {
			vstr = fmt.Sprintf(`"%s"`, vn)
		} else {
			vstr = fmt.Sprintf("%d", v.(int))
		}
		query = strings.Replace(query, ":"+k, vstr, -1)
		glog.Infoln("replace:", k, vstr, "out: ", query)
	}
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
