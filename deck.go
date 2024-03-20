package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type ReviewData struct {
	reviewing   bool
	complete    bool
	fillIn      bool
	currIx      int
	curr        *Card
	reviewCards []*Card
}

type Deck struct {
	keyMap         keyMap
	help           help.Model
	progress       progress.Model
	fillInResponse textinput.Model
	descShown      bool
	resultsShown   bool
	searching      bool

	numNew      int
	numLearning int
	numReview   int
	numComplete int

	Name         string     `json:"name"`
	Json         string     `json:"json"`
	Cards        list.Model `json:"-"`
	reviewData   ReviewData `json:"-"`
	deletedCards []*Card
}

func (d Deck) NumNew() string      { return strconv.Itoa(d.numNew) }
func (d Deck) NumLearning() string { return strconv.Itoa(d.numLearning) }
func (d Deck) NumReview() string   { return strconv.Itoa(d.numReview) }
func (d Deck) NumComplete() string { return strconv.Itoa(d.numComplete) }

func (d *Deck) StartReview() {
	d.reviewData.reviewing = true
	d.reviewData.fillIn = false
	d.reviewData.complete = false
	d.reviewData.reviewCards = d.GetReviewCards()
	d.reviewData.currIx = 0
	if len(d.reviewData.reviewCards) > 0 {
		d.reviewData.curr = d.reviewData.reviewCards[0]
	}
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

func (d *Deck) GetReviewCards() []*Card {
	var (
		timeNow = time.Now()

		c           *Card
		duration    time.Duration
		minutes     float64
		reviewCards []*Card
	)

	for _, card := range d.Cards.Items() {
		if card != nil {
			c = card.(*Card)
			if c.Status == New {
				reviewCards = append(reviewCards, c)
			} else {
				duration = timeNow.Sub(c.LastReviewed)
				minutes = math.Floor(duration.Minutes())
				if minutes >= float64(c.Interval) {
					reviewCards = append(reviewCards, c)
					if c.Status == Complete {
						c.Status = Review
					}
				}
			}
		}
	}

	rand.Shuffle(len(reviewCards), func(i, j int) {
		reviewCards[i], reviewCards[j] = reviewCards[j], reviewCards[i]
	})

	return reviewCards
}

func (d *Deck) UpdateReview() {
	d.reviewData.currIx++
	d.reviewData.complete = false
}

func NewDeck(name string, lst []list.Item) *Deck {
	d := &Deck{
		help:       help.New(),
		progress:   progress.New(),
		Name:       name,
		Cards:      list.New(lst, InitCustomDelegate(), 0, 0),
		keyMap:     DeckKeyMap(),
		reviewData: ReviewData{},
	}
	d.initCardInput()
	d.progress.ShowPercentage = false
	d.NameDeckJson()
	d.Cards.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{d.keyMap.New, d.keyMap.Edit, d.keyMap.Delete, d.keyMap.Undo, d.keyMap.Quit}
	}
	d.Cards.SetSize(screenWidth-40, screenHeight-4)
	d.searching = false
	d.descShown = true
	d.help.ShowAll = false
	d.UpdateStatus()
	return d
}

func (d *Deck) saveCards() {
	jsonData, err := json.Marshal(d.Cards.Items())
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(appDir+"/cards/"+d.Json, jsonData, 0644)
}

func (d *Deck) NameDeckJson() {
	id := uuid.New()
	d.Json = fmt.Sprintf("%s_%s%s", NameToFilename(d.Name), id, ".json")
}

func (d *Deck) DeleteCardsJson() {
	filePath := appDir + "/cards/" + d.Json

	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}
}

func (d Deck) Init() tea.Cmd {
	return nil
}

