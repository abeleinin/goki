package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	currUser *User

	appDir  string
	csvName string
	cli     bool

	sep = ','

	helpText = strings.TrimSpace(`

https://github.com/abeleinin/goki

Usage:
  goki                      - tui mode
  goki list                 - view deck index
  goki review <deck index>  - review deck from cli
		
Create:
  opt:                 - optional flags
    -n "deck name"     - assigned deck name to imported cards
    -t                 - assigns tab sep (default sep=',')

  goki opt < deck.txt  - import deck in using stdin`)
)

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
		case "-n":
			if len(args) > 1 {
				csvName = args[1]
			} else {
				fmt.Println("Please provide a deck name.")
				fmt.Println("Use 'goki help' for more info.")
			}
		case "-t":
			{
				sep = '\t'
			}
		default:
			fmt.Print(args[0], " is not a valid command. Use 'goki -h' for more information.")
		}

		if len(args) < 2 {
			return
		}

		if sep == 0 {
			for _, arg := range args[1:] {
				if arg == "-t" {
					sep = '\t'
					break
				}
			}
		}
	}

	response := readDeckStdin(sep)

	if response != "" {
		fmt.Println(response)
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
