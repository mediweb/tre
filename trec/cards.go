package tre

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func listCards(boardID string, save bool) {
	client := NewClientWithOsExitOnErr()

	board, err := client.Board(boardID)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.Cards()
	if err != nil {
		log.Fatal(err)
	}

	if len(cards) > 0 {
		fmt.Println("name", "\t", "list", "url")
		for _, card := range cards {
			fmt.Println(card.Name, card.IdList, card.ShortUrl)
		}

		if save {
			json, _ := json.Marshal(cards)
			err := ioutil.WriteFile("./cards.json", json, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
