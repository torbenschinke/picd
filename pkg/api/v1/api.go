package v1

import "time"

type Resolution struct {
	X int
	Y int
}

type ISO string

type CameraID string

type Capabilities struct {
	Resolution []Resolution `json:"resolution"`
	ISO        []ISO        `json:"iso"`
}

type StillCamera struct {
	ID           CameraID     `json:"id"`
	Capabilities Capabilities `json:"capabilities"`
}

type Settings struct {
	Camera     CameraID      `json:"id"`
	Resolution Resolution    `json:"resolution"`
	ISO        ISO           `json:"iso"`
	Shutter    time.Duration `json:"shutter"`
	Exposure   string        `json:"exposure"` // e.g. off, auto, night etc.
	Mode       string        `json:"mode"`     // not really defined, force sensor mode e.g. 0 = auto
	Rotation   int           `json:"rotation"` //0, 90, 180, or 270
	Quality    int           `json:"quality"`
}
