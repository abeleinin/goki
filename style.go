package main

import (
  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/lipgloss"
)

const (
  bullet   = "•"
  ellipsis = "…"
)

var (
  delegate *list.DefaultDelegate

  // Styles
  h = help.New()
  helpKeyColor  = h.Styles.ShortKey.Inline(true)
  helpSep       = h.Styles.ShortSeparator.Inline(true).Render(h.ShortSeparator)
  helpDescColor = h.Styles.ShortDesc.Inline(true)
)

func InitCustomDelegate() list.DefaultDelegate {
  temp := list.DefaultDelegate{}
  temp.ShowDescription = true
  temp.Styles = CustomItemStyles()
  temp.SetHeight(2)
  temp.SetSpacing(1)
  delegate = &temp
  return temp
}

func ViewTrueDescription() {
  delegate.ShowDescription = true
}

func ViewFalseDescription() {
  delegate.ShowDescription = false 
}

func CustomItemStyles() (s list.DefaultItemStyles) {
  s.NormalTitle = helpKeyColor.
    Padding(0, 0, 0, 1)

  s.NormalDesc = helpDescColor.
    Padding(0, 0, 0, 1)

  s.SelectedTitle = lipgloss.NewStyle().
    Bold(true).
    Border(lipgloss.NormalBorder(), false, false, false, true).
    BorderForeground(lipgloss.Color("2")).
    Foreground(lipgloss.Color("2")).
    Padding(0, 0, 0, 2)

  s.SelectedDesc = s.SelectedTitle.Copy().
    Foreground(lipgloss.Color("255")).
    Padding(0, 0, 0, 2)

  s.DimmedTitle = lipgloss.NewStyle().
    Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
    Padding(0, 0, 0, 2)

  s.DimmedDesc = s.DimmedTitle.Copy().
    Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

  s.FilterMatch = lipgloss.NewStyle().Underline(true)

  return s
}