package endpointsv1

import (
	"fmt"

	"../db"
)

// Treatments is the treatments API struct definition
type Treatments []Treatment

// Treatment is a singular treatment
type Treatment struct {
	ID          string  `json:"_id"`
	EventType   string  `json:"eventType"`
	Insulin     float32 `json:"insulin"`
	Carbs       float32 `json:"carbs"`
	Glucose     int     `json:"glucose"`
	GlucoseType string  `json:"glucoseType"`
	EnteredBy   string  `json:"enteredBy"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"created_at"`
}

// GenTreatmentsEndpoint is a placeholder which returns a fixed treatments output
func (v1 *EndpointsV1) GenTreatmentsEndpoint() Treatments {
	/*
		dummyData := Treatments{{
			EventType: "Meal Bolus",
			Insulin:   6.67,
			Carbs:     40,
			EnteredBy: "Diabetes-M (dm2nsc)",
			Notes:     "wrap (40)",
			CreatedAt: "2019-01-06 18:05:00-05:00",
		}, {
			EventType:   "<none>",
			EnteredBy:   "xdrip",
			Notes:       "note",
			CreatedAt:   "2019-01-03T07:22:10Z",
			Carbs:       0,
			Insulin:     0,
			Glucose:     100,
			GlucoseType: "finger",
		}}
	*/
	var out Treatments
	for _, tr := range v1.Db.GetTreatments() {
		out = append(out, DbTreatmentToTreatment(tr))
	}
	return out
}

// DbTreatmentToTreatment converts a database to a local object
func DbTreatmentToTreatment(t db.Treatment) Treatment {
	return Treatment{
		ID:          fmt.Sprintf("%d", t.ID),
		EventType:   t.EventType,
		Insulin:     t.Insulin,
		Carbs:       t.Carbs,
		Glucose:     t.Glucose,
		GlucoseType: t.GlucoseType,
		EnteredBy:   t.EnteredBy,
		Notes:       t.Notes,
		CreatedAt:   t.Time.String(),
	}
}
