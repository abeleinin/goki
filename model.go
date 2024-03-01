package main

import (
  "encoding/json"
  // "io/ioutil"
  "log"
  "os"

  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/table"
  "github.com/charmbracelet/bubbles/textinput"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var (
  focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
  blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
  cursorStyle         = focusedStyle.Copy()
  noStyle             = lipgloss.NewStyle()
  helpStyle           = blurredStyle.Copy()
  cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

  docStyle = lipgloss.NewStyle().Width(100).Height(100).Align(lipgloss.Center)
)

type User struct {
  id      string
  help    help.Model
  table   table.Model
  input   []textinput.Model
  decks   []*Deck // table -> decks
}

func (u *User) UpdateTable() {
  i := u.table.Cursor()
  currRows := u.table.Rows()
  
  rows := []table.Row{}
  for j, _ := range currRows {
    if j == i {
      rows = append(rows, table.Row{u.decks[i].Name(), 
                                    u.decks[i].NumNew(), 
                                    u.decks[i].NumLearning(),
                                    u.decks[i].NumReview()})
    } else {
      rows = append(rows, currRows[j])
    }
  }
  sg_user.table.SetRows(rows)
}

func NewUser() *User {
	help := help.New()
	help.ShowAll = true
	return &User{help: help}
}

func (u *User) Init() tea.Cmd {
  return nil
}

func (u User) writeJSON() {
  var result map[string]map[string]string

  if result == nil {
    result = make(map[string]map[string]string)
  }
  if result["Questions"] == nil {
      result["Questions"] = make(map[string]string)
  }

  // for _, card := range u.list.Items() {
  //   result["Questions"][card.(Card).question]= card.(Card).answer
  // }

  jsonData, err := json.Marshal(result)
  if err != nil {
    log.Fatalf("Error serializing to JSON: %s", err)
  }

  err = os.WriteFile("result.json", jsonData, 0644)
  if err != nil {
    log.Fatalf("Error writing JSON to file: %s", err)
  }
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch {
        case key.Matches(msg, keys.Quit):
          return u, tea.Quit
        case key.Matches(msg, keys.Enter):
          i := u.table.Cursor()
          return u.decks[i].Update(nil)
        // case key.Matches(msg, keys.Save):
        //   u.writeJSON()
      }
    case tea.WindowSizeMsg:
      h, v := docStyle.GetFrameSize()
      docStyle = docStyle.Width(msg.Width - h).Height(msg.Height - v)
    case Flashcard:
      i := sg_user.table.Cursor()
      card := u.decks[i].cards.SelectedItem()
      if msg.edit {
        u.decks[i].cards.SetItem(msg.index, msg.EditCard(card.(*Card)))
      } else {
        u.decks[i].cards.InsertItem(0, msg.CreateCard())
        u.decks[i].NumNewInc()
      }
      return u.decks[i].Update(nil)
  }

  cmd = u.updateInputs(msg)

  u.table, cmd = u.table.Update(msg)

  return u, cmd
}

func (u *User) View() string {
  logoStyle := lipgloss.NewStyle().
  Bold(true).
  Foreground(lipgloss.Color("0"))

  gokiLogo := `   ________        __    __  
  /  _____/  ____ |  | _|__|
 /   \  ___ /    \|  |/ /  |
 \    \_\  |  /\  |    <|  |
  \______  /\____/|__|_ \__|
         \/            \/   `

  pageLeft := lipgloss.JoinVertical(
    lipgloss.Left,
    u.table.View(),             // Render the table
    u.help.View(keys),          // Render the help
  )

  page := lipgloss.JoinVertical(
    lipgloss.Center,            // Center page
    logoStyle.Render(gokiLogo), // Render the logo
    pageLeft,
  )
  return docStyle.Render(page)
}

func (u *User) updateInputs(msg tea.Msg) tea.Cmd {
  cmds := make([]tea.Cmd, len(u.input))

  for i := range u.input {
    u.input[i], cmds[i] = u.input[i].Update(msg)
  }

  return tea.Batch(cmds...)
}

// func processJSON() []list.Item {
//   file, err := os.Open("result.json")
//   if err != nil {
//       log.Fatalf("Error opening file: %s", err)
//   }
//   defer file.Close()

//   byteValue, err := ioutil.ReadAll(file)
//   if err != nil {
//       log.Fatalf("Error reading file: %s", err)
//   }

//   var result map[string]map[string]string

//   err = json.Unmarshal(byteValue, &result)
//   if err != nil {
//       log.Fatalf("Error parsing JSON: %s", err)
//   }

//   cards := []list.Item{}
//   for q, a := range result["Questions"] {
//     cards = append(cards, Card{question: q, answer: a})
//   }

//   return cards
// }
