package main

import (
  "strconv"
  "time"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type Status int

const (
  New Status = iota
  Learning
  Review
  Complete
)

type Card struct {
  front    string 
  back     string 

  score    int
  status   Status
  reviewAt time.Time 
}

func NewCard(front, back string) *Card {
  return &Card{
    front: front,
    back: back,
    score: 0,
    status: New,
    reviewAt: time.Now(),
  }
}

type Deck struct {
  // Deck table information
  name        string 
  numNew      int
  numLearning int
  numReview   int
  numComplete int

  // Deck data
  json     string 
  cards    list.Model
}

func (d *Deck) UpdateStatus() {
  for _, card := range d.cards.Items() {
    c := card.(*Card)
    switch c.status {
    case New:
        d.numNew++
    case Learning:
        d.numLearning++
    case Review:
        d.numReview++
    case Complete:
        d.numReview++
    }
  }
}

func (d Deck) Name()        string { return d.name }
func (d Deck) NumNew()      string { return strconv.Itoa(d.numNew) }
func (d Deck) NumLearning() string { return strconv.Itoa(d.numLearning) }
func (d Deck) NumReview()   string { return strconv.Itoa(d.numReview) }
func (d Deck) NumComplete() string { return strconv.Itoa(d.numComplete) }

func (d *Deck) NumNewInc()         { d.numNew++ }

func (c Card) FilterValue() string { return c.front }
func (c Card) Title()       string { return c.front }
func (c Card) Description() string { return c.back }

func newDefaultDeck() *Deck {
  return NewDeck("Deck Name", list.New(nil, list.NewDefaultDelegate(), 0, 0))
}

func NewDeck(name string, cards list.Model) *Deck {
  d := &Deck{
    name: name,
    cards: cards,
  }
  d.UpdateStatus()
  return d
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
        sg_user.UpdateTable()
        return sg_user.Update(nil)
      case key.Matches(msg, keys.New):
        f := newDefaultFlashcard()
        f.edit = false
        return f.Update(nil)
      case key.Matches(msg, keys.Edit):
        card := d.cards.SelectedItem().(*Card)
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
