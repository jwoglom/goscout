package endpoints

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status  string `json:"status"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GenStatusEndpoint() Status {
	return Status{
		Status:  "ok",
		Name:    "Goscout",
		Version: "0",
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GenStatusEndpoint())
}
