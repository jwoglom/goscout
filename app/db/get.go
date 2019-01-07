package db

// GetTreatments returns all treatments
func (db *Db) GetTreatments() []Treatment {
	var out []Treatment
	db.dbMap.Select(&out, `SELECT
		id, enteredBy, carbs, insulin, glucose, notes, eventType, time
	FROM treatments`)

	return out
}

// GetTreatment returns a single treatment with the given ID
func (db *Db) GetTreatment(id int) *Treatment {
	var out *Treatment
	db.dbMap.SelectOne(&out, `SELECT
		id, enteredBy, carbs, insulin, glucose, notes, eventType, time
	FROM treatments WHERE id = :id`, map[string]interface{}{
		"id": id,
	})

	return out
}
