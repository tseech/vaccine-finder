package model

type CheckLocationResponse struct {
	Availability[]Availability `json:"availability"`
}
