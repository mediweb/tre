package tre

import (
	"fmt"
	"log"
)

func listCards(boardID string) {
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
		fmt.Println("name", "\t", "url")
		for _, card := range cards {
			fmt.Println(card.Name, card.Url)
		}
	}
}
