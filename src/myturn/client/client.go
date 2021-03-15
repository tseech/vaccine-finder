package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foo.com/vaccine-finder/model"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var jar, _ = cookiejar.New(nil)

func GetEligibility() model.EligibilityResponse {

	client := &http.Client{
		Jar: jar,
	}

	{
		getReq, _ := http.NewRequest("GET", "https://myturn.ca.gov", nil)
		getReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15")
		getReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		getResp, _ := client.Do(getReq)
		fmt.Println("Get 1: " + getResp.Status)
	}

	{
		getReq, _ := http.NewRequest("GET", "https://api.myturn.ca.gov/public/config?url=https:%2F%2Fmyturn.ca.gov%2F", nil)
		getReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15")
		getReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		getResp, _ := client.Do(getReq)
		fmt.Println("Get 2: " + getResp.Status)
	}

	url := "https://api.myturn.ca.gov/public/eligibility"
	data := []byte(`{"eligibilityQuestionResponse":[{"id":"q.screening.18.yr.of.age","value":["q.screening.18.yr.of.age"],"type":"multi-select"},{"id":"q.screening.health.data","value":["q.screening.health.data"],"type":"multi-select"},{"id":"q.screening.privacy.statement","value":["q.screening.privacy.statement"],"type":"multi-select"},{"id":"q.screening.eligibility.age.range","value":"75 and older","type":"single-select"},{"id":"q.screening.underlying.health.condition","value":"No","type":"single-select"},{"id":"q.screening.disability","value":"No","type":"single-select"},{"id":"q.screening.eligibility.industry","value":"Other","type":"single-select"},{"id":"q.screening.eligibility.county","value":"San Diego","type":"single-select"},{"id":"q.screening.accessibility.code","type":"text"}],"url":"https://myturn.ca.gov/screening"}`)

	postReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	postReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15")
	postReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	postReq.Header.Add("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(postReq)
	fmt.Println("POST: " + resp.Status)
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
		Location: model.Position{
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

	client := &http.Client{
		Jar: jar,
	}
	postReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	postReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15")
	postReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	postReq.Header.Add("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(postReq)
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
		EndDate:     time.Now().In(loc).Add(time.Hour * 24 * NumberOfDays).Format("2006-01-02"),
		VaccineData: VaccineData,
		DoseNumber:  DoseNumber,
		Url:         "https://myturn.ca.gov/appointment-select",
	}

	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error marshalling JSON. ", err)
	}

	client := &http.Client{
		Jar: jar,
	}
	postReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	postReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15")
	postReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	postReq.Header.Add("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(postReq)
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
