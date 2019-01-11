package endpointsv1

import (
	"fmt"
	"net/http"
	"time"
)

// Status contains information about this Goscout instance
type Status struct {
	Status            string                  `json:"status"`
	Name              string                  `json:"name"`
	Version           string                  `json:"version"`
	ServerTime        string                  `json:"serverTime"`
	ServerTimeEpoch   int64                   `json:"serverTimeEpoch"`
	APIEnabled        bool                    `json:"apiEnabled"`
	CareportalEnabled bool                    `json:"careportalEnabled"`
	BoluscalcEnabled  bool                    `json:"boluscalcEnabled"`
	Settings          *StatusSettings         `json:"settings"`
	ExtendedSettings  *StatusExtendedSettings `json:"extendedSettings"`
	Authorized        *bool                   `json:"authorized"`
}

type StatusSettings struct {
	Units                  string `json:"units"`
	TimeFormat             int    `json:"timeFormat"`
	NightMode              bool   `json:"nightMode"`
	EditMode               bool   `json:"editMode"`
	ShowRawbg              string `json:"showRawbg"`
	CustomTitle            string `json:"customTitle"`
	Theme                  string `json:"theme"`
	AlarmUrgentHigh        bool   `json:"alarmUrgentHigh"`
	AlarmUrgentHighMins    []int  `json:"alarmUrgentHighMins"`
	AlarmHigh              bool   `json:"alarmHigh"`
	AlarmHighMins          []int  `json:"alarmHighMins"`
	AlarmLow               bool   `json:"alarmLow"`
	AlarmLowMins           []int  `json:"alarmLowMins"`
	AlarmUrgentLow         bool   `json:"alarmUrgentLow"`
	AlarmUrgentLowMins     []int  `json:"alarmUrgentLowMins"`
	AlarmUrgentMins        []int  `json:"alarmUrgentMins"`
	AlarmWarnMins          []int  `json:"alarmWarnMins"`
	AlarmTimeagoWarn       bool   `json:"alarmTimeagoWarn"`
	AlarmTimeagoWarnMins   int    `json:"alarmTimeagoWarnMins"`
	AlarmTimeagoUrgent     bool   `json:"alarmTimeagoUrgent"`
	AlarmTimeagoUrgentMins int    `json:"alarmTimeagoUrgentMins"`
	AlarmPumpBatteryLow    bool   `json:"alarmPumpBatteryLow"`
	Language               string `json:"language"`
	ScaleY                 string `json:"scaleY"`
	ShowPlugins            string `json:"showPlugins"`
	ShowForecast           string `json:"showForecast"`
	FocusHours             int    `json:"focusHours"`
	Heartbeat              int    `json:"heartbeat"`
	BaseURL                string `json:"baseURL"`
	AuthDefaultRoles       string `json:"authDefaultRoles"`
	Thresholds             struct {
		BgHigh         int `json:"bgHigh"`
		BgTargetTop    int `json:"bgTargetTop"`
		BgTargetBottom int `json:"bgTargetBottom"`
		BgLow          int `json:"bgLow"`
	} `json:"thresholds"`
	DefaultFeatures []string `json:"DEFAULT_FEATURES"`
	AlarmTypes      []string `json:"alarmTypes"`
	Enable          []string `json:"enable"`
}

type StatusExtendedSettings struct {
	DeviceStatus struct {
		Advanced bool `json:"advanced"`
	} `json:"devicestatus"`
}

// GenStatusEndpoint is a placeholder which returns the fixed status output
func (v1 *EndpointsV1) GenStatusEndpoint(r *http.Request) interface{} {
	return Status{
		Status:            "ok",
		Name:              "Goscout",
		Version:           "0",
		ServerTime:        time.Now().String(),
		ServerTimeEpoch:   time.Now().UnixNano() / int64(time.Millisecond),
		APIEnabled:        true,
		CareportalEnabled: true,
		BoluscalcEnabled:  true,
		Settings:          &StatusSettings{},
		ExtendedSettings:  &StatusExtendedSettings{},
		Authorized:        nil,
	}
}

// GenStatusHTMLEndpoint is a placeholder which returns the fixed status output
func (v1 *EndpointsV1) GenStatusHTMLEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>STATUS OK</h1>`)
}
