package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	// "github.com/yi-jiayu/nationstates-secretary/nationstates"
	"github.com/Yi-Jiahe/BotNation/nationstates"
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

func AnswerIssues(name string, c nationstates.Client) {
	issues, err := c.GetIssues(name)
	if err != nil {
		panic(err)
	}

	var log string

	if len(issues) != 0 {
		for _, issue := range issues {

			// Pick a random option
			rand.Seed(time.Now().Unix())
			choice := rand.Intn(len(issue.Options))
			log += fmt.Sprintf("---------------ISSUE: %s-----------------\n", issue.Title)
			log += fmt.Sprintf("%s\n\n", issue.Text)
			for i, v := range issue.Options {
				if i == choice {
					log += fmt.Sprintf("(CHOSEN) ")
				}
				log += fmt.Sprintf("%d: %s\n", v.ID, v.Text)
			}
			log += "\n"

			consequences, err := c.AnswerIssue(name, issue.ID, choice)
			if err != nil {
				panic(err)
			}
			log += fmt.Sprintln("--------------------Consequences----------------------")
			if consequences.Error != "" {
				log += fmt.Sprintln(consequences.Error)
			} else {
				log += fmt.Sprintf("Talking Point: %s\n\n", consequences.Desc)
				log += fmt.Sprintln("Headlines:")
				for _, v := range consequences.Headlines {
					log += fmt.Sprintln(v)
				}
				log += "\n"

				// Trends are sorted by percentage change from + to -
				// We want to sort by the absolute magnitudes of the percentage change
				trends := consequences.Rankings
				sort.Slice(trends, func(i, j int) bool {
					return math.Abs(float64(trends[i].PChange)) > math.Abs(float64(trends[j].PChange))
				})
				log += fmt.Sprintln("Trends:")
				for _, v := range trends {
					var direction string
					switch {
					case v.PChange > 0:
						direction = "gained"
					case v.PChange < 0:
						direction = "lost"
					}
					log += fmt.Sprintf("%s %s %.2f%% by %.2f to %.2f\n", nationstates.CensusLabels[v.ID], direction, v.PChange, v.Change, v.Score)
				}
				log += fmt.Sprintln()
			}
		}
	}

	shards := []string{"nextissuetime"}
	nation, err := c.GetNation(name, shards, nil)
	if err != nil {
		panic(err)
	}
	t := time.Unix(nation.NextIssueTime, 0)
	if err != nil {
		panic(err)
	}
	log += fmt.Sprintf("Next issue at %v\n", t)

	fmt.Print(log)
}

func PerformCensus(name string, c nationstates.Client) {
	shard := c.CreateCensusShard(nil, nil)
	n, err := c.GetNation(name, []string{shard}, nil)
	if err != nil {
		panic(err)
	}
	scales := n.Scales
	if len(scales) == 0 {
		fmt.Println("No scales returned")
		return
	}
	var log string

	log += fmt.Sprintln("--------------------------National Census-----------------------")
	for _, v := range scales {
		log += fmt.Sprintf("------%s-------\n", nationstates.CensusLabels[v.ID])
		if len(v.Points) != 0 {
			log += fmt.Sprintln("I don't really want to print that")
		} else {
			if v.Score != 0 {
				log += fmt.Sprintf("Score: %.2f\n", v.Score)
			}
			if v.Rank != 0 {
				log += fmt.Sprintf("Ranked %dth in the world\n", v.Rank)
			}
			if v.PRank != 0 {
				log += fmt.Sprintf("%d%% in the world\n", v.PRank)
			}
			if v.RRank != 0 {
				log += fmt.Sprintf("Ranked %dth in the region\n", v.RRank)
			}
			if v.PRRank != 0 {
				log += fmt.Sprintf("%d%% in the region\n", v.PRRank)
			}
		}
	}
	fmt.Print(log)
}

func HandleRequest() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	client := nationstates.Client{
		Password: config.Password,
	}

	AnswerIssues(config.Name, client)
	PerformCensus(config.Name, client)

	return
}

func main() {
	HandleRequest()
}
