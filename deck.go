package main

import (
	"time"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type Card struct {
  id       string 
  front    string 
  back     string 

  status   int    // 0 = new, 1 = learning, 2 = review
  reviewAt time.Time 
}

type Deck struct {
	// Deck table information
  name     string 
  new      int
  learning int
  review   int

	// Deck data
  json 		 string 
  cards    list.Model
}

func (d Deck) Name() string        { return d.name }

func (c Card) FilterValue() string { return c.front }
func (c Card) Title()       string { return c.front }
func (c Card) Description() string { return c.back }

func newDefaultDeck() *Deck {
	return NewDeck("Deck Name", "JSON Data")
}

func NewDeck(name, json string) *Deck {
	return &Deck{
		name: name,
		json: json,
		// cards: list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
}

func (d Deck) Init() tea.Cmd {
	return nil
}
func (d Deck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
    switch {
      case key.Matches(msg, keys.Quit):
			  return d, tea.Quit
      case key.Matches(msg, keys.Back):
        return sg_user.Update(nil)
      case key.Matches(msg, keys.New):
        f := newDefaultFlashcard()
        f.edit = false
        return f.Update(nil)
      case key.Matches(msg, keys.Edit):
        card := d.cards.SelectedItem().(Card)
        f := NewFlashcard(card.front, card.back)
        f.index = d.cards.Index()
        f.edit = true
        return f.Update(nil)
    }
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		d.cards.SetSize(msg.Width-h, msg.Height-v)
	}
	d.cards.SetSize(100, 50)

	d.cards, cmd = d.cards.Update(msg)
	return d, cmd
}

func (d Deck) View() string {
	return listStyle.Render(d.cards.View())
}
