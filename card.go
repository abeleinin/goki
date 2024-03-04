package main

import "time"

type (
	Status     int
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
	Front string `json:"front"`
	Back  string `json:"back"`

	Score        int       `json:"score"`
	Interval     int       `json:"interval"`
	EaseFactor   float64   `json:"easeFactor"`
	Status       Status    `json:"status"`
	LastReviewed time.Time `json:"LastReviewed"`
}

func (c Card) FilterValue() string { return c.Front }
func (c Card) Title() string       { return c.Front }
func (c Card) Description() string { return c.Back }

func NewCard(front, back string) *Card {
	return &Card{
		Front:    front,
		Back:     back,
		Score:    0,
		Interval: 0,
		Status:   New,
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
	c.EaseFactor = c.EaseFactor + 0.1 - (5-float64(diff))*(0.08+(5-float64(diff))*0.02)
	c.LastReviewed = time.Now()
}
