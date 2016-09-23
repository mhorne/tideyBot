package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

const MAX_INCREASE = 3
const MAX_DECREASE = -3

type scoreCollection struct {
	userList []string
	scores   map[string]int
}

type PlusPlus struct {
	guildList   []string
	collections map[string]scoreCollection
}

// Create a new instance of PlusPlus
func Initialize(s *discordgo.Session) {

	p := new(PlusPlus)

	err := p.fillScores(s)
	if err != nil {
		logrus.Error(err)
		logrus.Error("PlusPlus module was not initialized!")
		return
	}

	// Add message event handler to the discord session
	s.AddHandler(p.HandleMessage)

	logrus.Info("Initialized PlusPlus module")
	return
}

// Message handler method to be invoked by the discordgo session
func (p *PlusPlus) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(m.Mentions) <= 0 {
		return
	}

	// Use a simple state machine to parse the incoming messages
	// Looks for " ++ " and " -- ",  and longer segments of these
	var state int = 0
	var mod int = 0

	for _, j := range m.Content {
		switch state {
		case 0:
			if strings.Compare(string(j), " ") == 0 {
				state = 1
			}
		case 1:
			if strings.Compare(string(j), "+") == 0 {
				state = 2
			} else if strings.Compare(string(j), "-") == 0 {
				state = 3
			}
		case 2:
			if strings.Compare(string(j), "+") == 0 {
				state = 4
				mod = 1
			} else {
				state = 0
			}
		case 3:
			if strings.Compare(string(j), "-") == 0 {
				state = 5
				mod = -1
			} else {
				state = 0
			}
		case 4:
			if strings.Compare(string(j), "+") == 0 {
				mod++
			} else if strings.Compare(string(j), " ") == 0 {
				break
			} else {
				state = 0
				mod = 0
			}
		case 5:
			if strings.Compare(string(j), "-") == 0 {
				mod--
			} else if strings.Compare(string(j), " ") == 0 {
				break
			} else {
				state = 0
				mod = 0
			}
		}
	}

	messageChan, err := s.Channel(m.ChannelID)
	guild := messageChan.GuildID
	if err != nil {
		logrus.Error(err)
	}

	//Loop through all mentioned users and update their score
	if mod != 0 {
		for i := range m.Mentions {
			p.modifyScore(s, guild, m.ChannelID, m.Mentions[i].Username, mod)
		}
	}

	return
}

//Function to modify an existing user's score
func (p *PlusPlus) modifyScore(s *discordgo.Session, guild string, channel string, user string, mod int) {

	//Cap the amount of points a user can gain or lose at once
	if mod > MAX_INCREASE {
		mod = MAX_INCREASE
	} else if mod < MAX_DECREASE {
		mod = MAX_DECREASE
	}

	//Update the score
	p.collections[guild].scores[user] += mod
	newScore := p.collections[guild].scores[user]

	// Send message to the channel
	var message string
	if mod >= 0 {
		message = "Nice! " + user + " just gained " + strconv.Itoa(mod) + " points. They now have a total of " + strconv.Itoa(newScore) + "!"
	} else {
		mod = -mod
		message = "Ouch! " + user + " just lost " + strconv.Itoa(mod) + " points. They now have a total of " + strconv.Itoa(newScore) + "!"
	}

	s.ChannelMessageSend(channel, message)

	return
}

// This method iterates through the guilds and their members
// to create a table of scores TODO: Clean up
func (p *PlusPlus) fillScores(s *discordgo.Session) error {

	guilds, err := s.UserGuilds()
	if err != nil {
		return err
	}

	p.guildList = make([]string, len(guilds))
	p.collections = make(map[string]scoreCollection)

	for i := range guilds {
		g, err := s.Guild(guilds[i].ID)
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
func (p *PlusPlus) printScores() {

	for i := range p.guildList {
		for j := range p.collections[p.guildList[i]].userList {

			u := p.collections[p.guildList[i]].userList[j]
			s := p.collections[p.guildList[i]].scores[u]
			fmt.Printf("User %s has a total of %d points.\n", u, s)
		}
	}
}
