package main

import "github.com/charmbracelet/bubbles/key"

func (u User) ShortHelp() []key.Binding {
  return []key.Binding{u.KeyMap.Help, u.KeyMap.Quit}
}

func (u User) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {u.KeyMap.Help, u.KeyMap.Up, u.KeyMap.Down, u.KeyMap.Quit},
    {u.KeyMap.New, u.KeyMap.Open, u.KeyMap.Edit, u.KeyMap.Review},
  }
}

func (f Form) ShortHelp() []key.Binding {
  return []key.Binding{f.keyMap.Enter, f.keyMap.Tab, f.keyMap.Back}
}

func (f Form) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {f.keyMap.Enter, f.keyMap.Tab, f.keyMap.Back},
  }
}

func (d Deck) ShortHelp() []key.Binding {
  return []key.Binding{d.keyMap.Again, d.keyMap.Good, d.keyMap.Easy}
}

func (d Deck) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {d.keyMap.New, d.keyMap.Edit, d.keyMap.Delete, d.keyMap.Up},
    {d.keyMap.Down, d.keyMap.Enter, d.keyMap.Help, d.keyMap.Quit}, 
    {d.keyMap.Back, d.keyMap.Tab, d.keyMap.Review, d.keyMap.Open},
  }
}

func (c Card) ShortHelp() []key.Binding {
  return []key.Binding{cardKeyMap.Open, cardKeyMap.Back}
}

func (c Card) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {cardKeyMap.Open, cardKeyMap.Back},
  }
}

type keyMap struct {
  New    key.Binding
  Edit   key.Binding
  Delete key.Binding
  Up     key.Binding
  Down   key.Binding
  Enter  key.Binding
  Help   key.Binding
  Quit   key.Binding
  Back   key.Binding
  Tab    key.Binding
  Review key.Binding
  Open   key.Binding
  Easy   key.Binding
  Good   key.Binding
  Again  key.Binding
  Save   key.Binding
  Search key.Binding

  ShowFullHelp  key.Binding
  CloseFullHelp key.Binding
}

var cardKeyMap = keyMap{
  Back: key.NewBinding(
    key.WithKeys("esc"),
    key.WithHelp("esc", "previous page"),
  ),
  Open: key.NewBinding(
    key.WithKeys("o"),
    key.WithHelp("o", "show back"),
  ),
}

func DeckKeyMap() keyMap {
	return keyMap{
    New: key.NewBinding(
      key.WithKeys("n"),
      key.WithHelp("n", "new card"),
    ),
    Edit: key.NewBinding(
      key.WithKeys("e"),
      key.WithHelp("e", "edit card"),
    ),
    Back: key.NewBinding(
      key.WithKeys("esc"),
      key.WithHelp("esc", "home page"),
    ),
    Open: key.NewBinding(
      key.WithKeys("o"),
      key.WithHelp("o", "hide/show description"),
    ),
    Again: key.NewBinding(
      key.WithKeys("1"),
      key.WithHelp("1", "again"),
    ),
    Good: key.NewBinding(
      key.WithKeys("2"),
      key.WithHelp("2", "good"),
    ),
    Easy: key.NewBinding(
      key.WithKeys("3"),
      key.WithHelp("3", "easy"),
    ),
    Save: key.NewBinding(
      key.WithKeys("ctrl+s"),
      key.WithHelp("ctrl+s", "save cards"),
    ),
    Delete: key.NewBinding(
      key.WithKeys("d"),
      key.WithHelp("d", "delete card"),
    ),
    Search: key.NewBinding(
      key.WithKeys("/"),
      key.WithHelp("/", "filter cards"),
    ),
    Enter: key.NewBinding(
      key.WithKeys("enter"),
      key.WithHelp("enter", "search"),
    ),
  }
}

func FormKeyMap() keyMap {
  return keyMap{
    Enter: key.NewBinding(
      key.WithKeys("enter"),
      key.WithHelp("enter", "next field/submit"),
    ),
    Tab: key.NewBinding(
      key.WithKeys("tab"),
      key.WithHelp("tab", "previous field"),
    ),
    Back: key.NewBinding(
      key.WithKeys("esc"),
      key.WithHelp("esc", "previous page"),
    ),
  }
}

func DefaultKeyMap() keyMap {
	return keyMap{
    New: key.NewBinding(
      key.WithKeys("N"),
      key.WithHelp("N", "new deck"),
    ),
    Edit: key.NewBinding(
      key.WithKeys("e"),
      key.WithHelp("e", "edit deck name"),
    ),
    Delete: key.NewBinding(
      key.WithKeys("d"),
      key.WithHelp("d", "delete deck"),
    ),
    Up: key.NewBinding(
      key.WithKeys("up", "k"),
      key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
      key.WithKeys("down", "j"),
      key.WithHelp("↓/j", "move down"),
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
      key.WithHelp("esc", "back to previous page"),
    ),
    Tab: key.NewBinding(
      key.WithKeys("tab"),
      key.WithHelp("tab", "move focus"),
    ),
    Review: key.NewBinding(
      key.WithKeys("r"),
      key.WithHelp("r", "review deck"),
    ),
    Open: key.NewBinding(
      key.WithKeys("o"),
      key.WithHelp("o", "view deck"),
    ),
    ShowFullHelp: key.NewBinding(
      key.WithKeys("?"),
      key.WithHelp("?", "more"),
    ),
    CloseFullHelp: key.NewBinding(
      key.WithKeys("?"),
      key.WithHelp("?", "close help"),
    ),
  }
}