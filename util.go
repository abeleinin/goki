package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/table"
  "github.com/charmbracelet/lipgloss"
)

func initTable() {
  header := []table.Column{
    {Title: "Decks", Width: 20},
    {Title: "New", Width: 10},
    {Title: "Learning", Width: 10},
    {Title: "Review", Width: 10},
  }

  rows := []table.Row{}
  for _, deck := range sg_user.decks {
    deck.Cards.Title = deck.Name()
    rows = append(rows, table.Row{deck.Name(), 
                                  deck.NumNew(), 
                                  deck.NumLearning(),
                                  deck.NumReview()})
  }

  sg_user.table = table.New(
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

  sg_user.table.SetStyles(s)
}

var mathCards []list.Item
var aboutCards []list.Item
var quizCards []list.Item

func initCards(saveCards bool) {
  if saveCards {
    mathCards = readCards("./cards/math.json")
    aboutCards = readCards("./cards/about.json")
    quizCards = readCards("./cards/test.json")
  } else {
    mathCards = []list.Item{
      NewCard("What's 2 * 2?", "4"),
      NewCard("What's 3 * 3?", "9"),
      NewCard("What's 4 * 4?", "16"),
      NewCard("What's 5 * 5?", "25"),
      NewCard("What's 6 * 6?", "36"),
    }
    aboutCards = []list.Item{
      NewCard("What's my name?", "Goki"),
      NewCard("What's my favorite Language?", "Go :)"),
    }
    quizCards = []list.Item{
      NewCard("What's JSON?", "JavaScript Object Notation"),
      NewCard("What's a struct?", "A collection of fields"),
      NewCard("What's a pointer?", "A memory address"),
    }
  }
}

func saveCards(d *Deck) {
  jsonData, err := json.Marshal(d.Cards.Items())
  if err != nil {
      log.Fatal(err)
  }
  err = os.WriteFile("./cards/" + d.json, jsonData, 0644)
}

func readCards(fileName string) []list.Item {
  file, err := os.Open(fileName)
  if err != nil {
      log.Fatalf("Error opening file: %s", err)
  }
  defer file.Close()

  byteValue, err := ioutil.ReadAll(file)
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
      Front: jsonCard.Front,
      Back: jsonCard.Back,
      Score: jsonCard.Score,
      Status: jsonCard.Status,
      ReviewAt: jsonCard.ReviewAt,
    }
    cards = append(cards, &card)
  }

  return cards
}