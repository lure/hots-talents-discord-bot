package main

import (
	"context"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// AIzaSyClDshTtfadZdct-PxZ8Oqx4w2fkbBzqwc
//https://console.cloud.google.com/apis/credentials?project=localtest-1025
//https://stackoverflow.com/questions/39691100/golang-google-sheets-api-v4-write-update-example

var FanBuilds map[string]map[string]string

func parseFanGoogleSheet() map[string]map[string]string {
	builds := make(map[string]map[string]string)
	apiKey := "AIzaSyClDshTtfadZdct-PxZ8Oqx4w2fkbBzqwc"
	srv, err := sheets.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	spreadSheetID := "1kiHfe0obByIt5qBvLNeVwrKr2SLl_laTCDqbfJwI9X8"
	readRange := []string{
		"B3:C40",
		"F3:G40",
		"J3:K40",
		"N3:O27",
	}

	for _, rng := range readRange {
		resp, err := srv.Spreadsheets.Values.Get(spreadSheetID, rng).Do()

		if err != nil {
			panic(err)
		}

		for _, row := range resp.Values {
			rawName := strings.Split(row[1].(string), ",")[1]
			hero := normalize(rawName)
			if _, ok := builds[hero]; !ok {
				builds[hero] = make(map[string]string)
			}
			builds[hero][row[0].(string)] = row[1].(string)
		}
	}

	FanBuilds = builds
	return builds

	// for key, value := range builds {
	// 	fmt.Println(key)
	// 	fmt.Println(value)
	// 	break
	// }
}
