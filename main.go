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
    NewCard("What's 2 * 2?", "4"),
    NewCard("What's 3 * 3?", "9"),
    NewCard("What's 4 * 4?", "16"),
    NewCard("What's 5 * 5?", "25"),
    NewCard("What's 6 * 6?", "36"),
  }
  cards2 := []list.Item{
    NewCard("What's my name?", "Abe Leininger"),
    NewCard("What's my favorite color?", "Blue"),
    NewCard("How old are you?", "22"),
  }

  deck := NewDeck("Math Deck", 
                  list.New(cards, list.NewDefaultDelegate(), 0, 0))

  deck2 := NewDeck("About Me Deck", 
                  list.New(cards2, list.NewDefaultDelegate(), 0, 0))

  sg_user.decks = append(sg_user.decks, deck)
  sg_user.decks = append(sg_user.decks, deck2)

  header := []table.Column{
    {Title: "Deck", Width: 20},
    {Title: "New", Width: 10},
    {Title: "Learning", Width: 10},
    {Title: "Review", Width: 10},
  }

  rows := []table.Row{}
  for _, deck := range sg_user.decks {
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
