package plusPlus

import (

	"github.com/bwmarrin/discordgo"
	"github.com/Sirupsen/logrus"
	"fmt"
)

var (
	filename string = "scores.js"
)

type userScore struct {
	user  	string 	`json:"username"`
	score 	int	`json:"score"`
}

type ScoreCollection struct {
	scores []userScore
}

func checkErr(e error) {
	if e != nil {
		logrus.Error(e)
	}
}

/*func populateScores(s []userScore) {

	//Some code to read userscores from the current file
	b, err := ioutil.ReadFile("scores.js")
	json.Unmarshal(b, &s)
	checkErr(err)

	for i := range s {
		fmt.Print(i)
		fmt.Printf((s[i].User))
	}
}*/

/*func writeToFile() {

	//Some code to update the file from current userscore
	f, err := os.Open(filename)
	checkErr(err)

	test := userScore{
		User: "Mitcho",
		Score: 1,
	}

	b, err := json.Marshal(test)

	fmt.Print(b)

	f.Write(b)
	f.Sync()
	f.Close()
}*/

func modifyScore(user discordgo.User, modifier int) {

	//Iterate through collection, find matching userscore
	//Modify userscore

	//Write userscore to file
	//Print new score
}

func FillScores(discordgo.Guild) *ScoreCollection{

	var s ScoreCollection = ScoreCollection{}

	return &s
}

func (s *ScoreCollection) PrintScores(n int) {

	if n != nil {
		if n > len(s.scores) {
			n = len(s.scores)
		}
	} else {
		n = len(s.scores)
	}

	for i := 0; i < n; i++ {
		fmt.Printf("User %s has a total of %d points.\n", s.scores[i].user, s.scores[i].score)
	}
}
