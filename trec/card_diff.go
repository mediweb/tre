package tre

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/VojtechVitek/go-trello"
)

type labels []struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type cardDiff struct {
	Id          string
	Name        string
	ShortUrl    string
	IdListIs    string
	ListName    string
	IdListWas   string
	ListNameWas string
	NewCard     bool
	Labels      labels
}

type cardDiffs []*cardDiff

type listMap map[string]string

type sortableList []trello.List

func (l sortableList) Len() int {
	return len(l)
}

func (l sortableList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l sortableList) Less(i, j int) bool {
	return l[i].Pos < l[j].Pos
}

func labelColor(name string) string {
	switch name {
	case "green":
		return "#61BD4F"
	case "yellow":
		return "#F2D600"
	case "orange":
		return "#FFAB4A"
	default:
		return "#B6BBBF"
	}
}

func showMovedCards(boardID string, slack bool, channel string) {
	byt, err := ioutil.ReadFile("cards.json")
	if err != nil {
		log.Fatal(err)
	}

	var savedCards []trello.Card
	if err := json.Unmarshal(byt, &savedCards); err != nil {
		log.Fatal(err)
	}

	client := NewClientWithOsExitOnErr()

	var slackToken string

	if slack {
		if slackToken = os.Getenv("SLACK_INCOMING_WEBHOOK"); slackToken == "" {
			log.Fatal("SLACK_INCOMING_WEBHOOK is not set")
		}

		if channel == "" {
			log.Fatal("Slack channel name is not given")
		}
	}

	board, err := client.Board(boardID)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.Cards()
	if err != nil {
		log.Fatal(err)
	}

	var lists sortableList
	var listErr error
	lists, listErr = board.Lists()
	if err != nil {
		log.Fatal(listErr)
	}

	listMap := make(listMap)

	for _, l := range lists {
		listMap[l.Id] = l.Name
	}

	var diff cardDiffs

	if len(cards) > 0 {
		for _, cc := range cards {
			newCard := true
			for _, oc := range savedCards {
				if cc.Id == oc.Id {
					newCard = false
					if cc.IdList != oc.IdList {
						diff = append(
							diff,
							&cardDiff{
								Id: cc.Id, Name: cc.Name,
								ShortUrl:    cc.ShortUrl,
								IdListIs:    cc.IdList,
								ListName:    listMap[cc.IdList],
								IdListWas:   oc.IdList,
								ListNameWas: listMap[oc.IdList],
								Labels:      cc.Labels,
							})
					}
				}
			}
			if newCard {
				diff = append(
					diff,
					&cardDiff{
						Id: cc.Id, Name: cc.Name,
						ShortUrl: cc.ShortUrl,
						IdListIs: cc.IdList,
						ListName: listMap[cc.IdList],
						NewCard:  true,
						Labels:   cc.Labels,
					})
			}
		}
	}

	if len(diff) > 0 {
		fmt.Println("name\tlist\turl\tnew?")
		sort.Sort(sort.Reverse(lists))
		atchs := make(attachments, 0)
		for _, list := range lists {
			for _, c := range diff {
				if c.IdListIs != list.Id {
					continue
				}
				new := ""
				if c.NewCard {
					new = "New!!"
				}
				fmt.Println(c.Name, "\t", c.ListName, " <= ", c.ListNameWas, "\t", c.ShortUrl, "\t", new)
				f := make(fields, 0)
				f = append(f, field{Title: "Moved From", Value: c.ListNameWas, Short: true})
				f = append(f, field{Title: "Moved To", Value: c.ListName, Short: true})
				if len(c.Labels) > 0 {
					color := labelColor(c.Labels[0].Color)
					atchs = append(atchs, attachment{Title: c.Name, TitleLink: c.ShortUrl, Fields: f, Color: color})
				} else {
					atchs = append(atchs, attachment{Title: c.Name, TitleLink: c.ShortUrl, Fields: f, Color: "#B6BBBF"})
				}
			}
		}

		if slack {
			msg := message{
				Username:    "Trello Moved Card Report",
				IconEmoji:   ":robot_face:",
				Channel:     channel,
				Attachments: atchs,
			}
			body, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err)
			}
			_, err = http.Post(slackToken, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Println("There's no update")
	}
}
