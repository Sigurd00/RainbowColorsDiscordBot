package bot

import (
	"RainbowColorsDiscordBot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

var BotId string
var goBot *discordgo.Session
var doneSignal = make(chan bool)
var ticker = time.NewTicker(10 * time.Second)

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running !")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	args := strings.Fields(m.Content)
	if args[0] == "!start" {
		if TooFewArguments(args, 2) {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Please specify the role you want to rainbow", m.Reference())
			if err != nil {
				fmt.Printf(err.Error())
			}
			return
		}
		roles, _ := s.GuildRoles(m.GuildID)
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(args[1]) {
				StartRainbow(s, m, role, ticker)
				break
			}
		}
	} else if args[0] == "!stop" {
		if TooFewArguments(args, 2) {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Please specify the role you want to stop", m.Reference())
			if err != nil {
				fmt.Printf(err.Error())
			}
			return
		}
		roles, _ := s.GuildRoles(m.GuildID)
		for _, role := range roles {
			if strings.ToLower(role.Name) == strings.ToLower(args[1]) {
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Stopping changing colors on role %s", role.Name))
				if err != nil {
					fmt.Printf(err.Error())
				}
				ticker.Stop()
				doneSignal <- true
				break
			}
		}
	}
}

func StartRainbow(s *discordgo.Session, m *discordgo.MessageCreate, role *discordgo.Role, ticker *time.Ticker) {
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Started changing colors on role %s", role.Name))
	if err != nil {
		fmt.Printf(err.Error())
	}
	go func() {
		rainbowColors := [7]int{
			0xff0000,
			0xffa500,
			0xffff00,
			0x008000,
			0x0000ff,
			0x4b0082,
			0xee82ee,
		}
		var colorIndex int = 0
		for {
			select {
			case <-doneSignal:
				return
			case t := <-ticker.C:
				fmt.Println(t)
				ChangeRoleColor(s, m, role, rainbowColors[colorIndex])
				colorIndex++
				if colorIndex >= len(rainbowColors) {
					colorIndex = 0
				}
			}
		}
	}()
}

func TooFewArguments(args []string, countIntRequired int) bool {
	return len(args) < countIntRequired
}

func ChangeRoleColor(s *discordgo.Session, m *discordgo.MessageCreate, role *discordgo.Role, color int) {
	role, err := s.GuildRoleEdit(m.GuildID, role.ID, role.Name, color, role.Hoist, role.Permissions, role.Mentionable)
	fmt.Println(role.Color)
	if err != nil {
		fmt.Println(err.Error())
	}
}
