package google

import (
	"context"
	"fmt"
	"go-discord-bot/stringutils"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// AIzaSyClDshTtfadZdct-PxZ8Oqx4w2fkbBzqwc
//https://console.cloud.google.com/apis/credentials?project=localtest-1025
//https://stackoverflow.com/questions/39691100/golang-google-sheets-api-v4-write-update-example

func FetchFanGoogleSheet(apiKey string, spreadSheetID string, readRange []string) (map[string]map[string]string, error) {
	builds := make(map[string]map[string]string)
	srv, err := sheets.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	for _, rng := range readRange {
		resp, err := srv.Spreadsheets.Values.Get(spreadSheetID, rng).Do()

		if err != nil {
			return nil, err
		}

		for _, row := range resp.Values {
			v, ok := row[1].(string)
			if !ok {
				// return nil, fmt.Errorf("can't process type %T", row[1])
				log.Printf("Skipping row '%v'", v[1])
				continue
			}
			idx := strings.IndexRune(v, ',')
			if idx < 0 {
				return nil, fmt.Errorf("can't split string '%s' by rune `,`", v)
			}
			rawName := v[idx+1:]
			hero := stringutils.Normalize(rawName)
			if _, ok := builds[hero]; !ok {
				builds[hero] = make(map[string]string)
			}
			builds[hero][row[0].(string)] = v
		}
	}

	return builds, nil
}
