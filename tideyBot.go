package main

import (
	"os"
	"os/signal"

	"tideyBot/modules"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	logrus.Info("Recieved READY payload")
	s.UpdateStatus(0, "fuck all y'all")
}

func main() {

	var (
		Token = "Bot MTcwMzI3MzkxOTA0ODU4MTEy.CglgAQ.fXXT3MkRP2B5N05pFwmByTQuUYI"
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

	discord.AddHandler(onReady)

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
	modules.Init(discord)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	discord.Close()
}
