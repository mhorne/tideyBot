package modules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/Sirupsen/logrus"
)

type scoreCollection struct {
	userList []string
	scores map[string]int
}

type PlusPlus struct {
	guildList []string
	collections map[string]scoreCollection
}

// Create a new instance of PlusPlus
func Initialize(s *discordgo.Session) {

	p := new(PlusPlus)

	err := p.fillScores(s)
	if err != nil {
		logrus.Error(err)
		logrus.Error("PlusPlus module was not initialized!")
	}

	// Add message event handler to the discord session
	s.AddHandler(p.HandleMessage)

	logrus.Info("Initialized PlusPlus module")

	return
}

// Message handler method to be invoked by the discordgo session
func (p *PlusPlus)HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(m.Mentions) > 0 {
		// TODO: add the parsing logic here
	}

	return
}

// This method iterates through the guilds and their members
// to create a table of scores TODO: Clean up
func (p *PlusPlus)fillScores(s *discordgo.Session) error {

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
func (p *PlusPlus)printScores() {

	for i := range p.guildList {
		for j := range p.collections[p.guildList[i]].userList {

			u := p.collections[p.guildList[i]].userList[j]
			s := p.collections[p.guildList[i]].scores[u]
			fmt.Printf("User %s has a total of %d points.\n", u, s)
		}
	}
}
