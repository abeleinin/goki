package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	keyMap keyMap

	help     help.Model
	question textarea.Model
	answer   textarea.Model

	index int
	edit  bool
}

func newDefaultForm() *Form {
	return NewForm("card front...", "card back...")
}

func NewForm(question, answer string) *Form {
	fc := Form{
		help:     help.New(),
		question: textarea.New(),
		answer:   textarea.New(),
		keyMap:   FormKeyMap(),
		edit:     false,
	}
	fc.question.ShowLineNumbers = false
	fc.answer.ShowLineNumbers = false
	fc.question.SetHeight(3)
	fc.answer.SetHeight(3)
	fc.help.ShowAll = false
	fc.question.Placeholder = question
	fc.answer.Placeholder = answer
	fc.question.Focus()
	return &fc
}

func EditForm(question, answer string) *Form {
	fc := Form{
		help:     help.New(),
		question: textarea.New(),
		answer:   textarea.New(),
		keyMap:   FormKeyMap(),
		edit:     false,
	}
	fc.question.ShowLineNumbers = false
	fc.answer.ShowLineNumbers = false
	fc.question.SetHeight(3)
	fc.answer.SetHeight(3)
	fc.help.ShowAll = false
	fc.question.SetValue(question)
	fc.answer.SetValue(answer)
	fc.question.Focus()
	return &fc
}

func (f Form) EditCard(card *Card) {
	front := WrapString(f.question.Value(), 40)
	back := WrapString(f.answer.Value(), 40)
	card.Front = front
	card.Back = back
}

func (f Form) CreateCard() *Card {
	front := WrapString(f.question.Value(), 40)
	back := WrapString(f.answer.Value(), 40)
	return NewCard(front, back)
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
				return f.Update(nil)
			}
			return currUser.Update(f)
		case key.Matches(msg, f.keyMap.Tab):
			if f.answer.Focused() {
				f.answer.Blur()
				f.question.Focus()
				return f.Update(nil)
			}
		}
	case tea.WindowSizeMsg:
		screenWidth, screenHeight = msg.Width, msg.Height
	}

	if f.question.Focused() {
		f.question, cmd = f.question.Update(msg)
		return f, cmd
	}

	f.answer, cmd = f.answer.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	var sections []string

	if f.edit {
		sections = append(sections, formTitleStyle("Edit Card"))
	} else {
		sections = append(sections, formTitleStyle("Create Card"))
	}
	sections = append(sections, pad("Front:"))
	sections = append(sections, f.question.View())
	sections = append(sections, pad("Back:"))
	sections = append(sections, f.answer.View())
	sections = append(sections, formFooterStyle.Render(f.help.View(f)))

	if screenWidth < 100 {
		viewStyle = viewStyle.Width(9 * screenWidth / 10)
		formStyle = formStyle.Margin(screenHeight/10, screenWidth/20, 0, screenWidth/20)
	} else {
		viewStyle = viewStyle.Width(screenWidth / 2)
		formStyle = formStyle.Margin(screenHeight/10, screenWidth/4, 0, screenWidth/4)
	}

	return formStyle.Render(viewStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...)))
}