func (d *Deck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keyMap.Quit):
			if !d.searching && !d.fillInResponse.Focused() {
				saveAll()
				return d, tea.Quit
			}
		case key.Matches(msg, d.keyMap.Back):
			if !d.fillInResponse.Focused() {
				if cli {
					saveAll()
					return d, tea.Quit
				} else if d.resultsShown {
					d.resultsShown = false
				} else {
					d.UpdateStatus()
					currUser.UpdateTable()
					return currUser.Update(nil)
				}
			}
		case key.Matches(msg, d.keyMap.New):
			if !d.searching && !d.reviewData.reviewing {
				f := newDefaultForm()
				f.edit = false
				return f.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Delete):
			if !d.searching && !d.reviewData.reviewing {
				d.deletedCards = append(d.deletedCards, d.Cards.Items()[d.Cards.Index()].(*Card))
				d.Cards.RemoveItem(d.Cards.Index())
				d.UpdateStatus()
				return d.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Undo):
			size := len(d.deletedCards)
			if size > 0 && !d.searching && !d.reviewData.reviewing {
				d.Cards.InsertItem(0, d.deletedCards[size-1])
				d.deletedCards = d.deletedCards[:size-1]
				return d.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Edit):
			if !d.searching && !d.reviewData.reviewing && len(d.Cards.Items()) > 0 {
				card := d.Cards.SelectedItem().(*Card)
				f := EditForm(card.Front, card.Back)
				f.index = d.Cards.Index()
				f.edit = true
				return f.Update(nil)
			}
			if !d.fillInResponse.Focused() {
				return d.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Search):
			if !d.searching {
				d.searching = true
			}
		case key.Matches(msg, d.keyMap.Open):
			if !d.searching && len(d.Cards.Items()) > 0 && !d.fillInResponse.Focused() {
				if d.reviewData.reviewing {
					d.reviewData.complete = true
				}
				return d.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Fill):
			if !d.searching && len(d.Cards.Items()) > 0 && !d.fillInResponse.Focused() {
				if d.reviewData.reviewing {
					d.reviewData.fillIn = true
					d.fillInResponse.Focus()
				}
				return d.Update(nil)
			}
		case key.Matches(msg, d.keyMap.Easy):
			if d.reviewData.complete && !d.reviewData.fillIn {
				d.reviewData.curr.SM2(Easy)
				d.UpdateReview()
			}
		case key.Matches(msg, d.keyMap.Good):
			if d.reviewData.complete && !d.reviewData.fillIn {
				d.reviewData.curr.SM2(Good)
				d.UpdateReview()
			}
		case key.Matches(msg, d.keyMap.Again):
			if d.reviewData.complete && !d.reviewData.fillIn {
				d.reviewData.curr.SM2(Again)
				d.UpdateReview()
			}
		case key.Matches(msg, d.keyMap.Enter):
			if d.reviewData.complete {
				d.reviewData.fillIn = false
				d.fillInResponse.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
				d.fillInResponse.SetValue("")
				d.UpdateReview()
			} else if d.fillInResponse.Focused() {
				s := d.fillInResponse.Value()
				if strings.ToLower(s) == strings.ToLower(d.reviewData.curr.Back) {
					d.fillInResponse.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04b52d"))
					d.reviewData.curr.SM2(Easy)
				} else {
					d.fillInResponse.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
					d.reviewData.curr.SM2(Again)
				}
				d.reviewData.complete = true
				d.fillInResponse.Blur()
				d.fillInResponse.PromptStyle = blurredStyle
			}
			if d.searching {
				d.searching = false
				d.resultsShown = true
			}
		}
	case tea.WindowSizeMsg:
		screenWidth, screenHeight = msg.Width, msg.Height
		cardStyle = cardStyle.MarginLeft(3 * screenWidth / 10).MarginTop(screenHeight / 10).
			Width(2 * screenWidth / 5).Height(screenHeight / 5)
	}

	if d.reviewData.reviewing {
		if d.fillInResponse.Focused() {
			var cmd tea.Cmd
			d.fillInResponse, cmd = d.fillInResponse.Update(msg)
			return d, cmd
		}
		if d.reviewData.currIx > len(d.reviewData.reviewCards)-1 {
			if cli {
				saveAll()
				return d, tea.Quit
			}
			d.reviewData.reviewing = false
			d.reviewData.complete = false
			d.UpdateStatus()
			currUser.UpdateTable()
			return currUser.Update(nil)
		} else {
			d.reviewData.curr = d.reviewData.reviewCards[d.reviewData.currIx]
		}
	}

	var cmd tea.Cmd
	d.Cards, cmd = d.Cards.Update(msg)
	return d, cmd
}

func (d Deck) View() string {
	if d.reviewData.reviewing {
		var sections []string

		front := WrapString(d.reviewData.curr.Front, 40)
		sections = append(sections, front)

		if d.reviewData.complete {
			back := WrapString(d.reviewData.curr.Back, 40)
			sections = append(sections, answerStyle.Render(back))
			if !d.reviewData.fillIn {
				sections = append(sections, helpKeyColor.Render("Card Difficulty:"))
				sections = append(sections, lipgloss.NewStyle().Inline(true).Render(d.help.View(d)))
			} else {
				sections = append(sections, promptCompleteStyle.Render(d.fillInResponse.View()))
			}
		} else if d.reviewData.fillIn {
			sections = append(sections, promptFooterStyle.Render(d.fillInResponse.View()))
		} else {
			sections = append(sections, deckFooterStyle.Render(d.help.View(d.reviewData.curr)))
		}

		progress := float64(d.reviewData.currIx) / float64(len(d.reviewData.reviewCards))
		sections = append(sections, progressStyle(d.progress.ViewAs(progress)))

		page := questionStyle.Render(lipgloss.JoinVertical(lipgloss.Center, sections...))

		if screenWidth < 100 {
			cardStyle = cardStyle.MarginLeft(1 * screenWidth / 10).MarginTop(screenHeight / 10).
				Width(4 * screenWidth / 5).Height(screenHeight / 5)
		} else {
			cardStyle = cardStyle.MarginLeft(3 * screenWidth / 10).MarginTop(screenHeight / 10).
				Width(2 * screenWidth / 5).Height(screenHeight / 5)
		}

		if cli {
			// display card in upper right corner of CLI
			cardStyle = cardStyle.Margin(0, 0, 1)
		}

		return cardStyle.Render(page)
	}

	listStyle = listStyle.Align(lipgloss.Left).MarginLeft((screenWidth - 60) / 2)

	return listStyle.Render(d.Cards.View())
}
