package endpointsv1

import (
	"fmt"
	"net/http"

	"../db"
	"github.com/gorilla/mux"
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

// GenTreatmentsEndpoint returns all treatments in the database
func (v1 *EndpointsV1) GenTreatmentsEndpoint(r *http.Request) interface{} {
	var out Treatments

	findArgs, count := db.FindArgumentsFromQuery(r.URL.Query(), mux.Vars(r))
	for _, tr := range v1.Db.GetTreatmentsWithFind(findArgs, count) {
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
