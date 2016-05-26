package main

import (

	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"tideyBot/modules"
)

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

	// Add Event Handlers
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

	// Load Modules
	modules.Initialize(discord)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	discord.Close()
}