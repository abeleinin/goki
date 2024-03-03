package main

import (
  "fmt"
  "os"

  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/textarea"
  "github.com/charmbracelet/bubbles/textinput"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"

  "golang.org/x/term"
)

type Form struct {
  keyMap   keyMap

  help     help.Model
  question textinput.Model
  answer   textinput.Model

  index    int  
  edit     bool
}

func newDefaultForm() *Form {
  return NewForm("card front...", "card back...")
}

func NewForm(question, answer string) *Form {
  fc := Form{
    help:       help.New(),
    question:   textinput.New(),
    answer:     textinput.New(),
    keyMap:     FormKeyMap(),
  }
  fc.help.ShowAll = false
  fc.question.Placeholder = question
  fc.answer.Placeholder = answer
  fc.question.Focus()
  return &fc
}

func EditForm(question, answer string) *Form {
  fc := Form{
    help:       help.New(),
    question:   textinput.New(),
    answer:     textinput.New(),
    keyMap:     FormKeyMap(),
  }
  fc.help.ShowAll = false
  fc.question.SetValue(question)
  fc.answer.SetValue(answer)
  fc.question.Focus()
  return &fc
}

func (f Form) EditCard(card *Card) {
  card.Front = f.question.Value()
  card.Back = f.answer.Value()
}

func (f Form) CreateCard() *Card {
  return NewCard(f.question.Value(), f.answer.Value())
}

func (f Form) Init() tea.Cmd {
  return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, f.keyMap.Back):
      i := currUser.table.Cursor()
      return currUser.decks[i].Update(nil)
    case key.Matches(msg, f.keyMap.Enter):
      if f.question.Focused() {
        f.question.Blur()
        f.answer.Focus()
        return f, textarea.Blink
      }
      return currUser.Update(f)
    case key.Matches(msg, f.keyMap.Tab):
      if f.answer.Focused() {
        f.answer.Blur()
        f.question.Focus()
        return f, textarea.Blink
      }
    }
  case tea.WindowSizeMsg:
    h, v := promptStyle.GetFrameSize()
    promptStyle = promptStyle.Width(msg.Width - h).Height(msg.Height - v)
  }

  if f.question.Focused() {
    f.question, cmd = f.question.Update(msg)
    return f, cmd
  }

  f.answer, cmd = f.answer.Update(msg)
  return f, cmd
}

func (f Form) View() string {
  helpStyle := lipgloss.NewStyle().Align(lipgloss.Center)

  fd := int(os.Stdout.Fd())
  width, height, err := term.GetSize(fd)
  if err != nil {
      fmt.Println("Error getting size:", err)
  }

  promptStyle = promptStyle.
                Align(lipgloss.Center).
                Width(width).
                Height(height)

  viewStyle := lipgloss.NewStyle().
                Border(lipgloss.RoundedBorder()).
                Padding(2, 2, 0, 2)

  padTop := lipgloss.NewStyle().PaddingTop(2)

  prompt := lipgloss.JoinVertical(
    lipgloss.Left,
    "Create new card:",
    f.question.View(),
    f.answer.View(),
  )

  ui := lipgloss.JoinVertical(
    lipgloss.Left,
    prompt,
    padTop.Render(helpStyle.Render(f.help.View(f))),
  )

  return promptStyle.Render(viewStyle.Render(ui))
}