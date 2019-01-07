package app

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
)

// EndpointFunc refers to a function accepting a http request and returning
// a struct of data which is rendered via, e.g., json
type EndpointFunc func(*http.Request) interface{}

// CSVEndpointFunc refers to a function accepting a http request and returning
// an array of CSV data
type CSVEndpointFunc func(*http.Request) [][]string

func (s *Server) addRoutes() {
	v1 := s.Package.EndpointsV1
	s.router.HandleFunc("/api/v1/status", jsonWrapper(v1.GenStatusEndpoint))
	s.router.HandleFunc("/api/v1/status.json", jsonWrapper(v1.GenStatusEndpoint))

	s.router.HandleFunc("/api/v1/entries", csvWrapper(v1.GenEntriesCSVEndpoint))
	s.router.HandleFunc("/api/v1/entries.csv", csvWrapper(v1.GenEntriesCSVEndpoint))
	s.router.HandleFunc("/api/v1/entries.json", jsonWrapper(v1.GenStatusEndpoint))

	s.router.HandleFunc("/api/v1/treatments", jsonWrapper(v1.GenTreatmentsEndpoint))
	s.router.HandleFunc("/api/v1/treatments.json", jsonWrapper(v1.GenTreatmentsEndpoint))

	s.router.HandleFunc("/api/v1/devicestatus", jsonWrapper(v1.GenDeviceStatusEndpoint))
	s.router.HandleFunc("/api/v1/devicestatus.json", jsonWrapper(v1.GenDeviceStatusEndpoint))
}

func jsonWrapper(endpoint EndpointFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(endpoint(r))
	}
}

func csvWrapper(endpoint CSVEndpointFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cw := csv.NewWriter(w)
		defer cw.Flush()
		for _, row := range endpoint(r) {
			cw.Write(row)
		}
	}
}
