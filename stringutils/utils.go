package stringutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"log"
)

type HName struct {
	Normalized  string
	Capitalized string
	HotsName    string
}

const psionicStormCalculatorLink = "https://psionic-storm.com/en/talent-calculator/%s/#talents=%s"
const icyVeinsCalculatorLink = "https://www.icy-veins.com/heroes/talent-calculator/%s#55.1!%s"

// transliterate normalized input user name to psion website talent page
var nameToPsion = map[string]string{
	"liming":      "li-ming",
	"lostvikings": "the-lost-vikings",
	"ltmorales":   "lt-morales",
	"sgthammer":   "sgt-hammer",
	"thebutcher":  "the-butcher",
}

// transliterate normalized input user name to icy-veins website talent page
var nametoIcyVeins = map[string]string{
	"deckard":     "deckard-cain",
	"etc":         "e-t-c",
	"ketlthuzad":  "kel-thuzad",
	"lili":        "li-li",
	"liming":      "li-ming",
	"ltmorales":   "lt-morales",
	"lostvikings": "the-lost-vikings",
	"thebutcher":  "the-butcher",
	"sgthammer":   "sgt-hammer",
}

// handle some common missprints
var inputToName = map[string]string{
	"butcher":        "thebutcher",
	"thelostvikings": "lostvikings",
	"vikings":        "lostvikings",
	"morales":        "ltmorales",
	"hammer":         "sgthammer",
}

var (
	nameNormalizeRegex = must(regexp.Compile(`[\s\'\-\[\]]`))

	numbersOnlyRegex = must(regexp.Compile(`\d+`))
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func Normalize(s string) string {
	return strings.ToLower(nameNormalizeRegex.ReplaceAllString(s, ""))
}

// Converts user input into lovercase char-only form, then uses it to find the internal representation if any,
// The internal representation is used to find service-specific variants in a later calls
func PrepareName(s string) string {

	str, _ := strings.CutPrefix(s, "!")
	normalizedName := Normalize(str)
	return GetOrDefault(normalizedName, inputToName)
}

// Converts the name offered into specific service name, or returns the
// original name if no susbstitute found
func GetOrDefault(key string, conversion map[string]string) string {
	converted, ok := conversion[key]
	if !ok {
		converted = key
	}
	return converted
}

// extracts numbers only from the string of form "[1232131, Valla]"
func getTalentsFromBuild(build string) (string, error) {
	numbers := numbersOnlyRegex.FindAllString(build, -1)
	if len(numbers) != 1 {
		return "", fmt.Errorf("incorrect data talents: %s", build)
	}
	return numbers[0], nil
}

// converts string "[1232131, Valla]"" to []int{1,2,3,2,1,3,1}
func BuildToSevenNumbers(talents string) ([7]int, error) {
	numbers := numbersOnlyRegex.FindAllString(talents, -1)
	var result [7]int
	// TODO shall we iterate over runes in the first digits string found?
	for i, n := range numbers[0] {
		v, err := strconv.Atoi(string(n))
		if err != nil {
			return result, fmt.Errorf("can't parse '%v' : %w", n, err)
		}
		result[i] = v
	}
	return result, nil
}

func GetExternalLinks(hero string, talents string) map[string]string {

	makePsionicTalents := func(hero string, numbers string) string {
		var result strings.Builder
		for _, n := range numbers {
			if result.Len() > 0 {
				result.WriteString("-")
			}
			result.WriteString(string(n))
		}

		substitutedHero := GetOrDefault(hero, nameToPsion)
		return fmt.Sprintf(psionicStormCalculatorLink, substitutedHero, result.String())
	}

	makeIcyTalents := func(hero string, numbers string) string {
		substitutedHero := GetOrDefault(hero, nametoIcyVeins)
		return fmt.Sprintf(icyVeinsCalculatorLink, substitutedHero, numbers)
	}

	numbers, err := getTalentsFromBuild(talents)
	if err != nil {
		log.Println(err)
		return map[string]string{}
	}

	result := make(map[string]string)
	result["psionic-storm"] = makePsionicTalents(hero, numbers)
	result["icy-veins"] = makeIcyTalents(hero, numbers)
	return result
}
