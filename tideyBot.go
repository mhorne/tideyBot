package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"tideyBot/modules/plusPlus"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

// Location of the main tideyBot configuration file
const ConfigPath = "./tidey.conf"

// Config is a struct to hold the configuration info
type Config struct {
	Token    string `toml:"token"`
	DBPath   string `toml:"database_path"`
	LogLevel string `toml:"log_level"`
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Recieved READY payload")
	s.UpdateStatus(0, "fuck all y'all")
}

func configure() (*Config, error) {
	config := new(Config)

	// Attempt to read configuration info from file
	_, err := toml.DecodeFile(ConfigPath, &config)
	if err != nil {
		return config, err
	}

	// Set log level
	config.LogLevel = strings.ToUpper(config.LogLevel)

	switch config.LogLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	}

	// Bot tokens should be of the form "Bot ..."
	if !strings.Contains(config.Token, "Bot ") {
		config.Token = "Bot " + config.Token
	}

	fmt.Println(config)
	log.Info("Configuration loaded successfully")
	return config, nil
}

func main() {
	// Read configuration data from file
	config, err := configure()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Create a discord session
	log.Info("Starting discord session...")
	discord, err := discordgo.New(config.Token)
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
	go plusPlus.Initialize(discord, db)
	//go soundPlayer.Initialize(discord)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
