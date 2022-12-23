package main

import (
	"fmt"
	"github.com/astralservices/go-dalle"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Bot Variables and DALL-E clients
var (
	token      = "MTA1NTc4MTAxMzU3OTYzMjY4Mg.GwuCrX.x7639_esF1KsTHn6P9Ii5qktLy4WBd4mDbZlP8"
	botPrefix  = "!"
	openApiKey = "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"
	client     = dalle.NewClient(openApiKey)
)

func main() {

	// Creation of new discord session
	discordBot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	// Specifying Handler for message creation events
	discordBot.AddHandler(messageCreate)

	// Specifying Event types that we want to listen to
	discordBot.Identify.Intents = discordgo.IntentsGuildMessages

	// Opening the connection and waking up the bot
	err = discordBot.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println("Bot is running!!")

	// Listening for interruptions or stop commands
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Closing the connection
	err = discordBot.Close()
	if err != nil {
		panic(err)
	}
}

// This function is called every time a message is sent on any channel that the bot can view
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all the messages by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Message Content exploded by space
	args := strings.Fields(m.Content)

	// Ping Reply
	if args[0] == botPrefix+"ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Pong reply
	if args[0] == botPrefix+"pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Image Generation by DALL-E
	if args[0] == botPrefix+"image" {

		// Conversion of Slice to String
		prompt := strings.Join(args[1:], " ")

		// Generation of Image
		data, err := client.Generate(prompt, nil, nil, nil, nil)
		if err != nil {
			_, err = s.ChannelMessageSend(m.ChannelID, "Unable to generate image!")
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(err)
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, data[0].URL)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
