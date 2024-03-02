package main

import (
  "fmt"
  "os"
  "time"
  "strconv"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"

  "golang.org/x/term"
)

var (
  listStyle = lipgloss.NewStyle().Margin(1, 2)
)

type Status int

const (
  New Status = iota
  Learning
  Review
  Complete
)

type Card struct {
  Front    string `json:"front"`
  Back     string `json:"back"`

  Score    int        `json:"score"`
  Status   Status     `json:"status"` 
  ReviewAt time.Time  `json:"reviewAt"`
}

func NewCard(front, back string) *Card {
  return &Card{
    Front: front,
    Back: back,
    Score: 0,
    Status: New,
    ReviewAt: time.Now(),
  }
}

type Deck struct {
  keyMap keyMap

  // Deck table information
  name        string 
  numNew      int
  numLearning int
  numReview   int
  numComplete int

  // Deck data
  json      string     
  Cards     list.Model
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
  temp := []list.Item{}
  for _, card := range d.Cards.Items() {
    if card != nil {
      c := card.(*Card)
      switch c.Status {
      case New:
          d.numNew++
      case Learning:
          d.numLearning++
      case Review:
          d.numReview++
      case Complete:
          d.numComplete++
      }
      temp = append(temp, c)
    }
  }
  d.Cards.SetItems(temp)
}

func (d *Deck) StartReview() {
  d.rdata.reviewing = true
  d.rdata.complete = false
  d.rdata.currIx = 0
  d.rdata.curr = d.Cards.Items()[d.rdata.currIx].(*Card)
}

func (d Deck) Name()        string { return d.name }
func (d Deck) NumNew()      string { return strconv.Itoa(d.numNew) }
func (d Deck) NumLearning() string { return strconv.Itoa(d.numLearning) }
func (d Deck) NumReview()   string { return strconv.Itoa(d.numReview) }
func (d Deck) NumComplete() string { return strconv.Itoa(d.numComplete) }

func (c Card) FilterValue() string { return c.Front }
func (c Card) Title()       string { return c.Front }
func (c Card) Description() string { return c.Back }

func NewDeck(name string, jsonName string, lst []list.Item) *Deck {
  d := &Deck{
    name: name,
    json: jsonName,
    Cards: list.New(lst, list.NewDefaultDelegate(), 0, 0),
    keyMap: DeckKeyMap(),
    rdata: ReviewData{},
  }
  d.UpdateStatus()
  return d
}

func (d Deck) Init() tea.Cmd {
  return nil
}

func (d Deck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch {
      case key.Matches(msg, d.keyMap.Quit):
        return d, tea.Quit
      case key.Matches(msg, d.keyMap.Back):
        return sg_user.Update(d)
      case key.Matches(msg, d.keyMap.New):
        f := newDefaultForm()
        f.edit = false
        return f.Update(nil)
      case key.Matches(msg, d.keyMap.Delete):
        d.Cards.RemoveItem(d.Cards.Index())
        return d.Update(nil)
      case key.Matches(msg, d.keyMap.Save):
        saveCards(&d)
      case key.Matches(msg, d.keyMap.Edit):
        if len(d.Cards.Items()) > 0 {
          card := d.Cards.SelectedItem().(*Card)
          f := EditForm(card.Front, card.Back)
          f.index = d.Cards.Index()
          f.edit = true
          return f.Update(nil)
        }
      case key.Matches(msg, d.keyMap.Open):
        d.rdata.complete = true
        return d.Update(nil)
      case key.Matches(msg, d.keyMap.Easy):
        if d.rdata.complete {
          d.rdata.curr.Score = 1
          d.rdata.curr.Status = Complete
          d.rdata.currIx++
          d.rdata.complete = false
        }
      case key.Matches(msg, d.keyMap.Medium):
        if d.rdata.complete {
          d.rdata.curr.Score = 0
          d.rdata.curr.Status = Learning
          d.rdata.currIx++
          d.rdata.complete = false
        }
      case key.Matches(msg, d.keyMap.Hard):
        if d.rdata.complete {
          d.rdata.curr.Score = 0
          d.rdata.curr.Status = New
          d.rdata.currIx++
          d.rdata.complete = false
        }
    }
  case tea.WindowSizeMsg:
    h, v := listStyle.GetFrameSize()
    d.Cards.SetSize(msg.Width-h, msg.Height-v)
  }

  if d.rdata.reviewing {
    if d.rdata.currIx > len(d.Cards.Items()) - 1 {
      d.rdata.reviewing = false
      d.rdata.complete = false
      i := sg_user.table.Cursor()
      sg_user.decks[i].UpdateStatus()
      sg_user.UpdateTable()
      return sg_user.Update(nil)
    } else {
      d.rdata.curr = d.Cards.Items()[d.rdata.currIx].(*Card)
    }
  }

  d.Cards.SetSize(100, 50)

  var cmd tea.Cmd
  d.Cards, cmd = d.Cards.Update(msg)
  return d, cmd
}

func (d Deck) View() string {
  if d.rdata.reviewing {

    fd := int(os.Stdout.Fd())
    width, height, err := term.GetSize(fd)
    if err != nil {
        fmt.Println("Error getting size:", err)
    }
    cardStyle := lipgloss.NewStyle().
                  Align(lipgloss.Center).
                  Width(width).
                  Height(height)

    var ui string
    questStyle := lipgloss.NewStyle().
                  Bold(true).
                  Foreground(lipgloss.Color("10")).
                  Border(lipgloss.RoundedBorder()).
                  Padding(5, 20)
    ansStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color("12"))

    if d.rdata.complete {
      ui = lipgloss.JoinVertical(
        lipgloss.Center,
        d.rdata.curr.Front,
        "",
        ansStyle.Render(d.rdata.curr.Back),
      )
    } else {
      ui = lipgloss.JoinVertical(
        lipgloss.Center,
        d.rdata.curr.Front,
        "",
        "",
      )
    }
    return cardStyle.Render(questStyle.Render(ui))
  } else {
    return listStyle.Render(d.Cards.View())
  }
}
