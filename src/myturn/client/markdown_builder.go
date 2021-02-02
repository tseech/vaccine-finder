package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func BuildMarkdownFile(Locations LocationAvailabilityByDose, NumberOfDaysSearched time.Duration) {
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
	writer.WriteString(fmt.Sprintf("*Date range: %s - %s*\n\n",
		time.Now().Format("Mon Jan 2 2006"),
		time.Now().Add(time.Hour*24*NumberOfDaysSearched).Format("Mon Jan 2 2006")))
	writer.WriteString(fmt.Sprintf("*Go to: <https://myturn.ca.gov> to schedule your appointment*\n\n"))
	writer.WriteString("\n")

	writer.WriteString("## Locations with both doses\n\n")
	writeLocations(Locations.Both, writer)

	writer.WriteString("## Locations with dose 1 only\n\n")
	writeLocations(Locations.Dose1Only, writer)

	writer.WriteString("## Locations with dose 2 only\n\n")
	writeLocations(Locations.Dose2Only, writer)

	writer.WriteString("## Locations with neither\n\n")
	writeLocations(Locations.Neither, writer)

	writer.Flush()
}

func writeLocations(availability[] LocationAvailability, writer *bufio.Writer) {
	for _, location := range availability {
		writer.WriteString(fmt.Sprintf(">### %s\n",location.Location.Name))
		writer.WriteString(fmt.Sprintf(">#### %s\n", location.Location.DisplayAddress))

		hasDose1 := 0
		dose1Dates := ""
		hasDose2 := 0
		dose2Dates := ""
		for i:=0; i<len(location.Dose1Availability) && i<len(location.Dose2Availability); i++ {
			if location.Dose1Availability[i].Available {
				hasDose1++
				if dose1Dates != "" {
					dose1Dates += ", "
				}
				dose1Dates += location.Dose1Availability[i].Date
			}
			if location.Dose2Availability[i].Available {
				hasDose2++
				if dose2Dates != "" {
					dose2Dates += ", "
				}
				dose2Dates += location.Dose2Availability[i].Date
			}
		}

		writer.WriteString(fmt.Sprintf(">- Dose 1 available on %d days\n", hasDose1))
		if dose1Dates != "" {
			writer.WriteString(fmt.Sprintf(">  - Days: %s\n", dose1Dates))
		}
		writer.WriteString(fmt.Sprintf(">- Dose 2 available on %d days\n", hasDose2))
		if dose2Dates != "" {
			writer.WriteString(fmt.Sprintf(">  - Days: %s\n", dose2Dates))
		}

		writer.WriteString("\n")
	}

	if len(availability) == 0 {
		writer.WriteString(">None\n\n")
	}
}
