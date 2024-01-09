package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CACHE OPERATIONS

type cache struct{}

func (*cache) makeCachePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can't get current directory - user not set? : %w", err)
	}

	path := filepath.Join(dir, "cache")

	if err = os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("can't create path %s : %w", path, err)
	}

	return filepath.Join(path, "talents_dictionary.json"), nil
}

func (c *cache) saveParsed(talents TalentsType) error {
	path, err := c.makeCachePath()
	if err != nil {
		return err
	}
	if bytes, err := json.Marshal(talents); err != nil {
		return fmt.Errorf("can't marshal json: %w", err)
	} else if err := os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("can't write file at '%s' : %w", path, err)
	}
	return nil
}

func (c *cache) loadFromCache() (TalentsType, error) {
	log.Println("Reading constants from local cache")
	path, err := c.makeCachePath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read file at %s : %w", path, err)
	}

	var result TalentsType
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json '%s' : %w", string(file), err)
	}

	return result, nil
}
