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

	appDir        string
	csvName       string
	cli           bool
	criticalError error

	sep    = ','
	prompt = false

	helpText = strings.TrimSpace(`

https://github.com/abeleinin/goki

Usage:
  goki                      - tui mode
  goki list                 - view deck index
  goki review <deck index>  - review deck from cli
		
Create:
  opt:                      - optional flags
    -n "deck name"          - assigned deck name to imported cards
    -t                      - assigns tab sep (default sep=',')

  goki opt < deck.txt       - import deck in using stdin
  goki --gpt < my_notes.txt - generate a deck of your notes using OpenAI API`)
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

	err := parseArgs(args)
	if err != nil {
		return
	}

	var response string
	if prompt {
		response = createDeckStdin()
	} else {
		response = readDeckStdin(sep)
	}

	if response != "" {
		fmt.Println(response)
		return
	}

	p := tea.NewProgram(currUser, tea.WithAltScreen())

	if _, err := p.Run(); err != nil || criticalError != nil {

		if criticalError != nil {
			err = criticalError
		}

		fmt.Println("Error running goki:", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "list":
			PrintDecks()
			// TODO: Not an error. Temp fix.
			return errors.New("")
		case "-h", "--help", "help":
			fmt.Println(gokiLogo)
			fmt.Println(helpText)
			return errors.New("")
		case "review":
			if i <= len(args)-2 {
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
		case "--gpt":
			prompt = true
			if i <= len(args)-2 {
				response := generateDeck(args[i+1])
				i++
				return errors.New(response)
			}
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
