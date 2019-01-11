package app

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

// EndpointFunc refers to a function accepting a http request and returning
// a struct of data which is rendered via, e.g., json
type EndpointFunc func(*http.Request) interface{}

// CSVEndpointFunc refers to a function accepting a http request and returning
// an array of CSV data
type CSVEndpointFunc func(*http.Request) [][]string

func (s *Server) addRoutes() {
	s.router.HandleFunc("/", indexPage)

	v1 := s.Package.EndpointsV1
	s.router.HandleFunc("/api/v1/status", v1.GenStatusHTMLEndpoint)
	s.router.HandleFunc("/api/v1/status.json", jsonWrapper(v1.GenStatusEndpoint))

	s.router.HandleFunc("/api/v1/entries", csvWrapper(v1.GenEntriesCSVEndpoint))
	s.router.HandleFunc("/api/v1/entries.csv", csvWrapper(v1.GenEntriesCSVEndpoint))
	s.router.HandleFunc("/api/v1/entries.json", jsonWrapper(v1.GenEntriesEndpoint))

	s.router.HandleFunc("/api/v1/entries/{type}.csv", csvWrapper(v1.GenEntriesCSVEndpoint))
	s.router.HandleFunc("/api/v1/entries/{type}.json", jsonWrapper(v1.GenEntriesEndpoint))
	s.router.HandleFunc("/api/v1/entries/{type}", csvWrapper(v1.GenEntriesCSVEndpoint))

	s.router.HandleFunc("/api/v1/treatments", jsonWrapper(v1.GenTreatmentsEndpoint))
	s.router.HandleFunc("/api/v1/treatments.json", jsonWrapper(v1.GenTreatmentsEndpoint))

	s.router.HandleFunc("/api/v1/devicestatus", jsonWrapper(v1.GenDeviceStatusEndpoint))
	s.router.HandleFunc("/api/v1/devicestatus.json", jsonWrapper(v1.GenDeviceStatusEndpoint))
}

func jsonWrapper(endpoint EndpointFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		json.NewEncoder(w).Encode(endpoint(r))
	}
}

func csvWrapper(endpoint CSVEndpointFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		cw := csv.NewWriter(w)
		defer cw.Flush()
		for _, row := range endpoint(r) {
			cw.Write(row)
		}
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
<body>
<ul>
	<li><a href="/api/v1/status">api/v1/status.json</a></li>
	<li><a href="/api/v1/entries">api/v1/entries.csv</a></li>
	<li><a href="/api/v1/entries.json">api/v1/entries.json</a></li>
	<li><a href="/api/v1/treatments.json">api/v1/treatments.json</a></li>
	<li><a href="/api/v1/devicestatus.json">api/v1/devicestatus.json</a></li>
</ul>
</body>
</html>`)
}
