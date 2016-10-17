package tre

import kingpin "gopkg.in/alecthomas/kingpin.v2"

var (
	cli = kingpin.New("tre", "A handy CLI for Trello")

	board = cli.Command("boards", "List boards")

	cards        = cli.Command("cards", "List cards in a board")
	cardsBoardID = cards.Arg("board_id", "Board ID").Required().String()

	lists        = cli.Command("lists", "List lists in a board")
	listsBoardID = lists.Arg("board_id", "Board ID").Required().String()
)

func Run(args []string) {
	kingpin.CommandLine.HelpFlag.Short('h')

	switch kingpin.MustParse(cli.Parse(args[1:])) {
	case "boards":
		listBoards()
	case "cards":
		listCards(*cardsBoardID)
	case "lists":
		listLists(*listsBoardID)
	}
}