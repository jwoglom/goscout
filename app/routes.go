package app

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"./endpointsv1"
)

// Endpoint is a placeholder interface for API endpoint structs
type Endpoint interface{}

// CSVEndpoint is an interface for API endpoint structs which support CSV output
type CSVEndpoint interface {
	CSV() [][]string
}

func (s *Server) addRoutes() {
	s.addJSONRoute("api/v1/status", endpointsv1.GenStatusEndpoint())
	s.addCSVRoute("api/v1/entries", endpointsv1.GenEntriesEndpoint())
	s.addJSONRoute("api/v1/treatments", endpointsv1.GenTreatmentsEndpoint())
	s.addJSONRoute("api/v1/devicestatus", endpointsv1.GenDeviceStatusEndpoint())
}

// addJSONRoute adds a route which has a default JSON output
func (s *Server) addJSONRoute(name string, generator Endpoint) {
	s.router.HandleFunc(fmt.Sprintf("/%s", name), jsonWrapper(generator))
	s.router.HandleFunc(fmt.Sprintf("/%s.json", name), jsonWrapper(generator))
}

// addCSVRoute adds a route which has a default CSV output
func (s *Server) addCSVRoute(name string, generator CSVEndpoint) {
	s.router.HandleFunc(fmt.Sprintf("/%s", name), csvWrapper(generator))
	s.router.HandleFunc(fmt.Sprintf("/%s.csv", name), csvWrapper(generator))
	s.router.HandleFunc(fmt.Sprintf("/%s.json", name), jsonWrapper(generator))
}

/*func httpWrapper(e Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, e.Http())
	}
}*/

func jsonWrapper(e Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(e)
	}
}

func csvWrapper(e CSVEndpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cw := csv.NewWriter(w)
		defer cw.Flush()
		for _, row := range e.CSV() {
			cw.Write(row)
		}
	}
}
