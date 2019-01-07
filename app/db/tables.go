package db

import (
	"math/rand"
	"time"

	"github.com/ttacon/glog"
)

// A Treatment is an entering of at least one of the following:
// carbs, insulin, glucose, or a note
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

// AddFakeTreatment adds a fake treatment
func (db *Db) AddFakeTreatment() {
	glog.FatalIf(db.dbMap.Insert(&Treatment{
		EnteredBy: "manual",
		Carbs:     float32(rand.Intn(50)),
		Insulin:   float32(rand.Intn(20)),
		Notes:     "test treatment",
		Time:      time.Now(),
	}))
}

// An Entry is a blood glucose entry derived from a sensor
// NOTE: Nightscout entries with type != "sgv" (e.g. mbg, cal) should be stored as a Treatment
// TODO: type=cal needs slope, intercept fields in new struct
type Entry struct {
	ID         int64     `db:"id, primarykey, autoincrement"`
	Device     string    `db:"device"`
	Time       time.Time `db:"time"`
	Sgv        int       `db:"sgv"`
	Delta      float32   `db:"delta"`
	Direction  string    `db:"direction"`
	Type       string    `db:"type"`
	Filtered   int       `db:"filtered"`
	Unfiltered int       `db:"unfiltered"`
	Rssi       int       `db:"rssi"`
	Noise      int       `db:"noise"`
}

func (db *Db) addTables() {
	db.dbMap.AddTableWithName(Treatment{}, "treatments").SetKeys(true, "ID")
	db.dbMap.AddTableWithName(Entry{}, "entries").SetKeys(true, "ID")
}

// AddFakeEntry adds a fake entry
func (db *Db) AddFakeEntry() {
	glog.FatalIf(db.dbMap.Insert(&Entry{
		Device:     "manual",
		Time:       time.Now(),
		Sgv:        100 + rand.Intn(100),
		Delta:      float32(rand.Intn(20)),
		Direction:  "FortyFiveUp",
		Type:       "sgv",
		Filtered:   91596,
		Unfiltered: 91154,
		Rssi:       100,
		Noise:      1,
	}))
}
