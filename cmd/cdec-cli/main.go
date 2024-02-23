package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const cdecBaseUrl = "https://cdec.water.ca.gov/dynamicapp/req/JSONDataServlet"

type queryOptions struct {
	station   string
	sensor    string
	dur_code  string
	startDate string
	endDate   string
}

func main() {
	fmt.Println("Welcome to the CDEC CLI!")

	if len(os.Args) < 2 {
		fmt.Println("Usage: cdec-cli <command> [arguments]")
		os.Exit(1)
	}

	var options queryOptions

	switch os.Args[1] {
	case "query":
		{
			queryCommand := flag.NewFlagSet("query", flag.ExitOnError)
			queryCommand.StringVar(&options.station, "station", "", "Station ID to query data for. Use 'cdec-cli station --list' for a complete list of available stations")
			queryCommand.StringVar(&options.sensor, "sensor", "", "Sensor number to query data for")
			queryCommand.StringVar(&options.dur_code, "duration", "", "duration code for data, can be one of (d)aily, (h)ourly, or (m)onthly")
			queryCommand.StringVar(&options.startDate, "startdate", "", "Query start date (YYYY-mm-dd)")
			queryCommand.StringVar(&options.endDate, "enddate", "", "Query end date (YYYY-mm-dd)")

			queryCommand.Parse(os.Args[2:])

			params := url.Values{}
			params.Set("Stations", options.station)
			params.Set("SensorNums", options.sensor)
			params.Set("dur_code", options.dur_code)
			params.Set("Start", options.startDate)
			params.Set("End", options.endDate)

			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

			u, err := url.Parse(cdecBaseUrl)
			if err != nil {
				logger.Println("error trying to encode the url")
				return
			}

			u.RawQuery = params.Encode()
			logger.Printf("Using the following url for retrieving data: \n%s\n\n", u.String())

			resp, err := http.Get(u.String())
			if err != nil {
				logger.Println("error trying to reach cdec")
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			fmt.Println(string(body))
		}
	case "stations":
		{
			fmt.Println("yet to be implemented!")
		}
	}
	flag.Parse()

}
