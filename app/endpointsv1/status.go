package endpointsv1

// Status contains information about this Goscout instance
type Status struct {
	Status  string `json:"status"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// GenStatusEndpoint is a placeholder which returns the fixed status output
func GenStatusEndpoint() Status {
	return Status{
		Status:  "ok",
		Name:    "Goscout",
		Version: "0",
	}
}
