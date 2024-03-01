package main

import (
  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/textarea"
  "github.com/charmbracelet/bubbles/textinput"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

type Form struct {
  help        help.Model
  question    textinput.Model
  answer      textinput.Model

  index       int  
  edit        bool
}

var promptStyle = lipgloss.NewStyle().Width(100).Align(lipgloss.Center).MarginTop(10)

func newDefaultForm() *Form {
  return NewForm("Write Question Here...", "Answer Here...")
}

func NewForm(question, answer string) *Form {
  fc := Form{
    help:       help.New(),
    question:   textinput.New(),
    answer:     textinput.New(),
  }
  fc.question.Placeholder = question
  fc.answer.Placeholder = answer
  fc.question.Focus()
  return &fc
}

func (f Form) EditCard(card *Card) {
  card.front = f.question.Value()
  card.back = f.answer.Value()
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
    case key.Matches(msg, keys.Back):
      i := sg_user.table.Cursor()
      return sg_user.decks[i].Update(nil)
    case key.Matches(msg, keys.Enter):
      if f.question.Focused() {
        f.question.Blur()
        f.answer.Focus()
        return f, textarea.Blink
      }
      return sg_user.Update(f)
    case key.Matches(msg, keys.Tab):
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
  prompt := lipgloss.JoinVertical(
    lipgloss.Top,
    "Create a new Form:",
    f.question.View(),
    f.answer.View(),
    f.help.View(keys),
  )

  return promptStyle.Render(prompt)
}