package plusPlus

import (
	"os"
	"io/ioutil"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/Sirupsen/logrus"
	"fmt"
)

var (
	filename string = "scores.js"
)

type userScore struct {
	User  	string 	`json:"username"`
	Score 	int	`json:"score"`
}

func checkErr(e error) {
	if e != nil {
		logrus.Error(e)
	}
}

func populateScores(s []userScore) {

	//Some code to read userscores from the current file
	b, err := ioutil.ReadFile("scores.js")
	json.Unmarshal(b, &s)
	checkErr(err)

	for i := range s {
		fmt.Print(i)
		fmt.Printf((s[i].User))
	}
}

func writeToFile() {

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
}

func modifyScore(user discordgo.User, modifier int) {

	//Iterate through collection, find matching userscore
	//Modify userscore

	//Write userscore to file
	//Print new score
}

func StartPlusPlus() {

	writeToFile()
}
