package main

import (
  "fmt"
  "os"

  tea "github.com/charmbracelet/bubbletea"
)

var currUser *User

func main() {
  currUser = NewUser()

  loadDecks()
  initTable()
  initInput()
  updateTableColumns()

  p := tea.NewProgram(currUser, tea.WithAltScreen())

  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}
