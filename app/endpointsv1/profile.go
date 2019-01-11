package endpointsv1

import (
	"net/http"
	"time"
)

// Profiles is the profile API struct definition
type Profiles []Profile

// Profile is an individual user profile
type Profile struct {
	ID             string `json:"_id"`
	DefaultProfile string `json:"defaultProfile"`
	Store          struct {
		Default struct {
			Dia       string `json:"dia"`
			Carbratio []struct {
				Time          string `json:"time"`
				Value         string `json:"value"`
				TimeAsSeconds string `json:"timeAsSeconds"`
			} `json:"carbratio"`
			CarbsHr string `json:"carbs_hr"`
			Delay   string `json:"delay"`
			Sens    []struct {
				Time          string `json:"time"`
				Value         string `json:"value"`
				TimeAsSeconds string `json:"timeAsSeconds"`
			} `json:"sens"`
			Timezone string `json:"timezone"`
			Basal    []struct {
				Time          string `json:"time"`
				Value         string `json:"value"`
				TimeAsSeconds string `json:"timeAsSeconds"`
			} `json:"basal"`
			TargetLow []struct {
				Time          string `json:"time"`
				Value         string `json:"value"`
				TimeAsSeconds string `json:"timeAsSeconds"`
			} `json:"target_low"`
			TargetHigh []struct {
				Time          string `json:"time"`
				Value         string `json:"value"`
				TimeAsSeconds string `json:"timeAsSeconds"`
			} `json:"target_high"`
			Units string `json:"units"`
		} `json:"Default"`
	} `json:"store"`
	StartDate time.Time `json:"startDate"`
	Mills     string    `json:"mills"`
	Units     string    `json:"units"`
	CreatedAt time.Time `json:"created_at"`
}

// GenProfileEndpoint is a stub returning all profiles
func (v1 *EndpointsV1) GenProfileEndpoint(r *http.Request) interface{} {
	out := Profiles{}

	return out
}

// GenCurrentProfileEndpoint is a stub returning the current profile
func (v1 *EndpointsV1) GenCurrentProfileEndpoint(r *http.Request) interface{} {
	var out Profile

	return out
}
