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

  initCards(false)

  mathDeck := NewDeck("Math", "math.json", mathCards)
  aboutDeck := NewDeck("About Me", "about.json", aboutCards)
  quizDeck := NewDeck("Quiz", "test.json", quizCards)

  sg_user.decks = append(sg_user.decks, mathDeck)
  sg_user.decks = append(sg_user.decks, aboutDeck)
  sg_user.decks = append(sg_user.decks, quizDeck)

  initTable()

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
