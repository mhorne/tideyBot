package main

import (
	//"database/sql"
	"os"
	"os/signal"

	//"tideyBot/modules/plusPlus"
	"tideyBot/modules/soundPlayer"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	//_ "github.com/mattn/go-sqlite3"

)

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "fuck all y'all")
}

func main() {

	var (
		Token = "Bot MTcwMzI3MzkxOTA0ODU4MTEy.CglgAQ.fXXT3MkRP2B5N05pFwmByTQuUYI"
	)

	// Connect to the database
	/*db, err := sql.Open("sqlite3", "./tidey.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()*/

	// Create a discord session
	log.Info("Starting discord session...")
	discord, err := discordgo.New(Token)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord session")
		return
	}

	discord.AddHandler(onReady)

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord websocket connection")
		return
	}
	defer discord.Close()

	// We're running!
	log.Info("TideyBot is up and running :')")

	// Load Modules
	//go plusPlus.Initialize(discord, db)
	go soundPlayer.Initialize(discord)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
