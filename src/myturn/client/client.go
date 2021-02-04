package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foo.com/vaccine-finder/model"
	"log"
	"net/http"
	"time"
)

func GetEligibility() model.EligibilityResponse {
	url := "https://api.myturn.ca.gov/public/eligibility"
	data := []byte(`{"eligibilityQuestionResponse":[{"id":"q.screening.18.yr.of.age","value":["q.screening.18.yr.of.age"],"type":"multi-select"},{"id":"q.screening.health.data","value":["q.screening.health.data"],"type":"multi-select"},{"id":"q.screening.eligibility.county","value":"San Diego","type":"single-select"},{"id":"q.screening.healthworker","value":"No","type":"single-select"},{"id":"q.screening.eligibility.age.range","value":"75 and older","type":"single-select"}],"url":"https://myturn.ca.gov/screening"}`)


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}


	var eligibility model.EligibilityResponse
	err = json.NewDecoder(resp.Body).Decode(&eligibility)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return eligibility
}

func GetLocations(VaccineData string, Lat float32, Lon float32) model.LocationsResponse {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	url := "https://api.myturn.ca.gov/public/locations/search"
	request := model.LocationsRequest{
		Location:    model.Position{
			Latitude:  Lat,
			Longitude: Lon,
		},
		FromDate:    time.Now().In(loc).Format("2006-01-02"),
		VaccineData: VaccineData,
		Url:         "https://myturn.ca.gov/location-select",
	}

	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error marshalling JSON. ", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}

	var locations model.LocationsResponse
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return locations
}

func checkLocation(VaccineData string, ExtId string, DoseNumber int, NumberOfDays time.Duration) model.CheckLocationResponse {
	url := fmt.Sprintf("https://api.myturn.ca.gov/public/locations/%s/availability", ExtId)

	loc, _ := time.LoadLocation("America/Los_Angeles")
	request := model.CheckLocationRequest{
		StartDate:   time.Now().In(loc).Format("2006-01-02"),
		EndDate:     time.Now().In(loc).Add(time.Hour*24*NumberOfDays).Format("2006-01-02"),
		VaccineData: VaccineData,
		DoseNumber:  DoseNumber,
		Url:         "https://myturn.ca.gov/appointment-select",
	}

	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error marshalling JSON. ", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}

	var location model.CheckLocationResponse
	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return location
}


