package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func updateRows() []table.Row {
	rows := []table.Row{}
	for _, deck := range currUser.decks {
		deck.Cards.Title = deck.Name
		rows = append(rows, table.Row{deck.Name, deck.NumNew(), deck.NumLearning(), deck.NumReview()})
	}
	return rows
}

func initTable() {
	header := []table.Column{
		{Title: "Decks", Width: 20},
		{Title: "New", Width: 10},
		{Title: "Learning", Width: 10},
		{Title: "Review", Width: 10},
	}

	rows := updateRows()

	currUser.table = table.New(
		table.WithColumns(header),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	currUser.table.SetStyles(s)
}

func initInput() {
	currUser.input = textinput.New()
	currUser.input.Placeholder = ""
	currUser.input.PromptStyle = blurredStyle
	currUser.input.CharLimit = 20
}

func saveAll() {
	saveDecks()
	for _, deck := range currUser.decks {
		deck.saveCards()
	}
}

func saveDecks() {
	jsonData, err := json.Marshal(currUser.Decks())
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("./decks/alldecks.json", jsonData, 0644)
}

func (d *Deck) saveCards() {
	jsonData, err := json.Marshal(d.Cards.Items())
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("./cards/"+d.Json, jsonData, 0644)
}

func loadDecks() {
	d := readDecks("./decks/alldecks.json")
	for _, curr := range d {
		cards := readCards("./cards/" + curr.Json)
		deck := NewDeck(curr.Name, curr.Json, cards)
		currUser.decks = append(currUser.decks, deck)
	}
}

func readDecks(fileName string) []*Deck {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var jsonDecks []Deck
	err = json.Unmarshal(byteValue, &jsonDecks)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	decks := []*Deck{}
	for _, jsonDeck := range jsonDecks {
		deck := Deck{
			Name: jsonDeck.Name,
			Json: jsonDeck.Json,
		}
		decks = append(decks, &deck)
	}

	return decks
}

func readCards(fileName string) []list.Item {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var jsonCards []Card
	err = json.Unmarshal(byteValue, &jsonCards)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	cards := []list.Item{}
	for _, jsonCard := range jsonCards {
		card := Card{
			Front:        jsonCard.Front,
			Back:         jsonCard.Back,
			Score:        jsonCard.Score,
			Interval:     jsonCard.Interval,
			EaseFactor:   jsonCard.EaseFactor,
			Status:       jsonCard.Status,
			LastReviewed: jsonCard.LastReviewed,
		}
		cards = append(cards, &card)
	}

	return cards
}

func updateTableColumns() {
	for _, deck := range currUser.decks {
		deck.GetReviewCards()
		deck.UpdateStatus()
	}
	currUser.UpdateTable()
}

func GetScreenDimensions() (int, int) {
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		log.Println("Error getting screen dimensions:", err)
	}
	return width, height
}

func (d *Deck) RenameCardsJson() {
	d.Json = NameToFilename(d.Name) + ".json"
}

func NameToFilename(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

func (d *Deck) DeleteCardsJson() {
	filePath := "./cards/" + d.Json

	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}
}
