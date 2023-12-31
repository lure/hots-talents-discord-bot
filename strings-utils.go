package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type HName struct {
	normalized  string
	capitalized string
	icon        string
}

const portraitUrl = "https://raw.githubusercontent.com/HeroesToolChest/heroes-images/master/heroesimages/heroportraits/ui_targetportrait_hero_%s.png"
const buildLink = "https://psionic-storm.com/en/talent-calculator/%s/#talents=%s"

// transliterate normalized input user name to hots avatar file names
var nameToHotsName = map[string]string{
	"brightwing": "faeriedragon",
	"blaze":      "firebat",
	"cassia":     "d2amazonf",
	"etc":        "l90etc",
	"greymane":   "genngreymane",
	"sgthammer":  "sgthammer",
	"kharazim":   "monk",
	"ltmorales":  "medic",
	"mei":        "meiow",
	"liming":     "wizard",
	"nazeebo":    "witchdoctor",
	"qhira":      "nexus2",
	"sonya":      "barbarian",
	"thebutcher": "butcher",
	"valla":      "demonhunter",
	"xul":        "necromancer",
}

// transliterate normalized input user name to psion website talent page
var nameToPsion = map[string]string{
	"thebutcher":  "the-butcher",
	"lostvikings": "the-lost-vikings",
	"ltmorales":   "lt-morales",
	"liming":      "li-ming",
	"sgthammer":   "sgt-hammer",
}

var inputToName = map[string]string{
	"butcher":        "thebutcher",
	"thelostvikings": "lostvikings",
	"vikings":        "lostvikings",
	"morales":        "ltmorales",
	"hammer":         "sgthammer",
}

var nameNormalizeRegex *regexp.Regexp
var psionBuildRegex *regexp.Regexp
var caser = cases.Title(language.AmericanEnglish)

func initStringUtils() {
	var err error
	nameNormalizeRegex, err = regexp.Compile(`[\s\'\-\[\]]`)
	if err != nil {
		panic(err)
	}
	psionBuildRegex, err = regexp.Compile(`\d+`)
	if err != nil {
		panic(err)
	}
}

func normalize(s string) string {
	return strings.ToLower(nameNormalizeRegex.ReplaceAllString(s, ""))
}

func prepareName(s string) HName {
	str, _ := strings.CutPrefix(s, "!")
	caps := caser.String(strings.ToLower(str))
	nameCandidat := normalize(str)
	norm, ok := inputToName[nameCandidat]
	if !ok {
		norm = nameCandidat
	}

	thumb, ok := nameToHotsName[norm]
	if !ok {
		thumb = norm
	}
	thumb = fmt.Sprintf(portraitUrl, thumb)
	return HName{
		normalized:  norm,
		capitalized: caps,
		icon:        thumb,
	}
}

func toArrayOfNumbers(talents string) []int {
	numbers := psionBuildRegex.FindAllString(talents, -1)
	result := make([]int, 7)
	for i, n := range numbers[0] {
		v, err := strconv.Atoi(string(n))
		if err != nil {
			panic(err)
		}
		result[i] = v
	}
	return result
}

func makePsionicTalents(hero string, talents string) string {
	var result string
	numbers := psionBuildRegex.FindAllString(talents, -1)
	if len(numbers) == 0 {
		fmt.Println("Incorrect data for hero:" + hero + " and talents: " + talents)
		return ""
	}
	for _, n := range numbers[0] {
		if len(result) > 0 {
			result += "-"
		}
		result += string(n)
	}

	psionHero, ok := nameToPsion[hero]
	if !ok {
		psionHero = hero
	}

	return fmt.Sprintf(buildLink, psionHero, result)
}
