package modules

import (

	"github.com/bwmarrin/discordgo"
	"fmt"
)

/*type userScore struct {
	user  	string 	`json:"username"`
	score 	int	`json:"score"`
}*/

type ScoreCollection struct {
	userList []string
	scores map[string]int
}

func FillScores(d *discordgo.Session) (*ScoreCollection, error) {

	var sc = new(ScoreCollection)

	guilds, err := d.UserGuilds()
	if err != nil {
		return sc, err
	}

	g, err := d.Guild(guilds[0].ID)
	if err != nil {
		return sc, err
	}

	// Fill collection with usernames and scores
	sc.userList = make([]string, len(g.Members))
	sc.scores = make(map[string]int)

	for i := range g.Members {
		sc.userList[i] = g.Members[i].User.Username
		sc.scores[g.Members[i].User.Username] = 0 // TODO: Grab value from file
	}

	return sc, err
}

// Prints all userscores in a supplied list
func (s *ScoreCollection) PrintScores() {

	for i := range s.userList {
		fmt.Printf("User %s has a total of %d points.\n", s.userList[i], s.scores[s.userList[i]])
	}
}
