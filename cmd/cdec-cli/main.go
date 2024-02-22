package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const cdecBaseUrl = "https://cdec.water.ca.gov/dynamicapp/wsSensorData"

type queryOptions struct {
	station   string
	sensor    string
	startDate string
	endDate   string
}

func main() {
	fmt.Println("Welcome to the CDEC CLI!")

	var options queryOptions

	flag.StringVar(&options.station, "station", "", "Station ID to query data for. Use 'cdec-cli station --list' for a complete list of available stations")
	flag.StringVar(&options.sensor, "sensor", "", "Sensor number to query data for")
	flag.StringVar(&options.startDate, "start-date", "", "Query start date (YYYY-mm-dd)")
	flag.StringVar(&options.endDate, "end-date", "", "Query end date (YYYY-mm-dd)")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	resp, err := http.Get(cdecBaseUrl)
	if err != nil {
		logger.Println("error trying to reach cdec")
		return
	}

	logger.Printf("the resp returned with status code: %s", resp.Body)
}
