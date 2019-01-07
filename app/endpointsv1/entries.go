package endpointsv1

import "fmt"

// Entries is the entries API struct definition
type Entries []Entry

// Entry is a singular SGV entry
type Entry struct {
	Device     string `json:"device"`
	Date       int64  `json:"date"`
	DateString string `json:"dateString"`
	Sgv        int    `json:"sgv"`
	Delta      int    `json:"delta"`
	Direction  string `json:"direction"`
	Type       string `json:"type"`
	Filtered   int    `json:"filtered"`
	Unfiltered int    `json:"unfiltered"`
	Rssi       int    `json:"rssi"`
	Noise      int    `json:"noise"`
	SysTime    string `json:"sysTime"`
}

// GenEntriesEndpoint is a placeholder which returns a fixed entries output
func (v1 *EndpointsV1) GenEntriesEndpoint() Entries {
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

// CSV returns a CSV output of entries
func (es Entries) CSV() [][]string {
	var out [][]string
	for _, e := range es {
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
