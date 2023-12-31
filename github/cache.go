package github

import (
	"encoding/json"
	"log"
	"os"
)

// CACHE OPERATIONS

type cache struct{}

func (*cache) makeCachePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := dir + string(os.PathSeparator) + "cache" + string(os.PathSeparator)

	if err = os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	return (path + "talents_dictionary.json"), nil
}

func (c *cache) saveParsed(talents TalentsType) {
	path, err := c.makeCachePath()
	if err != nil {
		log.Println("Failed to get current dir", err)
		return
	}
	bytes, err := json.Marshal(talents)
	if err != nil {
		log.Println("Failed to marshall talents", err)
		return
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		log.Println("Failed to save talent cache", err)
	}
}

func (c *cache) loadFromCache() (TalentsType, error) {
	log.Println("Reading constants from local cache")
	path, err := c.makeCachePath()
	if err != nil {
		log.Println("Failed to get current dir", err)
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result TalentsType
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
