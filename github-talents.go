package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

const talentsUrl = "https://raw.githubusercontent.com/HeroesToolChest/heroes-data/master/heroesdata/2.55.4.91368/data/herodata_91368_localized.json"
const constanstUrl = "https://github.com/HeroesToolChest/heroes-data/raw/master/heroesdata/2.55.4.91368/gamestrings/gamestrings_91368_enus.json"

// hero ->
//
//		[1] -> [1,2,3,4],
//	 [2]->[1,2,3,4]
type TalentsType map[string][][]string

var TalentsDictionary TalentsType

type RawJson map[string]interface{}

type GithubTalent struct {
}

func readTalentsFromGithub() TalentsType {
	rawJson := readGithubResource(talentsUrl)
	constants := readConsts()
	talents := make(TalentsType)
	TalentsDictionary = traverse(rawJson, talents, constants)
	return talents
}

func readConsts() map[string]string {
	rawJson := readGithubResource(constanstUrl)
	result := make(map[string]string)

	names := ((rawJson["gamestrings"].(map[string]interface{}))["abiltalent"].(map[string]interface{}))["name"].(map[string]interface{})

	for key, value := range names {
		result[key[:strings.Index(key, "|")]] = value.(string)
	}

	return result
}

func readGithubResource(url string) RawJson {
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
	var rawJson RawJson
	err = json.Unmarshal(content, &rawJson)
	if err != nil {
		panic(err)
	}
	return rawJson
}

// name -> "talents" -> level1 level4 level7 level10, level13, level16, level20
func traverse(data interface{}, talents TalentsType, readableNames map[string]string) TalentsType {
	switch v := data.(type) {
	case RawJson:
		for key, value := range v {
			hero := normalize(key)
			talents[hero] = make([][]string, 7)
			readTalents(value.(map[string]interface{}), talents, hero, readableNames)
		}
	}
	return talents
}

type RawTalent struct {
	order int
	id    string
}
type RawTalents []RawTalent

func (rt RawTalents) Len() int               { return len(rt) }
func (rt RawTalents) Swap(i int, j int)      { rt[i], rt[j] = rt[j], rt[i] }
func (rt RawTalents) Less(i int, j int) bool { return rt[i].order < rt[j].order }

func readTalents(heroData RawJson, result TalentsType, hero string, readableNames map[string]string) TalentsType {
	talentsLevelsObject := heroData["talents"]
	for talentLevel, talents := range talentsLevelsObject.(map[string]interface{}) {
		level := extractLevel(talentLevel)
		levelTalents := RawTalents{}

		for _, talenDescription := range talents.([]interface{}) {
			switch v := talenDescription.(type) {
			case map[string]interface{}:
				order := int(v["sort"].(float64))
				name := v["nameId"].(string)
				levelTalents = append(levelTalents, RawTalent{order, name})
			default:
				panic(talenDescription)
			}
		}
		sort.Sort(&levelTalents)
		sortedTalents := make([]string, len(levelTalents))
		for i, t := range levelTalents {
			readable, ok := readableNames[t.id]
			if !ok {
				panic(fmt.Sprintf("No constant for %s found\n", t.id))
			}
			sortedTalents[i] = readable
		}
		result[hero][level] = sortedTalents
	}
	return result
}

var levelMatch = map[string]int{
	"1":  0,
	"4":  1,
	"7":  2,
	"10": 3,
	"13": 4,
	"16": 5,
	"20": 6,
}

func extractLevel(s string) int {
	mbLevel, _ := strings.CutPrefix(s, "level")
	return levelMatch[mbLevel]
}
