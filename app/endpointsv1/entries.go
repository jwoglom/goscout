package endpointsv1

import (
	"fmt"
	"net/http"
	"time"

	"../db"
	"github.com/ttacon/glog"
)

// Entries is the entries API struct definition
type Entries []Entry

// Entry is a singular SGV entry
type Entry struct {
	ID         string  `json:"_id"`
	Device     string  `json:"device"`
	Date       int64   `json:"date"`
	DateString string  `json:"dateString"`
	Sgv        int     `json:"sgv"`
	Delta      float32 `json:"delta"`
	Direction  string  `json:"direction"`
	Type       string  `json:"type"`
	Filtered   int     `json:"filtered"`
	Unfiltered int     `json:"unfiltered"`
	Rssi       int     `json:"rssi"`
	Noise      int     `json:"noise"`
	SysTime    string  `json:"sysTime"`
}

// GenEntriesEndpointDummy is a placeholder which returns a fixed entries output
func (v1 *EndpointsV1) GenEntriesEndpointDummy(r *http.Request) interface{} {
	return Entries{{
		DateString: "2019-01-06T19:04:57.985-0500",
		Date:       1546819497985,
		Sgv:        157,
		Direction:  "FortyFiveUp",
		Device:     "xDrip-DexcomG5 G5 Native",
	}, {
		DateString: "2019-01-06T18:44:58.109-0500",
		Date:       1546818298109,
		Sgv:        131,
		Direction:  "Flat",
		Device:     "xDrip-DexcomG5 G5 Native",
	}}
}

// GenEntriesEndpoint returns all treatments in the database
func (v1 *EndpointsV1) GenEntriesEndpoint(r *http.Request) interface{} {
	var out Entries

	findArgs, count := db.FindArgumentsFromQuery(r.URL.Query())
	glog.Infoln("findArgs:", findArgs)
	for _, tr := range v1.Db.GetEntriesWithFind(findArgs, count) {
		out = append(out, DbEntryToEntry(tr))
	}
	return out
}

// DbEntryToEntry converts a database to a local object
func DbEntryToEntry(t db.Entry) Entry {
	return Entry{
		ID:         fmt.Sprintf("%d", t.ID),
		Device:     t.Device,
		Date:       t.Time.UnixNano() / int64(time.Millisecond),
		DateString: t.Time.String(),
		Sgv:        t.Sgv,
		Delta:      t.Delta,
		Direction:  t.Direction,
		Type:       t.Type,
		Filtered:   t.Filtered,
		Unfiltered: t.Unfiltered,
		Rssi:       t.Rssi,
		Noise:      t.Noise,
		SysTime:    t.Time.String(),
	}
}

// GenEntriesCSVEndpoint converts the output of GenEntriesEndpoint to CSV
func (v1 *EndpointsV1) GenEntriesCSVEndpoint(r *http.Request) [][]string {
	entries := v1.GenEntriesEndpoint(r).(Entries)
	var out [][]string
	for _, e := range entries {
		var row []string
		row = append(row, e.DateString)
		row = append(row, fmt.Sprintf("%d", e.Date))
		row = append(row, fmt.Sprintf("%d", e.Sgv))
		row = append(row, e.Direction)
		row = append(row, e.Device)
		out = append(out, row)
	}
	return out
}
