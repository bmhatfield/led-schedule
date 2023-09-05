package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/bmhatfield/led-schedule/jsonapi"
)

const TimeOfDay = "3:04:05 PM"

type Response struct {
	Result Sunset `json:"results"`
	Status string `json:"status"`
}

type Sunset struct {
	Sunrise    string `json:"sunrise"`
	Sunset     string `json:"sunset"`
	FirstLight string `json:"first_light"`
	LastLight  string `json:"last_light"`
	Dawn       string `json:"dawn"`
	Dusk       string `json:"dusk"`
	SolarNoon  string `json:"solar_noon"`
	GoldenHour string `json:"golden_hour"`
	DayLength  string `json:"day_length"`
	Timezone   string `json:"timezone"`
	UtcOffset  int    `json:"utc_offset"`
}

func (s Sunset) parse(t string) (time.Time, error) {
	now := time.Now()

	loc, err := time.LoadLocation(s.Timezone)
	if err != nil {
		return now, err
	}

	x, err := time.ParseInLocation(TimeOfDay, t, loc)
	if err != nil {
		return now, err
	}

	return x.AddDate(now.Year(), int(now.Month())-1, now.Day()-1), nil
}

func (s Sunset) SunsetTime() (time.Time, error) {
	return s.parse(s.Sunset)
}

func (s Sunset) SunriseTime() (time.Time, error) {
	return s.parse(s.Sunrise)
}

func sunset(lat, lon float64) (*Sunset, error) {
	if lat == 0 || lon == 0 {
		return nil, fmt.Errorf("invalid lat/lon")
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.sunrisesunset.io",
		Path:   "/json",
		RawQuery: url.Values{
			"lat": []string{fmt.Sprintf("%f", lat)},
			"lng": []string{fmt.Sprintf("%f", lon)},
		}.Encode(),
	}

	resp, err := jsonapi.Get[Response](u)
	if err != nil {
		return nil, err
	}

	return &resp.Result, nil
}
