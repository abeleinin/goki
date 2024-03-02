package main

import (
  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/textinput"
  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/table"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var (
  focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
  blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
  cursorStyle         = focusedStyle.Copy()
  noStyle             = lipgloss.NewStyle()
  helpStyle           = blurredStyle.Copy()
  cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

  docStyle = lipgloss.NewStyle().Width(100).Height(100).Align(lipgloss.Center)
)

type User struct {
  id      string
  help    help.Model
  KeyMap  keyMap
  table   table.Model
  input   textinput.Model
  decks   []*Deck
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
	help.ShowAll = false
	return &User{help: help, KeyMap: DefaultKeyMap(),}
}

func (u *User) Init() tea.Cmd {
  return nil
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch {
        case key.Matches(msg, u.KeyMap.Quit):
          if !u.input.Focused() {
            return u, tea.Quit
          }
        case key.Matches(msg, u.KeyMap.Open):
          if !u.input.Focused() {
            i := u.table.Cursor()
            u.decks[i].rdata = ReviewData{}
            return u.decks[i].Update(nil)
          }
        case key.Matches(msg, u.KeyMap.Review):
          if !u.input.Focused() {
            i := sg_user.table.Cursor()
            u.decks[i].StartReview()
            return u.decks[i].Update(nil)
          }
        case key.Matches(msg, u.KeyMap.New):
          if !u.input.Focused() {
            newDeck := NewDeck("New Deck", "new.json", []list.Item{})
            u.decks = append(u.decks, newDeck)
            u.table.SetRows(updateRows())
          }
        case key.Matches(msg, u.KeyMap.Back):
          if u.input.Focused() {
            u.input.PromptStyle = blurredStyle
            u.input.Blur()
            u.table.Focus()
            u.input.SetValue("")
          }
          return u.Update(nil)
        case key.Matches(msg, u.KeyMap.ShowFullHelp):
          fallthrough
        case key.Matches(msg, u.KeyMap.CloseFullHelp):
          if !u.input.Focused() {
            u.help.ShowAll = !u.help.ShowAll
          }
        case key.Matches(msg, u.KeyMap.Edit):
          if !u.input.Focused() {
            u.table.Blur()
            u.input.Focus()
            u.input.PromptStyle = focusedStyle
            return u, nil
          }
        case key.Matches(msg, u.KeyMap.Enter):
          if u.input.Focused() {
            s := u.input.Value()
            i := u.table.Cursor()
            u.decks[i].name = s
            u.decks[i].Cards.Title = s
            u.UpdateTable()
            u.input.Blur()
            u.table.Focus()
            u.input.SetValue("")
            u.input.PromptStyle = blurredStyle
          }
      }
    case tea.WindowSizeMsg:
      h, v := docStyle.GetFrameSize()
      docStyle = docStyle.Width(msg.Width - h).Height(msg.Height - v)
    case Form:
      i := sg_user.table.Cursor()
      if msg.edit {
        card := u.decks[i].Cards.Items()[msg.index]
        msg.EditCard(card.(*Card))
      } else {
        u.decks[i].Cards.InsertItem(0, msg.CreateCard())
        u.decks[i].UpdateStatus()
      }
      return u.decks[i].Update(nil)
    case Deck:
      i := sg_user.table.Cursor()
      u.decks[i].UpdateStatus()
      u.UpdateTable()
      return u.Update(nil)
  }

  if u.input.Focused() {
    u.input, cmd = u.input.Update(msg)
    return u, cmd
  }

  u.table, cmd = u.table.Update(msg)
  return u, cmd
}

func (u *User) View() string {
  logoStyle := lipgloss.NewStyle().
                Bold(true).
                MarginBottom(1)
  helpStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(58)

  gokiLogo := `   ________        __    __  
  /  _____/  ____ |  | _|__|
 /   \  ___ /    \|  |/ /  |
 \    \_\  |  /\  |    <|  |
  \______  /\____/|__|_ \__|
         \/            \/   `

  pageLeft := lipgloss.JoinVertical(
    lipgloss.Center,
    u.table.View(),
    helpStyle.Render(u.input.View()),
    helpStyle.Render(u.help.View(u)),
  )

  page := lipgloss.JoinVertical(
    lipgloss.Center,
    logoStyle.Render(gokiLogo),
    pageLeft,
    "",
  )
  return docStyle.Render(page)
}
