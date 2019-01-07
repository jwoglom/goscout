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
}
