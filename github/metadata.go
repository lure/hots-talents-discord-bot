package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-discord-bot/stringutils"
	"io"
	"log"
	"net/http"
	"sync"

	"strings"
)

// [hero] -> [
//
//			   [1] -> [1,2,3,4],
//			   [2]->[1,2,3,4]
//				...
//	     ]
type TalentsType map[string]HeroMeta
type HeroMeta struct {
	Portrait string
	Talents  [][]string
}

// / ----------------- Talent Localized Consts ---------------
type talentNamesJson struct {
	Gamestrings struct {
		Abiltalent struct {
			Names map[string]string `json:"name"`
		} `json:"abiltalent"`
	} `json:"gamestrings"`
}

// / ----------------- Talent system ---------------
type talentTreesJson map[string]heroStructJson
type heroStructJson struct {
	HyperlinkId string                        `json:"hyperlinkId"`
	Talents     map[string][]talentStructJson `json:"talents"`
	Portraits   struct {
		Target string `json:"target"`
	} `json:"portraits"`
}
type talentStructJson struct {
	Name   string `json:"nameId"`
	Button string `json:"buttonId"`
	Sort   int    `json:"sort"`
}

func ReadTalentSystemFromGithub(talentsUrl, constanstUrl string, useCache bool) (TalentsType, error) {
	var (
		writeCache     = true
		constantsCache cache
	)
	var localTalents TalentsType
	if useCache {
		talents, err := constantsCache.loadFromCache()
		if err != nil {
			log.Println("Failed to read cache", err)
			talents, err = fetchAndParseTalentConstants(talentsUrl, constanstUrl)
			if err != nil {
				return nil, err
			}
		} else {
			writeCache = false
		}
		localTalents = talents
	} else {
		if t, err := fetchAndParseTalentConstants(talentsUrl, constanstUrl); err != nil {
			return nil, err
		} else {
			localTalents = t
		}
	}

	if useCache && writeCache {
		if err := constantsCache.saveParsed(localTalents); err != nil {
			return nil, err
		}
	}
	return localTalents, nil
}

// TOOD goroutines - channels
func fetchAndParseTalentConstants(talentsUrl, constanstUrl string) (TalentsType, error) {
	log.Println("Reading metadata from github")
	var (
		wg          sync.WaitGroup
		talentTrees talentTreesJson
		constants   map[string]string
	)
	wg.Add(2)
	var errSys, errConst error
	go func() {
		defer wg.Done()
		talentTrees, errSys = fetchTalentSystem(talentsUrl)
	}()
	go func() {
		defer wg.Done()
		constants, errConst = fetchTalentConsts(constanstUrl)
	}()
	wg.Wait()
	if err := errors.Join(errSys, errConst); err != nil {
		return nil, err
	}
	return traverse(talentTrees, constants)
}

// returns a map of [talentId|buttonId]->localizedTalentNames
func fetchTalentConsts(constantsUrl string) (map[string]string, error) {
	rawString, err := fetchGithubResource(constantsUrl)
	if err != nil {
		return nil, err
	}
	var talentNames talentNamesJson
	if err := json.Unmarshal(rawString, &talentNames); err != nil {
		return nil, fmt.Errorf("can't unmarshal JSON '%s' : %w", string(rawString), err)
	}

	result := make(map[string]string)

	for key, value := range talentNames.Gamestrings.Abiltalent.Names {
		index, entry := 0, 0
		for entry != 2 && index < len(key) {
			if key[index] == '|' {
				entry++
			}
			index++
		}
		result[key[:index-1]] = value
	}

	return result, nil
}

func fetchTalentSystem(talentsUrl string) (talentTreesJson, error) {
	rawTalents, err := fetchGithubResource(talentsUrl)
	if err != nil {
		return nil, err
	}
	var talentTrees talentTreesJson
	if err := json.Unmarshal(rawTalents, &talentTrees); err != nil {
		return nil, fmt.Errorf("can't unmarshal JSON '%s' : %w", string(rawTalents), err)
	} else {
		return talentTrees, nil
	}
}

func fetchGithubResource(url string) ([]byte, error) {
	resp, err := githubClient.Get(url)
	switch {
	case err != nil:
		return nil, fmt.Errorf("can't load url: '%s': %w", url, err)
	case resp.StatusCode != http.StatusOK:
		defer resp.Body.Close()
		return nil, fmt.Errorf("can't read response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	if content, err := io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	} else {
		return content, nil
	}
}

// name -> "talents" -> level1 level4 level7 level10, level13, level16, level20
func traverse(data talentTreesJson, localizedTalentNames map[string]string) (TalentsType, error) {
	result := make(TalentsType)
	for _, value := range data {
		if _, err := readTalents(value, localizedTalentNames, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func readTalents(heroData heroStructJson, localizedTalentNames map[string]string, result TalentsType) (TalentsType, error) {
	hero := stringutils.Normalize(heroData.HyperlinkId)
	result[hero] = HeroMeta{
		Portrait: heroData.Portraits.Target,
		Talents:  make([][]string, len(heroData.Talents)),
	}

	for talentLevel, talents := range heroData.Talents {
		level := extractLevel(talentLevel)
		sortedTalents := make([]string, len(talents))
		for _, talentDescription := range talents {
			nameKey := talentDescription.Name + "|" + talentDescription.Button
			if readable, ok := localizedTalentNames[nameKey]; !ok {
				return nil, fmt.Errorf("no constant found: %s", nameKey)
			} else {
				sortedTalents[talentDescription.Sort-1] = readable
			}
		}
		result[hero].Talents[level] = sortedTalents
	}
	return result, nil
}

var levelMatch = map[string]int{"1": 0, "4": 1, "7": 2, "10": 3, "13": 4, "16": 5, "20": 6}

func extractLevel(s string) int {
	mbLevel, _ := strings.CutPrefix(s, "level")
	return levelMatch[mbLevel]
}
