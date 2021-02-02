package model

type LocationsResponse struct {
	Eligible bool       `json:"eligible"`
	VaccineData string  `json:"vaccineData"`
	Locations[]Location `json:"locations"`
}
