package tre

import (
	"fmt"
	"log"
)

func listLists(boardID string) {
	client := NewClientWithOsExitOnErr()

	board, err := client.Board(boardID)
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	if len(lists) > 0 {
		fmt.Println("name", "\t", "id", "position")
		for _, list := range lists {
			fmt.Println(list.Name, list.Id, list.Pos)
		}
	}
}
