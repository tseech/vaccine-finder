package main

import (
	"fmt"
	"foo.com/vaccine-finder/client"
	"time"
)

func main() {
	fmt.Println("Checking eligibility...")
	eligibility := client.GetEligibility()
	if !eligibility.Eligible {
		fmt.Println("Not eligible")
		return
	}
	fmt.Println("Eligibility confirmed")

	numberOfDays := time.Duration(30)

	fmt.Println("Finding locations...")
	southLocations := client.GetLocations(eligibility.VaccineData, 32.8250787, -117.091176)
	northLocations := client.GetLocations(eligibility.VaccineData, 33.1433723, -117.1661449)
	byDose := client.GetAvailabilityByDose(numberOfDays, southLocations, northLocations)


	client.BuildMarkdownFile(byDose, numberOfDays)
}

