package tre

import (
	"fmt"
	"log"
)

func listBoards() {
	client := NewClientWithOsExitOnErr()
	member, err := client.Member("me")

	if err != nil {
		log.Fatal(err)
	}

	boards, err := member.Boards()
	if err != nil {
		log.Fatal(err)
	}

	if len(boards) > 0 {
		fmt.Println("name", "\t", "id")
		for _, board := range boards {
			fmt.Println(board.Name, board.Id)
		}
	}
}
