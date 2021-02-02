package model

type CheckLocationRequest struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	VaccineData string `json:"vaccineData"`
	DoseNumber int `json:"doseNumber"`
	Url string `json:"url"`
}

