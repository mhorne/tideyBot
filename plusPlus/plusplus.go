package plusPlus

import (

	"github.com/bwmarrin/discordgo"
	"fmt"
)

type userScore struct {
	user  	string 	`json:"username"`
	score 	int	`json:"score"`
}

type ScoreCollection struct {
	scores []userScore
}

func FillScores(d *discordgo.Session) (*ScoreCollection, error) {

	var (
		sc = ScoreCollection{}
		guilds []*discordgo.Guild
		g *discordgo.Guild
	)

	guilds, err := d.UserGuilds()
	g, err = d.Guild(guilds[0].ID)

	if err == nil {

		fmt.Printf("%#v,guild", g)
		s := make([]userScore, len(g.Members))

		//Copy all members usernames
		for i := range g.Members {
			s[i].user = g.Members[i].User.Username
			s[i].score = 0
		}

		fmt.Printf("%#v,score", sc)

		sc := ScoreCollection{s}
		return &sc, err
	} else {
		return &sc, err
	}
}

// Prints all userscores in a supplied list
func (s *ScoreCollection) PrintScores() {

	for i:= range s.scores {
		fmt.Printf("User %s has a total of %d points.\n", s.scores[i].user, s.scores[i].score)
	}
}
