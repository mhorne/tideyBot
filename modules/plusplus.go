package modules

import (

	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/Sirupsen/logrus"
)

type userScore struct {
	user  	string 	`json:"username"`
	score 	int	`json:"score"`
}

type scoreCollection struct {
	scores []userScore
}

type PlusPlus struct {
	collections []scoreCollection
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

func (p *PlusPlus)fillScores(d *discordgo.Session) error {

	guilds, err := d.UserGuilds()
	if err != nil {
		return err
	}

	p.collections = make([]scoreCollection, len(guilds))

	for i := range guilds {
		g, err := d.Guild(guilds[i].ID)
		if err != nil {
			return err
		}

		s := make([]userScore, len(g.Members))

		for j := range g.Members {
			s[j].user = g.Members[j].User.Username
			s[j].score = 0
		}

		p.collections[i] = scoreCollection{s}
	}

	return err
}

// Prints all userScores in a supplied list
func (p *PlusPlus)PrintScores() {

	for i := range p.collections {
		for j := range p.collections[i].scores {
			u := p.collections[i].scores[j].user
			s := p.collections[i].scores[j].score
			fmt.Printf("User %s has a total of %d points.\n", u, s)
		}
	}
}
