package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	eligibility := getEligibility()
	if !eligibility.Eligible {
		fmt.Println("Not eligible")
		return
	}
	fmt.Println("Eligibility confirmed")

	fmt.Println("Finding locations...")
	locations := getLocations(eligibility.VaccineData)
	fmt.Printf("%d locations found\n", len(locations.Locations))

	numberOfDays := time.Duration(30)

	os.RemoveAll("./dist")
	os.Mkdir("./dist", 0777)
	file, err := os.Create("./dist/index.md")
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	loc, _ := time.LoadLocation("America/Los_Angeles")

	writer.WriteString("# San Diego Vaccine Appointments\n")
	writer.WriteString(fmt.Sprintf("*Last Updated: %s*\n\n", time.Now().In(loc).Format("Mon Jan 2 15:04:05 MST 2006")))
	writer.WriteString(fmt.Sprintf("*Date range: %s - %s*\n",
		time.Now().Format("Mon Jan 2 2006"),
		time.Now().Add(time.Hour*24*numberOfDays).Format("Mon Jan 2 2006")))
	writer.WriteString("\n")

	fmt.Printf("Checking the next %d days...\n", numberOfDays)
	for _, element := range locations.Locations  {
		dose1 := checkLocation(locations.VaccineData, element.ExtId, 1, numberOfDays)
		dose2 := checkLocation(locations.VaccineData, element.ExtId, 2, numberOfDays)

		hasDose1 := 0
		dose1Dates := ""
		hasDose2 := 0
		dose2Dates := ""
		for i:=0; i<len(dose1.Availability) && i<len(dose2.Availability); i++ {
			if dose1.Availability[i].Available {
				hasDose1++
				if dose1Dates != "" {
					dose1Dates += ", "
				}
				dose1Dates += dose1.Availability[i].Date
			}
			if dose2.Availability[i].Available {
				hasDose2++
				if dose2Dates != "" {
					dose2Dates += ", "
				}
				dose2Dates += dose2.Availability[i].Date
			}
		}

		doseStatus := ""
		if hasDose1 > 0 && hasDose2 == 0 {
			doseStatus = "Dose 1 only"
		} else if hasDose1 == 0 && hasDose2 > 0 {
			doseStatus = "Dose 2 only"
		} else if hasDose1 > 0 && hasDose2 > 0 {
			doseStatus = "Dose 1 & 2"
		} else {
			doseStatus = "No doses"
		}

		writer.WriteString(fmt.Sprintf("## *%s* - %s\n", doseStatus, element.Name))
		writer.WriteString(fmt.Sprintf("### %s\n", element.DisplayAddress))
		writer.WriteString(fmt.Sprintf("- Done 1 available on %d days\n", hasDose1))
		if dose1Dates != "" {
			writer.WriteString(fmt.Sprintf("  - Days: %s\n", dose1Dates))
		}
		writer.WriteString(fmt.Sprintf("- Done 2 available on %d days\n", hasDose2))
		if dose2Dates != "" {
			writer.WriteString(fmt.Sprintf("  - Days: %s\n", dose2Dates))
		}

		writer.WriteString("\n")
	}

	writer.Flush()
}

func getEligibility() EligibilityResponse {
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

