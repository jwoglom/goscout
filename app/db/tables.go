package db

import (
	"math/rand"
	"time"

	"github.com/ttacon/glog"
)

// A Treatment is an optional entering of carbs/insulin/glucose
type Treatment struct {
	ID          int64     `db:"id, primarykey, autoincrement"`
	EnteredBy   string    `db:"enteredBy"`
	Carbs       float32   `db:"carbs"`
	Insulin     float32   `db:"insulin"`
	Glucose     int       `db:"glucose"`
	GlucoseType string    `db:"glucoseType"`
	Notes       string    `db:"notes"`
	EventType   string    `db:"eventType"`
	Time        time.Time `db:"time"`
}

// EventTypes are the allowed event types
var EventTypes = []string{
	"",
	"Site Change",
	"Sensor Change",
}

func (db *Db) addTables() {
	db.dbMap.AddTableWithName(Treatment{}, "treatments").SetKeys(true, "ID")
}
func (db *Db) AddFakeTreatment() {
	glog.FatalIf(db.dbMap.Insert(&Treatment{
		EnteredBy: "manual",
		Carbs:     float32(rand.Intn(50)),
		Insulin:   float32(rand.Intn(20)),
		Notes:     "test treatment",
		Time:      time.Now(),
	}))
}
