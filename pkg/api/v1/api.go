package v1

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
	Camera     CameraID   `json:"id"`
	Resolution Resolution `json:"resolution"`
	ISO        ISO        `json:"iso"`
}
