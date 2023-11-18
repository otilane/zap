package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"web"
)

type Record struct {
	VrmConsumptionFc [][]json.Number `json:"vrm_consumption_fc"`
	SolarYieldFc     [][]json.Number `json:"solar_yield_forecast"`
}

type VRMForecast struct {
	Records *Record `json:"records"`

	Totals struct {
		VrmConsumptionFcTotal float64 `json:"vrm_consumption_fc"`
		SolarYieldFcTotal     float64 `json:"solar_yield_forecast"`
	} `json:"totals"`
}

type BatteryData struct {
	Records struct {
		Data map[string]struct {
			Code       string `json:"code"`
			ValueFloat string `json:"valueFloat"`
		} `json:"data"`
	} `json:"records"`
}

// Fetch historic and total VRM forecast data.
func fetchVRMForecast(idSite uint32, token string) (*map[int64]float64, *map[int64]float64, *map[int64][2]float64, error) {
	historicVRM := make(map[int64]float64)
	historicSolar := make(map[int64]float64)
	totals := make(map[int64][2]float64)

	// Set the params to be url encoded
	params := map[string]string{
		"interval": "hours",
		"type":     "forecast",
	}

	apiURL := web.SetParams(fmt.Sprintf("https://vrmapi.victronenergy.com/v2/installations/%d/stats?", idSite), params)

	// Attach the token header before fetching the VRM data
	headers := map[string]string{
		"x-authorization": fmt.Sprintf("Token %s", token),
	}

	body, err := web.FetchBody(web.Get, apiURL, nil, headers)
	if err != nil {
		fmt.Println(err)
		return &historicVRM, &historicSolar, &totals, err
	}

	// Unmarshal (unencode) the JSON into the VRMForecast struct
	var forecast VRMForecast

	if err = json.Unmarshal(body, &forecast); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return &historicVRM, &historicSolar, &totals, err
	}

	// Format historic VRM forecast into a return map
	for v := range forecast.Records.VrmConsumptionFc {
		timeStamp, err := forecast.Records.VrmConsumptionFc[v][0].Int64()
		if err != nil {
			fmt.Println(err)
			return &historicVRM, &historicSolar, &totals, err
		}
		value, err := forecast.Records.VrmConsumptionFc[v][1].Float64()
		if err != nil {
			fmt.Println(err)
			return &historicVRM, &historicSolar, &totals, err
		}

		historicVRM[timeStamp] = value 
	}

	// Format historic Solar forecast into a return map
	for v := range forecast.Records.SolarYieldFc {
		timeStamp, err := forecast.Records.SolarYieldFc[v][0].Int64()
		if err != nil {
			fmt.Println(err)
			return &historicVRM, &historicSolar, &totals, err
		}
		value, err := forecast.Records.SolarYieldFc[v][1].Float64()
		if err != nil {
			fmt.Println(err)
			return &historicVRM, &historicSolar, &totals, err
		}

		historicSolar[timeStamp] = value 
	}

	totals[time.Now().Unix()] = [2]float64{forecast.Totals.VrmConsumptionFcTotal, forecast.Totals.SolarYieldFcTotal}
	
	return &historicVRM, &historicSolar, &totals, err
}

// Fetch current battery percentage
func fetchBatteryData(idSite uint32, token string) (*map[int64]map[string]float64, error) {

	data := make(map[int64]map[string]float64)

	// Set the apiURL
	apiURL := fmt.Sprintf("https://vrmapi.victronenergy.com/v2/installations/%d/widgets/BatterySummary", idSite)

	// Attach the token header before fetching the VRM data
	headers := map[string]string{
		"x-authorization": fmt.Sprintf("Token %s", token),
	}

	body, err := web.FetchBody(web.Get, apiURL, nil, headers)
	if err != nil {
		fmt.Println(err)
		return &data, err
	}

	// Unmarshal (unencode) the JSON into the BatteryPercentage struct
	// Error checking should be done here, but not all fields have "valueFloat", which'll just keep throwing errors we don't care about
	// It is easier to just pretend that everything is fine and parse out the data after
	var batteryData BatteryData
	json.Unmarshal(body, &batteryData)


	// Parse out the battery data into a map to return
	for _, battery := range batteryData.Records.Data {
		if battery.Code == "" || battery.ValueFloat == "" {
			continue
		}

		if data[time.Now().Unix()] == nil {
			data[time.Now().Unix()] = make(map[string]float64)
		}
		
		data[time.Now().Unix()][battery.Code], _ = strconv.ParseFloat(battery.ValueFloat, 64)
	}
	return &data, err
}
