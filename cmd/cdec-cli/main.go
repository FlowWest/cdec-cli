package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

var cdecUrls = map[string]string{
	"query":    "https://cdec.water.ca.gov/dynamicapp/req/JSONDataServlet",
	"stations": "https://cdec.water.ca.gov/dynamicapp/staMeta",
}

type cliCmd struct {
	cmd  string
	args []string
}

type queryOptions struct {
	station   string
	sensor    string
	dur_code  string
	startDate string
	endDate   string
}

// For sure overkill but just to stay consistent
type stationOptions struct {
	stationId string
}

func findTables(n *html.Node) []*html.Node {
	var tables []*html.Node

	var findTable func(n *html.Node)
	findTable = func(*html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tables = append(tables, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTable(c)
		}

	}
	return tables
}

func parseHTMLTable(n *html.Node) (string, string) {

}

func main() {
	fmt.Println("Welcome to the CDEC CLI!")

	if len(os.Args) < 2 {
		fmt.Println("Usage: cdec-cli <command> [arguments]")
		os.Exit(1)
	}

	if os.Args[1] == "--help" {
		fmt.Println("Usage: cdec-cli <command> [arguments]")
		os.Exit(1)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	switch os.Args[1] {
	case "query":
		{
			var options queryOptions
			// set the query command
			queryCommand := flag.NewFlagSet("query", flag.ExitOnError)

			// append to this base command arguments needed to query the data
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

			u, err := url.Parse(cdecUrls["query"])
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
				logger.Println("Error reading response body:", err)
				return
			}

			fmt.Println(string(body))
		}
	case "stations":
		{
			var stationOptions stationOptions
			stationCommand := flag.NewFlagSet("station", flag.ExitOnError)

			stationCommand.StringVar(&stationOptions.stationId, "stationID", "", "Station ID to retrieva metadata for")

			stationCommand.Parse(os.Args[2:])

			params := url.Values{}
			params.Set("station_id", stationOptions.stationId)

			// build the url for stations
			u, err := url.Parse(cdecUrls["stations"])
			if err != nil {
				logger.Println("error trying to parse url")
				return
			}

			u.RawQuery = params.Encode()

			resp, err := http.Get(u.String())
			if err != nil {
				logger.Println("error trying to reach the statiosn endpoint")
				return
			}

			doc, err := html.Parse(resp.Body)
			if err != nil {
				logger.Println("error trying to parse the body")
				return
			}

			findTable(doc)
			logger.Printf("you are trying to get info for %s", stationOptions.stationId)
		}
	}
	flag.Parse()

}
