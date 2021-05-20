package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/yi-jiayu/nationstates-secretary/nationstates"
)

type Config struct {
	Autologin string `json:"autologin"`
	Password  string `json:"password"`
	Name      string `json:"name"`
}

func getConfig() (Config, error) {
	configFile, err := os.Open("secrets.json")
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func AnswerIssue(name string, c nationstates.Client) {
	issues, err := c.GetIssues(name)
	if err != nil {
		panic(err)
	}

	// No issues to be addressed
	if len(issues) == 0 {
		shards := []string{"nextissuetime"}
		nation, err := c.GetNation(name, shards, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", nation)
		return
	}

	// Address the first issue
	issue := issues[0]
	// Pick a random option
	rand.Seed(time.Now().Unix())
	choice := rand.Intn(len(issue.Options))
	fmt.Printf("---------------ISSUE: %s-----------------\n", issue.Title)
	fmt.Printf("%s\n\n", issue.Text)
	for i, v := range issue.Options {
		if i == choice {
			fmt.Printf("(CHOSEN) ")
		}
		fmt.Printf("%d: %s\n", v.ID, v.Text)
	}
	fmt.Println()

	consequences, err := c.AnswerIssue(name, issue.ID, choice)
	if err != nil {
		panic(err)
	}
	fmt.Println("--------------------Consequences----------------------")
	if consequences.Error != "" {
		fmt.Println(consequences.Error)
	} else {
		fmt.Printf("Talking Point: %s\n\n", consequences.Desc)
		fmt.Println("Headlines:")
		for _, v := range consequences.Headlines {
			fmt.Println(v)
		}
		fmt.Println()

		// Trends are sorted by percentage change from + to -
		// We want to sort by the absolute magnitudes of the percentage change
		trends := consequences.Rankings
		sort.Slice(trends, func(i, j int) bool {
			return math.Abs(float64(trends[i].PChange)) > math.Abs(float64(trends[j].PChange))
		})
		fmt.Println("Trends:")
		for _, v := range trends {
			var direction string
			switch {
			case v.PChange > 0:
				direction = "gained"
			case v.PChange < 0:
				direction = "lost"
			}
			fmt.Printf("%s %s %.2f%% by %.2f to %.2f\n", nationstates.CensusLabels[v.ID], direction, v.PChange, v.Change, v.Score)
		}
	}
}

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	client := nationstates.Client{
		Password: config.Password,
	}

	AnswerIssue(config.Name, client)
}
