package main

const botApiKey = "BOT_APIKEY"
const googleApiKey = "BOT_GOOGLEKEY"
const portraitUrl = "https://raw.githubusercontent.com/HeroesToolChest/heroes-images/master/heroesimages/heroportraits/%s"
const talentsUrl = "https://raw.githubusercontent.com/HeroesToolChest/heroes-data/master/heroesdata/2.55.4.91368/data/herodata_91368_localized.json"
const constanstUrl = "https://github.com/HeroesToolChest/heroes-data/raw/master/heroesdata/2.55.4.91368/gamestrings/gamestrings_91368_enus.json"
const spreadSheetID = "1kiHfe0obByIt5qBvLNeVwrKr2SLl_laTCDqbfJwI9X8"

// This ranges are looked up from the document referenced "spreadSheetID"
var readRange = []string{
	"B3:C40",
	"F3:G40",
	"J3:K40",
	"N3:O27",
}

const reaction = "📨"
