package main

import (
	"net/url"

	"github.com/bmhatfield/led-schedule/jsonapi"
)

type Geo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func geo() (*Geo, error) {
	u := url.URL{
		Scheme: "http",
		Host:   "ip-api.com",
		Path:   "/json",
	}

	return jsonapi.Get[Geo](u)
}
