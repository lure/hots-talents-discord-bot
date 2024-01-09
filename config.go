package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

// Simple config reader
// Loads the configuration option set passed as argument from envinroment and
// configuration file. The env has precedence over the file

func readConfig(credentialsPath string, options map[string]string) (result map[string]string, err error) {
	log.Println("Looking up credentials")

	fromEnvFunc := func() map[string]string {
		result = make(map[string]string)
		for name := range options {
			if apikey, ok := os.LookupEnv(name); ok {
				log.Println("\tFound " + name + " in env")
				result[name] = apikey
			}
		}
		return result
	}

	fromFileFunc := func(filename string) (map[string]string, error) {
		log.Println("Reading credentials from " + credentialsPath)
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lines := make(map[string]string)
		for scanner.Scan() {
			pair := strings.Split(scanner.Text(), "=")
			if len(pair) != 2 {
				return nil, errors.New("credentials file has incorrect format")
			}
			lines[strings.Trim(pair[0], " ")] = strings.Trim(pair[1], " ")
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		} else {
			return lines, nil
		}
	}

	credsFromEnv := fromEnvFunc()
	credsFromFile, err := fromFileFunc(credentialsPath)
	if err != nil {
		log.Println(err)
	}

	var missedKeys []string
	for option := range options {
		if c, ok := credsFromEnv[option]; ok {
			options[option] = c
		} else {
			if c, ok := credsFromFile[option]; ok {
				options[option] = c
			}
		}

		if options[option] == "" {
			missedKeys = append(missedKeys, option)
		}

	}
	if len(missedKeys) > 0 {
		return nil, errors.New("Some of Api keys are missing! The required keys are [" + strings.Join(missedKeys, ",") + "]")
	}

	return options, nil
}
