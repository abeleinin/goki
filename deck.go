package main

import (
  "time"
  "strconv"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var (
  listStyle = lipgloss.NewStyle().Margin(1, 2)
  cardStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

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
  json      string 
  cards     list.Model
  rdata     ReviewData
}

type ReviewData struct {
  reviewing bool
  complete  bool
  curr      *Card 
  currIx    int
}

func (d *Deck) UpdateStatus() {
  d.numNew, d.numLearning, d.numReview, d.numComplete = 0, 0, 0, 0
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
        d.numComplete++
    }
  }
}

func (d *Deck) StartReview() {
  d.rdata.reviewing = true
  d.rdata.complete = false
  d.rdata.currIx = 0
  d.rdata.curr = d.cards.Items()[d.rdata.currIx].(*Card)
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
        f := newDefaultForm()
        f.edit = false
        return f.Update(nil)
      case key.Matches(msg, keys.Edit):
        card := d.cards.SelectedItem().(*Card)
        f := NewForm(card.front, card.back)
        f.index = d.cards.Index()
        f.edit = true
        return f.Update(nil)
      case key.Matches(msg, keys.Open):
        d.rdata.complete = true
        return d.Update(nil)
      case key.Matches(msg, keys.Easy):
        if d.rdata.complete {
          d.rdata.curr.score = 1
          d.rdata.curr.status = Complete
          d.rdata.currIx++
          d.rdata.complete = false
        }
      case key.Matches(msg, keys.Medium):
        if d.rdata.complete {
          d.rdata.curr.score = 0
          d.rdata.curr.status = Learning
          d.rdata.currIx++
          d.rdata.complete = false
        }
      case key.Matches(msg, keys.Hard):
        if d.rdata.complete {
          d.rdata.curr.score = 0
          d.rdata.curr.status = New
          d.rdata.currIx++
          d.rdata.complete = false
        }
    }
  case tea.WindowSizeMsg:
    h, v := cardStyle.GetFrameSize()
    cardStyle = cardStyle.Width(msg.Width - h).Height(msg.Height - v)
    d.cards.SetSize(msg.Width-h, msg.Height-v)
  }

  if d.rdata.currIx > len(d.cards.Items()) - 1 {
    d.rdata.reviewing = false
    d.rdata.complete = false
    i := sg_user.table.Cursor()
    sg_user.decks[i].UpdateStatus()
    sg_user.UpdateTable()
    return sg_user.Update(nil)
  } else {
    d.rdata.curr = d.cards.Items()[d.rdata.currIx].(*Card)
  }

  h, v := cardStyle.GetFrameSize()
  cardStyle = cardStyle.MarginLeft(h/2).Height(v/2)

  d.cards.SetSize(100, 50)

  d.cards, cmd = d.cards.Update(msg)
  return d, cmd
}

func (d Deck) View() string {
  if d.rdata.reviewing {
    var ui string
    questStyle := lipgloss.NewStyle().
                  Bold(true).
                  Foreground(lipgloss.Color("10")).
                  Border(lipgloss.RoundedBorder()).
                  MarginTop(10).
                  Padding(5, 20)

    if d.rdata.complete {
      ui = lipgloss.JoinVertical(
        lipgloss.Left,
        questStyle.Render(d.rdata.curr.front),
        "",
        d.rdata.curr.back,
      )
    } else {
      ui = lipgloss.JoinVertical(
        lipgloss.Center,
        questStyle.Render(d.rdata.curr.front),
        "",
        "",
      )
    }
    return cardStyle.Render(ui)
  }
  return listStyle.Render(d.cards.View())
}
