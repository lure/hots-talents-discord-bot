package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	sess, err := discordgo.New("Bot ")
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot || (!strings.HasPrefix(m.Content, "!")) {
			return
		}

		heroName := prepareName(m.Content)
		if len(FanBuilds[heroName.normalized]) == 0 {
			return
		}

		thumbnail := discordgo.MessageEmbedThumbnail{
			URL: heroName.icon,
		}

		flds := make([]*discordgo.MessageEmbedField, 0)

		for buildName, talents := range FanBuilds[heroName.normalized] {
			var buffer bytes.Buffer
			link := makePsionicTalents(heroName.normalized, talents)

			// ------------------- START talents as a list
			talentsAslist := toArrayOfNumbers(talents)
			hotsname, ok := nameToHotsName[heroName.normalized]
			if !ok {
				hotsname = heroName.normalized
			}
			heroData := TalentsDictionary[hotsname]
			for talentLevel, talentOrder := range talentsAslist {
				buffer.WriteString(fmt.Sprintf("**[%d]** %s\n", talentOrder, heroData[talentLevel][talentOrder-1]))
			}
			// ------------------- END talents as a list

			buffer.WriteString("\n")
			buffer.WriteString("[Link to build](" + link + ")")
			buffer.WriteString("\n")
			buffer.WriteString(talents)

			fld := discordgo.MessageEmbedField{
				Name:   buildName,
				Value:  buffer.String(),
				Inline: true,
			}
			flds = append(flds, &fld)
		}

		footer := discordgo.MessageEmbedFooter{
			Text: "React to the 📧 emoji to get this message sent to your DMs!",
		}

		embed := discordgo.MessageEmbed{
			Thumbnail:   &thumbnail,
			Type:        discordgo.EmbedTypeRich,
			Description: "[Support me here!](https://local.js)",
			Fields:      flds,
			Footer:      &footer,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	initStringUtils()
	parseFanGoogleSheet()
	readTalentsFromGithub()

	fmt.Println("The bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
