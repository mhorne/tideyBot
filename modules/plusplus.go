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

func New(d *discordgo.Session) *PlusPlus {

	p := new(PlusPlus)

	err := p.fillScores(d)
	if err != nil {
		logrus.Error(err)
		logrus.Error("PlusPlus module was not initialized!")
	}

	return p
}

func (p *PlusPlus)HandleMessage(m *discordgo.Message) {

}

// This method iterates through the guilds and their members
// to create a table of scores TODO: Clean up
func (p *PlusPlus)fillScores(d *discordgo.Session) error {

	guilds, err := d.UserGuilds()
	if err != nil {
		return err
	}

	p.guildList = make([]string, len(guilds))
	p.collections = make(map[string]scoreCollection)

	for i := range guilds {
		g, err := d.Guild(guilds[i].ID)
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

		fmt.Println(sc)
		p.collections[p.guildList[i]] = *sc
	}

	return err
}

// Prints all userScores in a supplied list; only for testing right now
func (p *PlusPlus)PrintScores() {

	fmt.Println("Printing?")

	for i := range p.guildList {
		for j := range p.collections[p.guildList[i]].userList {

			u := p.collections[p.guildList[i]].userList[j]
			s := p.collections[p.guildList[i]].scores[u]
			fmt.Printf("User %s has a total of %d points.\n", u, s)
		}
	}
}
