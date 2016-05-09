package main

import (

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"tideyBot/plusPlus"
	"os/signal"
	"os"
)

func messageParser(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message recieved.")

	//Check the message and decide what type of command it is
	//Will be empty for now, until modules that use this are added
	if m.Content[0] == '!' && len(m.Content) > 1 {
		switch m.Content {
		}
	}

	//PlusPlus doesn't use a '!' command, so check for that

	if len(m.Mentions) > 0 {

	}

	return
}

func main() {

	var (
		Token = "MTcwMzI3MzkxOTA0ODU4MTEy.CglgAQ.fXXT3MkRP2B5N05pFwmByTQuUYI"
	)

	// Create a discord session
	logrus.Info("Starting discord session...")
	discord, err := discordgo.New(Token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create discord session")
		return
	}

	//Load Modules
	plusPlus.StartPlusPlus()

	discord.AddHandler(messageParser)

	err = discord.Open()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create discord websocket connection")
		return
	}

	// We're running!
	logrus.Info("TideyBot is up and running :')")

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	discord.Close()
}
