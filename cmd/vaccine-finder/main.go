package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EligibilityResponse struct {
	Eligible bool `json:"eligible"`
	VaccineData string `json:"vaccineData"`
}

type LocationsRequest struct {
	Location Position `json:"location"`
	FromDate string `json:"fromDate"`
	VaccineData string `json:"vaccineData"`
	Url string `json:"url"`
}

type Position  struct {
	Latitude float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type Location struct {
	DisplayAddress string `json:"displayAddress"`
	ExtId string `json:"extId"`
	Name string `json:"name"`
}

type LocationsResponse struct {
	Eligible bool `json:"eligible"`
	VaccineData string `json:"vaccineData"`
	Locations[] Location `json:"locations"`
}

type CheckLocationRequest struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	VaccineData string `json:"vaccineData"`
	DoseNumber int `json:"doseNumber"`
	Url string `json:"url"`
}

type CheckLocationResponse struct {
	Availability[] Availability `json:"availability"`
}

type Availability struct {
	Date string `json:"date"`
	Available bool `json:"available"`
}


func main() {

	fmt.Println("Checking eligibility...")
	eligibility := getEligbility()
	if !eligibility.Eligible {
		fmt.Println("Not eligible")
		return
	}
	fmt.Println("Eligibility confirmed")

	fmt.Println("Finding locations...")
	locations := getLocations(eligibility.VaccineData)
	fmt.Printf("%d locations found\n", len(locations.Locations))

	numberOfDays := time.Duration(30)

	fmt.Printf("Checking the next %d days...\n", numberOfDays)
	for index, element := range locations.Locations  {
		fmt.Printf("%d - %s\n", index+1, element.Name)
		dose1 := checkLocation(locations.VaccineData, element.ExtId, 1, numberOfDays)
		dose2 := checkLocation(locations.VaccineData, element.ExtId, 2, numberOfDays)

		for i:=0; i<len(dose1.Availability) && i<len(dose2.Availability); i++ {
			if !dose1.Availability[i].Available && !dose2.Availability[i].Available {
				//fmt.Printf("   %s: None\n", dose1.Availability[i].Date)
			} else if dose1.Availability[i].Available && !dose2.Availability[i].Available {
				fmt.Printf("   %s: Done 1 only\n", dose1.Availability[i].Date)
			} else if !dose1.Availability[i].Available && dose2.Availability[i].Available {
				fmt.Printf("   %s: Done 2 only\n", dose1.Availability[i].Date)
			} else {
				fmt.Printf("   %s: Both!!!!!\n", dose1.Availability[i].Date)
			}
		}

		fmt.Println()
	}
}

func getEligbility() EligibilityResponse {
	url := "https://api.myturn.ca.gov/public/eligibility"
	data := []byte(`{"eligibilityQuestionResponse":[{"id":"q.screening.18.yr.of.age","value":["q.screening.18.yr.of.age"],"type":"multi-select"},{"id":"q.screening.health.data","value":["q.screening.health.data"],"type":"multi-select"},{"id":"q.screening.eligibility.county","value":"San Diego","type":"single-select"},{"id":"q.screening.healthworker","value":"No","type":"single-select"},{"id":"q.screening.eligibility.age.range","value":"65 - 74","type":"single-select"}],"url":"https://myturn.ca.gov/screening"}`)


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}


	var eligibility EligibilityResponse
	err = json.NewDecoder(resp.Body).Decode(&eligibility)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return eligibility
}

func getLocations(VaccineData string) LocationsResponse {
	url := "https://api.myturn.ca.gov/public/locations/search"
	request := LocationsRequest{
		Location:    Position{
			Latitude:  32.8250787,
			Longitude: -117.091176,
		},
		FromDate:    time.Now().Format("2006-01-02"),
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

	var locations LocationsResponse
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return locations
}

func checkLocation(VaccineData string, ExtId string, DoseNumber int, NumberOfDays time.Duration) CheckLocationResponse {
	url := fmt.Sprintf("https://api.myturn.ca.gov/public/locations/%s/availability", ExtId)

	request := CheckLocationRequest{
		StartDate:   time.Now().Format("2006-01-02"),
		EndDate:     time.Now().Add(time.Hour*24*NumberOfDays).Format("2006-01-02"),
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

	var location CheckLocationResponse
	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	return location
}

