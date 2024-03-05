package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpText = strings.TrimSpace(`

https://github.com/abeleinin/goki

Usage: goki
  goki                        - tui mode
  goki list                   - view deck index
  goki review <deck index>    - review a deck in cli`)
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
		case "-h", "--help":
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
			fmt.Print(args[0], " is not a valid command. Use 'goki -' for more information.")
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
	p := tea.NewProgram(currUser.decks[i])
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PrintDecks() {
	var section []string
	section = append(section, "\nDecks:\n")
	for i, deck := range currUser.decks {
		section = append(section, strconv.Itoa(i)+". "+deck.Name)
	}
	section = append(section, "use 'goki review <deck index>' to review a deck.\n")
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, section...))
}
