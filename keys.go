package main

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
  return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {k.Up, k.Down, k.Left, k.Right}, // first column
    {k.Help, k.Quit},                // second column
  }
}

type keyMap struct {
  New    key.Binding
  Edit   key.Binding
  Delete key.Binding
  Up     key.Binding
  Down   key.Binding
  Right  key.Binding
  Left   key.Binding
  Enter  key.Binding
  Help   key.Binding
  Quit   key.Binding
  Back   key.Binding
  Tab    key.Binding
  Review key.Binding
  Open   key.Binding
  Easy   key.Binding
  Medium key.Binding
  Hard   key.Binding
}

var keys = keyMap{
  New: key.NewBinding(
    key.WithKeys("n"),
    key.WithHelp("n", "new"),
  ),
  Edit: key.NewBinding(
    key.WithKeys("e"),
    key.WithHelp("e", "edit"),
  ),
  Delete: key.NewBinding(
    key.WithKeys("d"),
    key.WithHelp("d", "delete"),
  ),
  Up: key.NewBinding(
    key.WithKeys("up", "k"),
    key.WithHelp("↑/k", "move up"),
  ),
  Down: key.NewBinding(
    key.WithKeys("down", "j"),
    key.WithHelp("↓/j", "move down"),
  ),
  Right: key.NewBinding(
    key.WithKeys("right", "l"),
    key.WithHelp("→/l", "move right"),
  ),
  Left: key.NewBinding(
    key.WithKeys("left", "h"),
    key.WithHelp("←/l", "move left"),
  ),
  Enter: key.NewBinding(
    key.WithKeys("enter"),
    key.WithHelp("enter", "enter"),
  ),
  Help: key.NewBinding(
    key.WithKeys("?"),
    key.WithHelp("?", "toggle help"),
  ),
  Quit: key.NewBinding(
    key.WithKeys("q", "ctrl+c"),
    key.WithHelp("q/ctrl+c", "quit"),
  ),
  Back: key.NewBinding(
    key.WithKeys("esc"),
    key.WithHelp("esc", "back"),
  ),
  Tab: key.NewBinding(
    key.WithKeys("tab"),
    key.WithHelp("tab", "move focus"),
  ),
  Review: key.NewBinding(
    key.WithKeys("r"),
    key.WithHelp("r", "Review selected deck"),
  ),
  Open: key.NewBinding(
    key.WithKeys("o"),
    key.WithHelp("o", "Reveal card"),
  ),
  Easy: key.NewBinding(
    key.WithKeys("1"),
    key.WithHelp("1", "Card easy"),
  ),
  Medium: key.NewBinding(
    key.WithKeys("2"),
    key.WithHelp("2", "Card medium"),
  ),
  Hard: key.NewBinding(
    key.WithKeys("3"),
    key.WithHelp("3", "Card hard"),
  ),
}