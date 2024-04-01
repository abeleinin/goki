package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type User struct {
	help       help.Model
	KeyMap     keyMap
	table      table.Model
	input      textinput.Model
	spinner    spinner.Model
	decks      []*Deck
	del        bool
	gpt        bool
	gptLoading bool
}

func (u *User) Decks() []*Deck {
	return u.decks
}

func (u *User) UpdateTable() {
	i := u.table.Cursor()
	currRows := u.table.Rows()

	rows := []table.Row{}
	for j := range currRows {
		if j == i {
			rows = append(rows, table.Row{u.decks[i].Name,
				u.decks[i].NumNew(),
				u.decks[i].NumLearning(),
				u.decks[i].NumReview()})
		} else if currRows[j] != nil {
			rows = append(rows, currRows[j])
		}
	}
	currUser.table.SetRows(rows)
}

func NewUser() *User {
	help := help.New()
	help.ShowAll = false
	spinner := spinner.New()
	spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	return &User{
		help:    help,
		KeyMap:  DefaultKeyMap(),
		del:     false,
		spinner: spinner,
	}
}

func (u *User) Init() tea.Cmd {
	return nil
}

func asyncGpt(u *User, s string) {
	u.gptLoading = true
	deck, err := gptClient(s)
	if err != nil {
		// Set critical error so that the error is printed after tui exit
		criticalError = err
		u.Update(tea.Quit)
	}
	u.decks = append(u.decks, deck)
	u.table.SetRows(updateRows())
	u.Update(nil)
	u.gptLoading = false
	u.gpt = false
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, u.KeyMap.Quit):
			if !u.input.Focused() {
				saveAll()
				return u, tea.Quit
			}
		case key.Matches(msg, u.KeyMap.Open):
			if !u.input.Focused() && len(u.decks) > 0 {
				i := u.table.Cursor()
				u.decks[i].reviewData = ReviewData{}
				return u.decks[i].Update(nil)
			}
		case key.Matches(msg, u.KeyMap.Review):
			if !u.input.Focused() && len(u.decks) > 0 {
				i := currUser.table.Cursor()
				u.decks[i].StartReview()
				if len(u.decks[i].reviewData.reviewCards) > 0 {
					return u.decks[i].Update(nil)
				} else {
					u.decks[i].reviewData = ReviewData{}
					return u.Update(nil)
				}
			}
		case key.Matches(msg, u.KeyMap.New):
			if !u.input.Focused() {
				newDeck := NewDeck("New Deck", []list.Item{})
				newDeck.NameDeckJson()
				u.decks = append(u.decks, newDeck)
				u.table.SetRows(updateRows())
				return u.Update(nil)
			}
		case key.Matches(msg, u.KeyMap.Delete):
			if !u.input.Focused() && len(u.decks) > 0 {
				u.del = true
				u.table.Blur()
				u.input.Focus()
				u.input.PromptStyle = focusedStyle
				return u, nil
			}
		case key.Matches(msg, u.KeyMap.Back):
			if u.input.Focused() {
				u.input.PromptStyle = blurredStyle
				u.input.Blur()
				u.table.Focus()
				u.input.SetValue("")
				u.del = false
				u.gpt = false
			}
			return u.Update(nil)
		case key.Matches(msg, u.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, u.KeyMap.CloseFullHelp):
			if !u.input.Focused() {
				u.help.ShowAll = !u.help.ShowAll
			}
		case key.Matches(msg, u.KeyMap.Edit):
			if !u.input.Focused() && len(u.decks) > 0 {
				u.table.Blur()
				u.input.Focus()
				u.input.PromptStyle = focusedStyle
				return u, nil
			}
		case key.Matches(msg, u.KeyMap.Gpt):
			if !u.input.Focused() && !u.gptLoading {
				u.table.Blur()
				u.input.Focus()
				u.input.PromptStyle = focusedStyle
				u.gpt = true
				return u, nil
			}
		case key.Matches(msg, u.KeyMap.Enter):
			if u.input.Focused() {
				s := u.input.Value()
				i := u.table.Cursor()
				if u.del {
					if s == "yes" {
						temp := u.table.Cursor()
						u.table.SetCursor(temp - 1)
						u.decks[i].DeleteCardsJson()
						u.decks = append(u.decks[:i], u.decks[i+1:]...)
						u.table.SetRows(updateRows())
					}
				} else if u.gpt {
					go asyncGpt(u, s)
					u.gpt = false
				} else if len(s) > 0 {
					u.decks[i].Name = s
					u.decks[i].Cards.Title = s
					u.UpdateTable()
					u.decks[i].DeleteCardsJson()
					u.decks[i].NameDeckJson()
					u.decks[i].saveCards()
				}
				saveDecks()
				u.del = false
				u.input.Blur()
				u.table.Focus()
				u.input.SetValue("")
				u.input.PromptStyle = blurredStyle
				if u.gptLoading {
					return u, u.spinner.Tick
				}
			}
		}
	case tea.WindowSizeMsg:
		screenHeight, screenWidth = msg.Height, msg.Width
		h, v := docStyle.GetFrameSize()
		docStyle = docStyle.Width(msg.Width - h).Height(msg.Height - v)
	case spinner.TickMsg:
		var cmd tea.Cmd
		u.spinner, cmd = u.spinner.Update(msg)
		return u, cmd
	case Form:
		i := currUser.table.Cursor()
		if msg.edit {
			card := u.decks[i].Cards.Items()[msg.index]
			msg.EditCard(card.(*Card))
		} else {
			u.decks[i].Cards.InsertItem(0, msg.CreateCard())
			u.decks[i].UpdateStatus()
		}
		return u.decks[i].Update(nil)
	}

	if u.input.Focused() {
		u.input, cmd = u.input.Update(msg)
		return u, cmd
	}

	u.table, cmd = u.table.Update(msg)
	return u, cmd
}

func (u *User) View() string {
	var (
		sections []string
		footer   []string
		msg      string
	)

	if u.del {
		msg = "Type 'yes' to confirm deletion:"
	} else if u.gptLoading {
		msg = u.spinner.View() + " Generating deck..."
	} else if u.gpt {
		msg = "Prompt GPT to generate a deck:"
	} else if len(u.decks) == 0 {
		msg = "No decks.\nPress 'N' to create a new deck.\nPress 'G' to generate a new deck using GPT."
	}

	footer = append(footer, homeFooterStyle.Render(msg))
	footer = append(footer, homeFooterStyle.Render(u.input.View()))
	footer = append(footer, homeFooterStyle.Render(u.help.View(u)+"\n"))
	footerStack := lipgloss.JoinVertical(lipgloss.Center, footer...)

	logoHeight := lipgloss.Height(gokiLogo)
	footerHeight := lipgloss.Height(footerStack)
	tableStyle = tableStyle.Height(screenHeight - logoHeight - footerHeight - 2)
	docStyle = docStyle.Width(screenWidth).Height(screenHeight)

	sections = append(sections, logoStyle.Render(gokiLogo))
	sections = append(sections, tableStyle.Render(u.table.View()))
	sections = append(sections, footerStack)

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Center, sections...))
}
