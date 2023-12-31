package github

import (
	"encoding/json"
	"fmt"
	"go-discord-bot/stringutils"
	"io"
	"log"
	"net/http"

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

var TalentsDictionary TalentsType

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

func ReadTalentSystemFromGithub(talentsUrl, constanstUrl string, useCache bool) TalentsType {
	var talents TalentsType
	var err error
	writeCache := true
	var constantsCache cache
	if useCache {
		talents, err = constantsCache.loadFromCache()
		if err != nil {
			log.Println("Failed to read cache", err)
			talents = fetchAndParseTalentConstants(talentsUrl, constanstUrl)
		} else {
			writeCache = false
		}
	} else {
		talents = fetchAndParseTalentConstants(talentsUrl, constanstUrl)
	}

	TalentsDictionary = talents

	if useCache && writeCache {
		constantsCache.saveParsed(talents)
	}
	return talents
}

// TOOD goroutines - channels
func fetchAndParseTalentConstants(talentsUrl, constanstUrl string) TalentsType {
	log.Println("Reading metadata from github")
	talentTrees := fetchTalentSystem(talentsUrl)
	constants := fetchTalentConsts(constanstUrl)
	return traverse(talentTrees, constants)
}

// returns a map of [talentId|buttonId]->localizedTalentNames
func fetchTalentConsts(constanstUrl string) map[string]string {
	rawString := fetchGithubResource(constanstUrl)
	var talentNames talentNamesJson
	err := json.Unmarshal(rawString, &talentNames)
	if err != nil {
		panic(err)
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

	return result
}

func fetchTalentSystem(talentsUrl string) talentTreesJson {
	rawTalents := fetchGithubResource(talentsUrl)
	var talentTrees talentTreesJson
	err := json.Unmarshal(rawTalents, &talentTrees)
	if err != nil {
		panic(err)
	}
	return talentTrees
}

func fetchGithubResource(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status code %d is not OK for %s\n", resp.StatusCode, url)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return content
}

// name -> "talents" -> level1 level4 level7 level10, level13, level16, level20
func traverse(data talentTreesJson, localizedTalentNames map[string]string) TalentsType {
	result := make(TalentsType)
	for _, value := range data {
		readTalents(value, localizedTalentNames, result)
	}
	return result
}

func readTalents(heroData heroStructJson, localizedTalentNames map[string]string, result TalentsType) TalentsType {
	hero := stringutils.Normalize(heroData.HyperlinkId)
	result[hero] = HeroMeta{
		Portrait: heroData.Portraits.Target,
		Talents:  make([][]string, len(heroData.Talents)),
	}

	for talentLevel, talents := range heroData.Talents {
		level := extractLevel(talentLevel)
		sortedTalents := make([]string, len(talents))
		for _, talenDescription := range talents {
			nameKey := talenDescription.Name + "|" + talenDescription.Button
			readable, ok := localizedTalentNames[nameKey]
			if !ok {
				panic(fmt.Sprintf("No constant for %s found\n", nameKey))
			}
			sortedTalents[talenDescription.Sort-1] = readable
		}
		result[hero].Talents[level] = sortedTalents
	}
	return result
}

var levelMatch = map[string]int{"1": 0, "4": 1, "7": 2, "10": 3, "13": 4, "16": 5, "20": 6}

func extractLevel(s string) int {
	mbLevel, _ := strings.CutPrefix(s, "level")
	return levelMatch[mbLevel]
}
