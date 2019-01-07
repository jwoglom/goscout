package endpointsv1

// DeviceStatus is the devicestatus API struct definition
type DeviceStatus []DevStatus

// DevStatus is a singular status of a device
type DevStatus struct {
	Device    string      `json:"device"`
	Uploader  DevUploader `json:"uploader"`
	CreatedAt string      `json:"created_at"`
}

// DevUploader provides details about the uploader for a DevStatus
type DevUploader struct {
	Battery int `json:"battery"`
}

// GenDeviceStatusEndpoint is a placeholder which returns a fixed devicestatus output
func GenDeviceStatusEndpoint() DeviceStatus {
	return DeviceStatus{{
		Device: "Google Pixel 2 XL",
		Uploader: DevUploader{
			Battery: 91,
		},
		CreatedAt: "2019-01-07T01:00:02.336Z",
	}, {
		Device: "Google Pixel 2 XL",
		Uploader: DevUploader{
			Battery: 90,
		},
		CreatedAt: "2019-01-07T00:55:02.424Z",
	}}
}
