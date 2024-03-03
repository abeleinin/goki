package main

import (
  "fmt"
  "os"
  "time"
  "strconv"

  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/key"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"

  "golang.org/x/term"
)

var (
)

type (
  Status int
  Difficulty int
)

const (
  New Status = iota
  Learning
  Review
  Complete
)

const (
  Again Difficulty = iota
  Good
  Easy
)

type Card struct {
  Front        string    `json:"front"`
  Back         string    `json:"back"`

  Score        int       `json:"score"`
  Interval     int       `json:"interval"`
  EaseFactor   float64   `json:"easeFactor"`
  Status       Status    `json:"status"` 
  LastReviewed time.Time `json:"LastReviewed"`
}

func NewCard(front, back string) *Card {
  return &Card{
    Front: front,
    Back: back,
    Score: 0,
    Interval: 0,
    Status: New,
  }
}

func (c *Card) SM2(diff Difficulty) {
  switch diff {
    case Again:
      c.Score = 0
      c.Interval = 0
      c.Status = Learning
    case Good:
      if c.Score == 0 {
        c.Interval = 10
      } else {
        c.Interval = int(float64(c.Interval) * c.EaseFactor)
      }
      c.Status = Complete
      c.Score++
    case Easy:
      if c.Score == 0 {
        c.Interval = 20
      } else {
        c.Interval = int(float64(c.Interval) * c.EaseFactor)
      }
      c.Status = Complete
      c.Score++
  }
  c.EaseFactor = c.EaseFactor + 0.1 - (5 - float64(diff)) * (0.08 + (5 - float64(diff)) * 0.02)
  c.LastReviewed = time.Now()
}

type Deck struct {
  keyMap       keyMap
  help         help.Model
  descShown    bool
  resultsShown bool
  searching    bool

  numNew       int
  numLearning  int
  numReview    int
  numComplete  int

  Name         string    `json:"name"`
  Json         string     `json:"json"`
  Cards        list.Model `json:"-"`
  reviewData   ReviewData `json:"-"`
}

type ReviewData struct {
  reviewing   bool
  complete    bool
  currIx      int
  curr        *Card 
  reviewCards []*Card
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
  d.reviewData.reviewing = true
  d.reviewData.complete = false
  d.reviewData.reviewCards = d.GetReviewCards()
  d.reviewData.currIx = 0
  if len(d.reviewData.reviewCards) > 0 {
    d.reviewData.curr = d.reviewData.reviewCards[0]
  }
}

func (d *Deck) UpdateReview() {
  d.reviewData.currIx++
  d.reviewData.complete = false
}

func (d Deck) NumNew()      string { return strconv.Itoa(d.numNew) }
func (d Deck) NumLearning() string { return strconv.Itoa(d.numLearning) }
func (d Deck) NumReview()   string { return strconv.Itoa(d.numReview) }
func (d Deck) NumComplete() string { return strconv.Itoa(d.numComplete) }

func (c Card) FilterValue() string { return c.Front }
func (c Card) Title()       string { return c.Front }
func (c Card) Description() string { return c.Back }

