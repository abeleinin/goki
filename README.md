# Goki

A terminal-based spaced repetition flashcard tool.

![Screenshot example of Goki](img/goki_main.png)

Theme: [material default-community](https://github.com/kaicataldo/material.vim)

Goki is an intelligent flashcard management tool inspired by 
[Anki](https://apps.ankiweb.net/) built in the terminal!

Goki features a [Spaced Repetion Algorithm](https://en.wikipedia.org/wiki/Spaced_repetition)
which uses user feedback on card difficulty efficienty plans for the next study time. 

## Table of contents

- [Installation](#installation)
- [Tutorial](#tutorial)
- [Commands](#commands)
- [Resources](#resources)

## Installation

Using `go`:

```
go install github.com/abeleinin/goki@latest
```

Build from source (go 1.13+)

```
git clone https://github.com/abeleinin/goki.git
cd goki
go build
```

## Tutorial

Refer the the `help` command or TUI help menu by pressing `?` for more info
on avaiable actions.

### Creating Decks

Press `N` in the home page. Use `e` to edit the currently selected deck.

![Create new deck](img/create_deck.gif)

### Creating Cards

Press `o` to view the cards in a deck. Press `n` to create a new card.

![Create new flashcard](img/create_card.gif)

### Reviewing Cards

Press `r` on the selected deck you want to review on the home page. Or
use the command `goki review <deck index>` to review from the CLI.

**Review from TUI:**

![Review deck in tui](img/review.gif)

**Review from CLI:**

![Review deck in cli](img/review_cli.gif)

## Commands

```
Usage: goki
  goki                        - tui mode
  goki list                   - view deck index
  goki review <deck index>    - review deck from cli`)
```

## Resources

- [Augmenting Long-term Memory ](https://augmentingcognition.com/ltm.html) by Michael Nielsen, Y Combinator Research, July 2018
- Created using [Charm](https://charm.sh/).

