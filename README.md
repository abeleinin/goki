# Goki

A terminal-based spaced repetition flashcard tool.

![Screenshot example of Goki](img/goki_main.png)

Theme: [material default-community](https://github.com/kaicataldo/material.vim)

Goki is an intelligent flashcard management tool inspired by 
[Anki](https://apps.ankiweb.net/) built in the terminal!

## Table of contents

- [TUI Demo](#tui-demo)
- [Key Mappings](#key-mappings)
- [Installation](#installation)
- [Examples](#examples)
- [Commands](#commands)
- [Resources](#resources)

## TUI Demo

Launch by running `goki`:

https://github.com/abeleinin/goki/assets/95333017/91d46c80-6580-4550-b6c2-405342a48290

Goki features a [Spaced Repetion Algorithm](https://en.wikipedia.org/wiki/Spaced_repetition)
which uses user feedback on card difficulty to effectively space out practice sessions.

## Key Mappings

<details>
<summary>Home Page</summary>

| Action       | Keybinding |
|--------------|------------|
| Review Decks Flashcards | `r`        |
| Create New Deck         | `N`        |
| View Deck Card List     | `o`        |
| Edit Deck Name          | `e`        |
| Delete Deck      | `d`        |
| Move Up          | `up arrow`,`k`      |
| Move Down        | `down arrow`,`j`      |
| Toggle Help Menu      | `?`        |
| Quit             | `q`,`ctrl+c` |

</details>

<details>
<summary>Flashcard List Page</summary>

| Action           | Keybinding |
|------------------|------------|
| Move Up          | `up arrow`,`k`      |
| Move Down        | `down arrow`,`j`      |
| Next page        | `right arrow`,`l`      |
| Previous Page    | `left arrow`,`h`      |
| Search Flashcards | `/`        |
| New Card     | `n`        |
| Edit Card    | `e`        |
| Delete Card  | `d`        |
| Undo Deleted Card | `u`    |

</details>

<details>
<summary>Create/Edit Flashcard Form</summary>

| Action         | Keybinding |
|----------------|------------|
| Next Field / Submit | `enter`    |
| Previous Field      | `tab`      |
| Exit Form           | `esc`      |

</details>

<details>
<summary>Flashcard</summary>

| Action      | Keybinding |
|-------------|------------|
| Exit Review | `esc`      |
| Show Back   | `o`        |
| Flashcard needs repeated again | `1`        |
| Flashcard took some thought | `2`        |
| Flashcard was easy to remember | `3`        |

</details>

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

## Examples

### Reviewing Flashcards

Press `r` on the selected deck you want to review on the home page. Or
use the command `goki review <deck index>` to review from the CLI.

**Review from CLI:**

![Review deck in cli](img/review_cli.gif)

**Review from TUI:**

![Review deck in tui](img/review.gif)

### Creating Decks

Press `N` in the home page. Use `e` to edit the currently selected deck.

![Create new deck](img/create_deck.gif)

### Creating Flashcards

Press `o` to view the cards in a deck. Press `n` to create a new card.

![Create new flashcard](img/create_card.gif)

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

