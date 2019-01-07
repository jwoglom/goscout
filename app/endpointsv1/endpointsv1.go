package endpointsv1

import "../db"

// EndpointsV1 allows the endpoints to access the database and
// other parts of the application
type EndpointsV1 struct {
	Db *db.Db
}
