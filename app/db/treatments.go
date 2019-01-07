package db

// Fields is the db columns used for treatments
const treatmentFields = `id, enteredBy, carbs, insulin, glucose, notes, eventType, time`

// GetAll returns all treatments
func (db *Db) GetAllTreatments() []Treatment {
	var out []Treatment
	db.dbMap.Select(&out, `SELECT `+treatmentFields+` FROM treatments`)

	return out
}

// GetAll returns all treatments
func (db *Db) GetTreatments(limit int) []Treatment {
	var out []Treatment
	db.dbMap.Select(&out, `SELECT `+treatmentFields+` FROM treatments ORDER BY time DESC LIMIT :limit`,
		map[string]interface{}{
			"limit": limit,
		})

	return out
}

// GetFromID returns a single treatment with the given ID
func (db *Db) GetTreatmentWithID(id int) *Treatment {
	var out *Treatment
	db.dbMap.SelectOne(&out, `SELECT `+treatmentFields+` FROM treatments WHERE id = :id`,
		map[string]interface{}{
			"id": id,
		})

	return out
}
