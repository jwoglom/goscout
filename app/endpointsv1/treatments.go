package endpointsv1

// Treatments is the treatments API struct definition
type Treatments []Treatment

// Treatment is a singular treatment
type Treatment struct {
	EventType string  `json:"eventType"`
	Insulin   float32 `json:"insulin"`
	Carbs     float32 `json:"carbs"`
	EnteredBy string  `json:"enteredBy"`
	Notes     string  `json:"notes"`
	CreatedAt string  `json:"created_at"`
}

// GenTreatmentsEndpoint is a placeholder which returns a fixed treatments output
func GenTreatmentsEndpoint() Treatments {
	return Treatments{{
		EventType: "Meal Bolus",
		Insulin:   6.67,
		Carbs:     40,
		EnteredBy: "Diabetes-M (dm2nsc)",
		Notes:     "wrap (40)",
		CreatedAt: "2019-01-06 18:05:00-05:00",
	}, {
		EventType: "<none>",
		EnteredBy: "xdrip",
		Notes:     "note",
		CreatedAt: "2019-01-03T07:22:10Z",
		Carbs:     0,
		Insulin:   0,
	}}
}
