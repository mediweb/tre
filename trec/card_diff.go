package tre

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/VojtechVitek/go-trello"
)

type cardDiff struct {
	Id          string
	Name        string
	ShortUrl    string
	IdListIs    string
	ListName    string
	IdListWas   string
	ListNameWas string
	NewCard     bool
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

func showMovedCards(boardID string, slack bool) {
	byt, err := ioutil.ReadFile("cards.json")
	if err != nil {
		log.Fatal(err)
	}

	var savedCards []trello.Card
	if err := json.Unmarshal(byt, &savedCards); err != nil {
		log.Fatal(err)
	}

	client := NewClientWithOsExitOnErr()

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
					})
			}
		}
	}

	if len(diff) > 0 {
		fmt.Println("name\tlist\turl\tnew?")
		sort.Sort(sort.Reverse(lists))
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
			}
		}
	} else {
		fmt.Println("There's no update")
	}
}
