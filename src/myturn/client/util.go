package client

import (
	"foo.com/vaccine-finder/model"
	"time"
)

type LocationAvailability struct {
	Location          model.Location
	Dose1Availability []model.Availability
	Dose2Availability []model.Availability
}

type LocationAvailabilityByDose struct {
	Both[] LocationAvailability
	Dose1Only[] LocationAvailability
	Dose2Only[] LocationAvailability
	Neither[] LocationAvailability
}

func GetAvailabilityByDose(NumberOfDaysSearched time.Duration, LocationsResponses...model.LocationsResponse) LocationAvailabilityByDose {
	locationMap := make(map[string]LocationAvailability)
	byDose := LocationAvailabilityByDose {
		Both: make([]LocationAvailability, 0, 100),
		Dose1Only: make([]LocationAvailability, 0, 100),
		Dose2Only: make([]LocationAvailability, 0, 100),
		Neither: make([]LocationAvailability, 0, 100),
	}

	for _, locationsResponse := range LocationsResponses {
		for _, location := range locationsResponse.Locations {
			if _, ok := locationMap[location.ExtId]; !ok {
				dose1 := checkLocation(locationsResponse.VaccineData, location.ExtId, 1, NumberOfDaysSearched)
				dose2 := checkLocation(locationsResponse.VaccineData, location.ExtId, 2, NumberOfDaysSearched)

				 locationAvailability := LocationAvailability{
					Location:          location,
					Dose1Availability: dose1.Availability,
					Dose2Availability: dose2.Availability,
				}
				locationMap[location.ExtId] = locationAvailability

				hasDone1 := hasDoses(dose1.Availability)
				hasDone2 := hasDoses(dose2.Availability)
				if hasDone1 && hasDone2 {
					byDose.Both = append(byDose.Both, locationAvailability)
				} else if hasDone1 && !hasDone2 {
					byDose.Dose1Only = append(byDose.Dose1Only, locationAvailability)
				} else if !hasDone1 && hasDone2{
					byDose.Dose2Only = append(byDose.Dose2Only, locationAvailability)
				} else {
					byDose.Neither = append(byDose.Neither, locationAvailability)
				}
			}
		}
	}

	return byDose
}

func hasDoses(availability[] model.Availability) bool {
	for _, date := range availability {
		if date.Available {
			return true
		}
	}
	return false
}