func NewDeck(name string, jsonName string, lst []list.Item) *Deck {
  d := &Deck{
    help: help.New(),
    Name: name,
    Json: jsonName,
    Cards: list.New(lst, InitCustomDelegate(), 0, 0),
    keyMap: DeckKeyMap(),
    reviewData: ReviewData{},
  }
  d.Cards.AdditionalFullHelpKeys = func() []key.Binding {
    return []key.Binding{d.keyMap.Edit, d.keyMap.Delete, d.keyMap.New, d.keyMap.Open, d.keyMap.Save}
  }
  d.searching = false
  d.descShown = true
  d.help.ShowAll = false
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
        if !d.searching {
          return d, tea.Quit
        }
      case key.Matches(msg, d.keyMap.Back):
        if d.resultsShown {
          d.resultsShown = false
        } else {
          return currUser.Update(d)
        }
      case key.Matches(msg, d.keyMap.New):
        if !d.searching && !d.reviewData.reviewing {
          f := newDefaultForm()
          f.edit = false
          return f.Update(nil)
        }
      case key.Matches(msg, d.keyMap.Delete):
        if !d.searching && !d.reviewData.reviewing {
          d.Cards.RemoveItem(d.Cards.Index())
          return d.Update(nil)
        }
      case key.Matches(msg, d.keyMap.Save):
        if !d.searching && !d.reviewData.reviewing {
          saveCards(&d)
        }
      case key.Matches(msg, d.keyMap.Edit):
        if !d.searching && !d.reviewData.reviewing && len(d.Cards.Items()) > 0 {
          card := d.Cards.SelectedItem().(*Card)
          f := EditForm(card.Front, card.Back)
          f.index = d.Cards.Index()
          f.edit = true
          return f.Update(nil)
        }
        return d.Update(nil)
      case key.Matches(msg, d.keyMap.Search):
        if !d.searching {
          d.searching = true
        }
      case key.Matches(msg, d.keyMap.Open):
        if !d.searching {
          if d.reviewData.reviewing {
            d.reviewData.complete = true
          } else if d.descShown {
            ViewFalseDescription()
            d.descShown = !d.descShown
            d.Cards.SetDelegate(delegate)
          } else {
            ViewTrueDescription()
            d.descShown = !d.descShown
            d.Cards.SetDelegate(delegate)
          }
          return d.Update(nil)
        }
      case key.Matches(msg, d.keyMap.Easy):
        if d.reviewData.complete {
          d.reviewData.curr.SM2(Easy)
          d.UpdateReview()
        }
      case key.Matches(msg, d.keyMap.Good):
        if d.reviewData.complete {
          d.reviewData.curr.SM2(Good)
          d.UpdateReview()
        }
      case key.Matches(msg, d.keyMap.Again):
        if d.reviewData.complete {
          d.reviewData.curr.SM2(Again)
          d.UpdateReview()
        }
      case key.Matches(msg, d.keyMap.Enter):
        if d.searching {
          d.searching = false
          d.resultsShown = true
        }
    }
  case tea.WindowSizeMsg:
    h, v := listStyle.GetFrameSize()
    d.Cards.SetSize(msg.Width-h, msg.Height-v)
  }

  if d.reviewData.reviewing {
    if d.reviewData.currIx > len(d.reviewData.reviewCards) - 1 {
      d.reviewData.reviewing = false
      d.reviewData.complete = false
      i := currUser.table.Cursor()
      currUser.decks[i].UpdateStatus()
      currUser.UpdateTable()
      return currUser.Update(nil)
    } else {
      d.reviewData.curr = d.reviewData.reviewCards[d.reviewData.currIx]
    }
  }

  d.Cards.SetSize(100, 50)

  var cmd tea.Cmd
  d.Cards, cmd = d.Cards.Update(msg)
  return d, cmd
}

func (d Deck) View() string {
  fd := int(os.Stdout.Fd())
  width, height, err := term.GetSize(fd)
  if err != nil {
      fmt.Println("Error getting size:", err)
  }
  if d.reviewData.reviewing {
    cardStyle := lipgloss.NewStyle().
                  Align(lipgloss.Center).
                  Width(width).
                  Height(height)

    questStyle := lipgloss.NewStyle().
                  Bold(true).
                  Foreground(lipgloss.Color("10")).
                  Border(lipgloss.RoundedBorder()).
                  Padding(5, 20, 0, 20)

    ansStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color("15")).
                Margin(0, 0, 1, 0)

    footerStyle := lipgloss.NewStyle().MarginTop(3)

    var footer string
    if d.reviewData.complete {
      footerStyle = footerStyle.MarginTop(1)
      footer = lipgloss.JoinVertical(
        lipgloss.Center,
        ansStyle.Render(d.reviewData.curr.Back),
        helpKeyColor.Render("Card Difficulty:"),
        lipgloss.NewStyle().Inline(true).Render(d.help.View(d)),
      )
    } else {
      footer = lipgloss.JoinVertical(
        lipgloss.Center,
        d.help.View(d.reviewData.curr),
      )
    }

    ui := lipgloss.JoinVertical(
      lipgloss.Center,
      d.reviewData.curr.Front,
      footerStyle.Render(footer),
    )
    return cardStyle.Render(questStyle.Render(ui))
  } else {
    h, v := listStyle.GetFrameSize()
    listStyle = listStyle.MarginLeft(3*width/10)
    d.Cards.SetSize(width-h, height-v)
    return listStyle.Render(d.Cards.View())
  }
}
