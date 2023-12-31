package google

import (
	"context"
	"go-discord-bot/stringutils"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// AIzaSyClDshTtfadZdct-PxZ8Oqx4w2fkbBzqwc
//https://console.cloud.google.com/apis/credentials?project=localtest-1025
//https://stackoverflow.com/questions/39691100/golang-google-sheets-api-v4-write-update-example

var FanBuilds map[string]map[string]string

func FetchFanGoogleSheet(apiKey string, spreadSheetID string, readRange []string) map[string]map[string]string {
	builds := make(map[string]map[string]string)
	srv, err := sheets.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	for _, rng := range readRange {
		resp, err := srv.Spreadsheets.Values.Get(spreadSheetID, rng).Do()

		if err != nil {
			panic(err)
		}

		for _, row := range resp.Values {
			rawName := strings.Split(row[1].(string), ",")[1]
			hero := stringutils.Normalize(rawName)
			if _, ok := builds[hero]; !ok {
				builds[hero] = make(map[string]string)
			}
			builds[hero][row[0].(string)] = row[1].(string)
		}
	}

	FanBuilds = builds
	return builds
}
