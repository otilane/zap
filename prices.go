package main

import (
	"encoding/json"
	"fmt"
	"time"
	"web"
)

// Prices struct for the API
type Prices struct {
	DateTime           string  `json:"fromDateTime"`
	CentsPerKwhWithVat float64 `json:"centsPerKwhWithVat"`
}

type Resolution string

// Enums for resolution
const (
	Hour  Resolution = "one_hour"
	Day   Resolution = "one_day"
	Week  Resolution = "one_week"
	Month Resolution = "one_month"
	Year  Resolution = "one_year"
)

func fetchPrices(startDateTime time.Time, endDateTime time.Time, resolution Resolution) (*map[int64]float64, error) {
	data := make(map[int64]float64)

	// Set the params to be url encoded
	params := map[string]string{
		"startDateTime": startDateTime.Format(time.RFC3339),
		"endDateTime":   endDateTime.Format(time.RFC3339),
		"resolution":    fmt.Sprint(resolution),
	}

	apiURL := web.SetParams("https://estfeed.elering.ee/api/public/v1/energy-price/electricity?", params)

	// Fetch the prices data
	body, err := web.FetchBody(web.Get, apiURL, nil, nil)
	if err != nil {
		fmt.Println(err)
		return &data, err
	}

	// Unmarshal (unencode) the JSON into the Prices struct
	var prices []Prices
	if err = json.Unmarshal(body, &prices); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return &data, err
	}

	// Squash starting DateTime into a UNIX timestamp and build the map to return
	for v := range prices {
		toTime, err := time.Parse(time.RFC3339, prices[v].DateTime)
		if err == nil {
			data[toTime.Unix()] = prices[v].CentsPerKwhWithVat
		}
	}

	return &data, err
}
