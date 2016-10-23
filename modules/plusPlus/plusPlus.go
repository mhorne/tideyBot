package plusPlus

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Name of the tideyBot module
	moduleName = "PlusPlus"
)

var (
	maxIncrease = 3
	maxDecrease = -3
)

// scoreCollection contains a list of
// usernames mapped to a list of scores
type scoreCollection struct {
	userList []string
	scores   map[string]int
}

//PlusPlus is a struct of the PlusPlus module
type PlusPlus struct {
	session *discordgo.Session
	db      *sql.DB
}

// GetModuleName returns the name of the module
func (p *PlusPlus) GetModuleName() string {
	return moduleName
}

// Initialize PlusPlus module
func Initialize(s *discordgo.Session, db *sql.DB) {

	p := new(PlusPlus)
	p.session = s
	p.db = db

	// Check that the database is set up
	err := p.checkDB()
	if err != nil {
		log.Error("PlusPlus was not initialized!")
		return
	}

	// Add message event handler to the discord session
	p.session.AddHandler(p.handleMessage)

	log.Info("Initialized PlusPlus module")
	return
}

// Message handler method to be invoked by the discordgo session
func (p *PlusPlus) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(m.Mentions) <= 0 {
		return
	}

	// Use a regular expression to check the message
	// for instances of " ++ " or " -- "
	re, err := regexp.Compile(`\s([++]+|[--]+)([\s]|\z)`)
	match := re.FindString(m.Content)

	// If we have a match, continue
	if match == "" {
		return
	}

	// Find channel and guild IDs
	messageChan, err := p.session.Channel(m.ChannelID)
	guildID := messageChan.GuildID
	if err != nil {
		log.Error(err)
		return
	}

	// Determine the amount of points users will gain or lose
	var mod int
	if strings.Contains(match, "+") {
		mod = strings.Count(match, "+") - 1
	} else {
		mod = -strings.Count(match, "-") + 1
	}

	// Open a new SQL transaction
	tx, err := p.db.Begin()
	if err != nil {
		log.Error(err)
		return
	}

	// Loop through all mentioned users and update their score
	if mod != 0 {
		for i := range m.Mentions {
			p.modifyScore(tx, m.ChannelID, guildID, m.Mentions[i].ID, mod)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
	}
}

// Function to modify an existing user's score
func (p *PlusPlus) modifyScore(tx *sql.Tx, channelID string, guildID string, userID string, mod int) {

	// Cap the amount of points a user can gain or lose at once
	if mod > maxIncrease {
		mod = maxIncrease
	} else if mod < maxDecrease {
		mod = maxDecrease
	}

	// Check if has an exisiting score in the database
	var (
		score int
		stmt  *sql.Stmt
		err   error
	)

	query := "SELECT score FROM scores WHERE guild_id=? AND user_id=?"
	err = p.db.QueryRow(query, guildID, userID).Scan(&score)

	// Add SQL statement to transaction either
	// updating an entry or inserting a new one
	if err != nil {
		score = mod

		stmt, err = tx.Prepare("INSERT INTO scores(guild_id, user_id, score) VALUES(?, ?, ?)")
		if err != nil {
			log.Error(err)
			return
		}

		_, err = stmt.Exec(guildID, userID, score)
		if err != nil {
			log.Error("Score could not be updated")
			log.Error(err)
			return
		}
		stmt.Close()
	} else {
		score += mod

		stmt, err = tx.Prepare("UPDATE scores SET score=? WHERE guild_id=? AND user_id=?")
		_, err = stmt.Exec(score, guildID, userID)
		if err != nil {
			log.Error("Score could not be updated")
			log.Error(err)
			return
		}
		stmt.Close()
	}

	// Get user data
	user, err := p.session.User(userID)
	if err != nil {
		log.Error(err)
		return
	}

	// Send message to the channel
	var message string
	if mod >= 0 {
		message = fmt.Sprintf("Nice! %s just gained %d points. They now have a total of %d!", user.Username, mod, score)
	} else {
		mod = -mod
		message = fmt.Sprintf("Ouch! %s just lost %d points. They now have a total of %d!", user.Username, mod, score)
	}

	p.session.ChannelMessageSend(channelID, message)
}

func (p *PlusPlus) checkDB() error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS scores (guild_id TEXT NOT NULL, user_id TEXT NOT NULL, score INTEGER NOT NULL, PRIMARY KEY (guild_id, user_id));`
	_, err := p.db.Exec(sqlStmt)

	return err
}

func (p *PlusPlus) findLeaderBoard(lim int) []scoreCollection {
	sqlStmt := `SELECT user_id, score
	            FROM scores
				ORDER BY score DESC
				LIMIT ?`
	rows, err := p.db.Query(sqlStmt, lim)
	if err != nil {
		return nil
	}

	rows.Close()
	return nil
}
