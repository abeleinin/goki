package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	appDir   string
	helpText = strings.TrimSpace(`

https://github.com/abeleinin/goki

Usage: goki
  goki                        - tui mode
  goki list                   - view deck index
  goki review <deck index>    - review deck from cli`)
)

var currUser *User
var cli bool

func main() {
	runCLI(os.Args[1:])
}

func runCLI(args []string) {
	currUser = NewUser()

	loadDecks()
	initTable()
	initInput()
	updateTableColumns()

	if len(args) > 0 {
		switch args[0] {
		case "list":
			PrintDecks()
		case "-h", "--help", "help":
			fmt.Println(gokiLogo)
			fmt.Println(helpText)
		case "review":
			if len(args) > 1 {
				ReviewCLI(args[1])
			} else {
				fmt.Println("Not enough args to run 'goki review <deck index>.'")
				fmt.Println("Use 'goki list' to view deck index.")
			}
		default:
			fmt.Print(args[0], " is not a valid command. Use 'goki -h' for more information.")
		}
		return
	}

	p := tea.NewProgram(currUser, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running goki:", err)
		os.Exit(1)
	}
}

func ReviewCLI(s string) {
	i, _ := strconv.Atoi(s)
	if i < 0 || i >= len(currUser.decks) {
		fmt.Println("Invalid deck index when running 'goki review <deck index>'.")
		fmt.Println("Use 'goki list' to view deck index.")
		return
	}
	cli = true
	currUser.decks[i].StartReview()
	if len(currUser.decks[i].reviewData.reviewCards) > 0 {
		p := tea.NewProgram(currUser.decks[i])
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("No cards to review in the deck: " + currUser.decks[i].Name + ".")
	}
}
