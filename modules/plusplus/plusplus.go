package plusplus

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Module_Name = "PlusPlus"
)

var (
	max_increase = 3
	max_decrease = -3
)

// scoreCollection contains a list of
// usernames mapped to a list of scores
type scoreCollection struct {
	userList []string
	scores   map[string]int
}

type plusPlus struct {
	session     *discordgo.Session
	db          *sql.DB
	guildList   []string
	collections map[string]scoreCollection
}

func GetModuleName() string {
	return Module_Name
}

// Create a new instance of PlusPlus
func Initialize(s *discordgo.Session, db *sql.DB) {

	p := new(plusPlus)
	p.session = s
	p.db = db

	// Populate the map of scores
	err := p.checkDB()
	if err != nil {
		log.Error(err)
		log.Error("PlusPlus module was not initialized!")
		return
	}

	// Add message event handler to the discord session
	s.AddHandler(p.HandleMessage)

	log.Info("Initialized PlusPlus module")
	return
}

// Message handler method to be invoked by the discordgo session
func (p *plusPlus) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(m.Mentions) <= 0 {
		return
	}

	// Use a simple state machine to parse the incoming messages
	// Looks for " ++ " and " -- ",  and longer segments of these
	var state int = 0
	var mod int = 0

	for _, i := range m.Content {
		switch state {
		case 0:
			if strings.Compare(string(i), " ") == 0 {
				state = 1
			}
		case 1:
			if strings.Compare(string(i), "+") == 0 {
				state = 2
			} else if strings.Compare(string(i), "-") == 0 {
				state = 3
			}
		case 2:
			if strings.Compare(string(i), "+") == 0 {
				state = 4
				mod = 1
			} else {
				state = 0
			}
		case 3:
			if strings.Compare(string(i), "-") == 0 {
				state = 5
				mod = -1
			} else {
				state = 0
			}
		case 4:
			if strings.Compare(string(i), "+") == 0 {
				mod++
			} else if strings.Compare(string(i), " ") == 0 {
				break
			} else {
				state = 0
				mod = 0
			}
		case 5:
			if strings.Compare(string(i), "-") == 0 {
				mod--
			} else if strings.Compare(string(i), " ") == 0 {
				break
			} else {
				state = 0
				mod = 0
			}
		}
	}

	messageChan, err := p.session.Channel(m.ChannelID)
	guildID := messageChan.GuildID
	if err != nil {
		log.Error(err)
		return
	}

	// Loop through all mentioned users and update their score
	tx, err := p.db.Begin()

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
func (p *plusPlus) modifyScore(tx *sql.Tx, channelID string, guildID string, userID string, mod int) {

	// Cap the amount of points a user can gain or lose at once
	if mod > max_increase {
		mod = max_increase
	} else if mod < max_decrease {
		mod = max_decrease
	}

	guild, err := p.session.Guild(guildID)
	if err != nil {
		log.Error(err)
	}
	user, err := p.session.User(userID)
	if err != nil {
		log.Error(err)
	}

	// Check if has an exisiting score in the database
	var score int
	query := "SELECT score FROM scores WHERE guild_id=? AND user_id=?"
	err = p.db.QueryRow(query, guildID, userID).Scan(&score)

	if err != nil {
		score = mod

		stmt, err := tx.Prepare("INSERT INTO SCORES VALUES (?, ?, ?, ?, ?)")
		_, err = stmt.Exec(guildID, userID, guild.Name, user.Username, score)
		if err != nil {
			log.Error("Score could not be updated")
			log.Error(err)
			return
		}
		stmt.Close()
	} else {
		stmt, err := tx.Prepare("UPDATE SCORES SET SCORE=? WHERE GUILD_ID=? AND USER_ID=?")
		_, err = stmt.Exec(score+mod, guildID, userID)
		if err != nil {
			log.Error("Score could not be updated")
			log.Error(err)
			return
		}
		stmt.Close()
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

func (p *plusPlus) checkDB() error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS SCORES (	GUILD_ID INTEGER NOT NULL,
													USER_ID TEXT NOT NULL,
													GUILD_NAME TEXT,
													USER_NAME TEXT,
													SCORE INTEGER NOT NULL
													PRIMARY KEY (GUILD_ID, USER_ID));`

	_, err := p.db.Exec(sqlStmt)

	return err
}

// This method iterates through the guilds and their members
// to create a table of scores TODO: Clean up
func (p *plusPlus) fillScores() error {

	guilds, err := p.session.UserGuilds()
	if err != nil {
		return err
	}

	p.guildList = make([]string, len(guilds))
	p.collections = make(map[string]scoreCollection)

	for i := range guilds {
		g, err := p.session.Guild(guilds[i].ID)
		if err != nil {
			return err
		}

		sc := new(scoreCollection)
		p.guildList[i] = g.ID

		// Fill collection with usernames and scores
		sc.userList = make([]string, len(g.Members))
		sc.scores = make(map[string]int)

		for j := range g.Members {
			sc.userList[j] = g.Members[j].User.Username
			sc.scores[g.Members[j].User.Username] = 0 // TODO: Grab value from file
		}

		p.collections[p.guildList[i]] = *sc
	}

	return err
}

// Prints all userScores in a supplied list; only for testing right now
func (p *plusPlus) printScores() {

	for i := range p.guildList {
		for j := range p.collections[p.guildList[i]].userList {

			u := p.collections[p.guildList[i]].userList[j]
			s := p.collections[p.guildList[i]].scores[u]
			fmt.Printf("User %s has a total of %d points.\n", u, s)
		}
	}
}
