package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Config struct {
	IdSite uint32
	Token string
}

func main() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var config Config
	json.Unmarshal(jsonData, &config)


	fmt.Println(fetchPrices(time.Now().UTC(), time.Now().UTC(), Hour))
	fmt.Println(fetchVRMForecast(config.IdSite, config.Token))
	fmt.Println(fetchBatteryData(config.IdSite, config.Token))
}
