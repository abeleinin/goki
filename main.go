package main

import (
  "fmt"
  "os"

  "github.com/charmbracelet/bubbles/textinput"
  tea "github.com/charmbracelet/bubbletea"
)

var sg_user *User

func main() {
  sg_user = NewUser()

  initCards(true)

  mathDeck := NewDeck("Math", "math.json", mathCards)
  aboutDeck := NewDeck("About Me", "about.json", aboutCards)
  quizDeck := NewDeck("Quiz", "test.json", quizCards)

  sg_user.decks = append(sg_user.decks, mathDeck)
  sg_user.decks = append(sg_user.decks, aboutDeck)
  sg_user.decks = append(sg_user.decks, quizDeck)

  initTable()

  sg_user.input = textinput.New()
  sg_user.input.Placeholder = ""
  sg_user.input.PromptStyle = blurredStyle
  sg_user.input.CharLimit = 20
  p := tea.NewProgram(sg_user, tea.WithAltScreen())

  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}
