package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const gokiLogo = `   ________        __    __
  /  _____/  ____ |  | _|__|
 /   \  ___ /    \|  |/ /  |
 \    \_\  |  /\  |    <|  |
  \______  /\____/|__|_ \__|
         \/            \/   `

var (
	screenWidth, screenHeight = GetScreenDimensions()

	// model.go
	logoStyle       = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	blurredStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	docStyle        = lipgloss.NewStyle().Width(100).Height(100).Align(lipgloss.Center)
	tableStyle      = lipgloss.NewStyle().Height(screenHeight - lipgloss.Height(gokiLogo) - 2)
	homeFooterStyle = lipgloss.NewStyle().Align(lipgloss.Left).Width(58)

	// deck.go
	cardStyle = lipgloss.NewStyle().MarginTop(screenHeight / 10).MarginLeft(3 * screenWidth / 10).Width(2 * screenWidth / 5).
			Height(screenHeight / 5).Border(lipgloss.RoundedBorder()).Align(lipgloss.Center)
	listStyle       = lipgloss.NewStyle().Align(lipgloss.Left).MarginLeft((screenWidth - 60) / 2).Padding(2)
	questionStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).MarginTop(2)
	answerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).MarginTop(4).MarginBottom(4)
	deckFooterStyle = lipgloss.NewStyle().MarginTop(10)
	progressStyle   = lipgloss.NewStyle().MarginTop(1).Render

	// form.go
	formTitleStyle  = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).MarginTop(1).Render
	formStyle       = lipgloss.NewStyle().Align(lipgloss.Center).Margin(screenHeight/10, screenWidth/4, 0, screenWidth/4).Padding(2, 2)
	viewStyle       = lipgloss.NewStyle().Width(2*screenWidth/5).Border(lipgloss.RoundedBorder()).Padding(0, 2)
	formFooterStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingTop(2)
	pad             = lipgloss.NewStyle().PaddingTop(1).Render

	// other
	helpKeyColor  = help.New().Styles.ShortKey.Inline(true)
	helpDescColor = help.New().Styles.ShortDesc.Inline(true)
)

func InitCustomDelegate() list.DefaultDelegate {
	delegate := list.DefaultDelegate{}
	delegate.ShowDescription = true
	delegate.Styles = CustomItemStyles()
	delegate.SetHeight(2)
	delegate.SetSpacing(1)
	return delegate
}

func CustomItemStyles() (s list.DefaultItemStyles) {
	s.NormalTitle = helpKeyColor.
		Padding(0, 0, 0, 1).Inline(false)

	s.NormalDesc = helpDescColor.
		Padding(0, 0, 0, 1).Inline(false)

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
