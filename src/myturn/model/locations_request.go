package model

type LocationsRequest struct {
	Location    Position `json:"location"`
	FromDate    string   `json:"fromDate"`
	VaccineData string   `json:"vaccineData"`
	Url         string   `json:"url"`
}
