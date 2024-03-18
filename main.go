package main

import (
	"errors"
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

	// TODO: Find better solution to seperate cli and tui actions
	cli = false
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
	initGoki(os.Args[1:])

	err := readDeckStdin(sep)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !cli {
		runTUI()
	}
}

func initGoki(args []string) {
	currUser = NewUser()

	loadDecks()
	initTable()
	initInput()
	updateTableColumns()

	err := parseArgs(args)
	if err != nil {
		return
	}
}

func runTUI() {
	p := tea.NewProgram(currUser, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running goki:", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "list":
			cli = true
			PrintDecks()
		case "-h", "--help", "help":
			cli = true
			fmt.Println(gokiLogo)
			fmt.Println(helpText)
		case "review":
			if i <= len(args)-2 {
				cli = true
				ReviewCLI(args[i+1])
				i++
			} else {
				fmt.Println("Not enough args to run 'goki review <deck index>.'")
				fmt.Println("Use 'goki list' to view deck index.")
				return errors.New("Input Error")
			}
		case "-n":
			if i <= len(args)-2 {
				csvName = args[i+1]
				i++
			} else {
				fmt.Println("Please provide a deck name.")
				fmt.Println("Use 'goki help' for more info.")
				return errors.New("Input Error")
			}
		case "-t":
			sep = '\t'
		default:
			fmt.Print(args[i], " is not a valid command. Use 'goki -h' for more information.")
			return errors.New("Input Error")
		}

	}

	if sep == 0 {
		for _, arg := range args[1:] {
			if arg == "-t" {
				sep = '\t'
				break
			}
		}
	}

	return nil
}

func ReviewCLI(s string) {
	i, _ := strconv.Atoi(s)
	if i < 0 || i >= len(currUser.decks) {
		fmt.Println("Invalid deck index when running 'goki review <deck index>'.")
		fmt.Println("Use 'goki list' to view deck index.")
		return
	}
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
