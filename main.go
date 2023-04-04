package main

import (
	"flag"
	"log"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/bwmarrin/discordgo"
	"github.com/rssnyder/discord-stock-ticker/utils"
)

func main() {
	token := flag.String("token", "", "discord bot token")
	nickname := flag.Bool("nickname", true, "change bot nickname")
	activity := flag.String("activity", "", "bot activity")
	status := flag.Int("status", 0, "0: playing, 1: listening")
	refresh := flag.Int("refresh", 300, "seconds between refresh")
	flag.Parse()

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		log.Println(err)
		*nickname = false
	}
	if len(guilds) == 0 {
		*nickname = false
	}

	ticker := time.NewTicker(time.Duration(*refresh) * time.Second)

	for {
		select {
		case <-ticker.C:
			circulating, err := GetCirculating()
			if err != nil {
				log.Printf("Error getting circulating data %s\n", err)
				continue
			}

			priceData, err := utils.GetCryptoPrice("dsquared-finance")
			if err != nil {
				log.Printf("Error getting price data %s\n", err)
				continue
			}

			marketcapRaw := priceData.MarketData.CurrentPrice.USD * circulating

			p := message.NewPrinter(language.English)
			marketcap := p.Sprintf("$%.2f\n", marketcapRaw)

			if *nickname {
				for _, g := range guilds {
					err = dg.GuildMemberNickname(g.ID, "@me", marketcap)
					if err != nil {
						log.Println(err)
						continue
					} else {
						log.Printf("Set nickname in %s: %s\n", g.Name, marketcap)
					}
				}
			} else {
				switch *status {
				case 0:
					err = dg.UpdateGameStatus(0, marketcap)
				case 1:
					err = dg.UpdateListeningStatus(marketcap)
				}
				if err != nil {
					log.Printf("Unable to set activity: %s\n", err)
				} else {
					log.Printf("Set activity: %s\n", marketcap)
				}
			}
			if *activity != "" {
				switch *status {
				case 0:
					err = dg.UpdateGameStatus(0, *activity)
				case 1:
					err = dg.UpdateListeningStatus(*activity)
				}
				if err != nil {
					log.Printf("Unable to set activity: %s\n", err)
				} else {
					log.Printf("Set activity: %s\n", *activity)
				}
			}
		}
	}
}
