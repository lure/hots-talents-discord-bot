package main

import (
	"fmt"
	"go-discord-bot/github"
	"go-discord-bot/google"
	"go-discord-bot/stringutils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var appSettings = map[string]string{
	botApiKey:    "",
	googleApiKey: "",
}

var FanBuilds map[string]map[string]string

var TalentsDictionary github.TalentsType

// https://discord.com/api/oauth2/authorize?client_id=1189981976841699411&permissions=2112&scope=bot
func main() {

	var err error
	_, err = readConfig("local.creds", appSettings)
	if err != nil {
		log.Fatal(err)
		return
	}
	if fb, err := google.FetchFanGoogleSheet(appSettings[googleApiKey], spreadSheetID, readRange); err != nil {
		log.Fatal("Can't fetch google spreadsheet", err)
	} else {
		FanBuilds = fb
	}
	if talents, err := github.ReadTalentSystemFromGithub(talentsUrl, constanstUrl, true); err != nil {
		log.Fatal("Can't read talents from github", err)
	} else {
		TalentsDictionary = talents
	}

	sess, err := discordgo.New("Bot " + appSettings[botApiKey])
	if err != nil {
		log.Fatal(err)
	}

	installMessageListener(sess)
	installReactionListener(sess)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	log.Println("The bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func installMessageListener(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot || (!strings.HasPrefix(m.Content, "!")) {
			return
		}

		heroName := stringutils.PrepareName(m.Content)
		if len(FanBuilds[heroName]) == 0 {
			return
		}

		prt := fmt.Sprintf(portraitUrl, TalentsDictionary[heroName].Portrait)
		thumbnail := discordgo.MessageEmbedThumbnail{
			URL: prt,
		}

		flds := makeMessageFields(heroName)

		footer := discordgo.MessageEmbedFooter{
			Text: "React to the 📨 emoji to get this message sent to your DMs!",
		}

		embed := discordgo.MessageEmbed{
			Thumbnail:   &thumbnail,
			Type:        discordgo.EmbedTypeRich,
			Description: "[Support me here!](https://www.buymeacoffee.com/alexlt)",
			Fields:      flds,
			Footer:      &footer,
		}
		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			log.Println(err)
		} else {
			err = s.MessageReactionAdd(m.ChannelID, msg.ID, reaction)
			if err != nil {
				log.Println(err)
			}
		}
	})
}

func installReactionListener(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == s.State.User.ID {
			return
		}

		if m.Emoji.ID == "" && m.Emoji.Name == reaction {
			origMessage, err := s.ChannelMessage(m.ChannelID, m.MessageID)
			if err != nil {
				log.Println(err)
				return
			}
			if origMessage.Author.ID != s.State.User.ID {
				return
			}

			channel, err := s.UserChannelCreate(m.UserID)
			if err != nil {
				log.Println("Cant open direct messages to user"+m.UserID, err)
				return
			}
			if channel.ID != origMessage.ChannelID {
				_, err = s.ChannelMessageSendEmbed(channel.ID, origMessage.Embeds[0])
				if err != nil {
					log.Println(err)
				}
			}
		}
	})
}

func makeMessageFields(heroName string) []*discordgo.MessageEmbedField {
	flds := make([]*discordgo.MessageEmbedField, 0)

	for buildName, talents := range FanBuilds[heroName] {
		var buffer strings.Builder
		// ------------------- START talents as a list
		talentsAslist, err := stringutils.BuildToSevenNumbers(talents)
		if err != nil {
			log.Println("can't process seven numbers", err)
			continue
		}
		for talentLevel, talentOrder := range talentsAslist {
			heroData := TalentsDictionary[heroName].Talents[talentLevel][talentOrder-1]
			buffer.WriteString(fmt.Sprintf("**[%d]** %s\n", talentOrder, heroData))
		}
		// ------------------- END talents as a list
		// ------------------- START external links
		buffer.WriteString("\n")
		for a, b := range stringutils.GetExternalLinks(heroName, talents) {
			buffer.WriteString(fmt.Sprintf("[%s🔗](%s)\n", a, b))
		}
		// ------------------- END external links
		buffer.WriteString("\n")
		buffer.WriteString(talents)

		fld := discordgo.MessageEmbedField{
			Name:   buildName,
			Value:  buffer.String(),
			Inline: true,
		}
		flds = append(flds, &fld)
	}
	return flds
}
