package main

import (
  "fmt"
  "os"

  "github.com/charmbracelet/bubbles/textinput"
  "github.com/charmbracelet/bubbles/table"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "github.com/charmbracelet/bubbles/list"
)

var sg_user *User

func main() {
  sg_user = NewUser()

  cards := []list.Item{
    Card{front: "2 * 5", back: "10"},
    Card{front: "50 * 3", back: "150"},
    Card{front: "7 * 5", back: "35"},
    Card{front: "9 * 9", back: "81"},
  }
  cards2 := []list.Item{
    Card{front: "What's my name?", back: "Abe Leininger"},
    Card{front: "How old are you?", back: "22"},
  }

  deck := Deck{
    name: "Math Deck",
    cards: list.New(cards, list.NewDefaultDelegate(), 0, 0),
  }

  deck2 := Deck{
    name: "About Me Deck",
    cards: list.New(cards2, list.NewDefaultDelegate(), 0, 0),
  }

  sg_user.decks = append(sg_user.decks, deck)
  sg_user.decks = append(sg_user.decks, deck2)

  header := []table.Column{
    {Title: "Deck", Width: 20},
    {Title: "New", Width: 10},
    {Title: "Learning", Width: 10},
    {Title: "Due", Width: 10},
  }

  rows := []table.Row{}
  for _, deck := range sg_user.decks {
    rows = append(rows, table.Row{deck.Name(), "0", "0", "0"})
  }

  // data := []table.Row{
  //   {"Deck 1", "10", "20", "30"},
  //   {"Deck 2", "40", "50", "60"},
  //   {"Deck 3", "400", "5", "12"},
  //   {"Deck 4", "8", "7", "10"},
  //   {"Deck 5", "4", "53", "62"},
  // }

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

  sg_user.input = make([]textinput.Model, 1)

  var txt textinput.Model
  for i := range sg_user.input {
    txt = textinput.New()
    txt.CursorStyle = cursorStyle
    txt.CharLimit = 32

    switch i {
      case 0:
        txt.PromptStyle = focusedStyle
        txt.TextStyle = focusedStyle
    }

    sg_user.input[0] = txt
  }

  p := tea.NewProgram(sg_user, tea.WithAltScreen())

  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}
