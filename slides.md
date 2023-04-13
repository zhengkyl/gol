---
author: puhack.horse/gol-gh
date: ""
---

# Game of Life Terminal App

## Install Go

https://go.dev/doc/install

### OR

## Use Github Codespaces

https://puhack.horse/gol-gh

---

# Demo

## Just run it locally

## Hosted multiplayer version

`ssh gool.fly.dev`

---

# Libaries by [Charm](https://charm.sh)

## [bubbletea](https://github.com/charmbracelet/bubbletea)

terminal app framework

## [lipgloss](https://github.com/charmbracelet/lipgloss)

colors and styling for the terminal

---

# Go variables

```go
var five int
five = 4

var four = 5
const val = true // rune, string, bool, number

// var sentence = "i like 冰淇淋😳"
sentence := "i like 冰淇淋😳"
```

---

# Go functions

```go
func add(a int, b int) int {
  return a + b
}

func swap(c, d string) (string, string) {
  return d, c
}

```

---

# Go loops

```go
slice := []int{ 1, 2, 3 }

for index, value := range slice {
}

for i := 0; i < len(slice); i++ {
}
```

---

# Learn Go

This workshop is NOT representative of Go as a language.

Try the official Go tutorial. It's actually good!

https://go.dev/tour

---

# Setup new project

Initialize module `go mod init MODULE_NAME`

Create `main` package

---

# life.go

Handle logic for Conway's Game of Life

- `NewBoard()`

- `NextBoard()`

---

# model

holds all "information" (state) needed to render the ui

must implement tea.Model

- `Init()`
- `Update()`
- `View()`

---

# Dependencies

VSCode Go Extension should do automagically

First, try `go mod tidy`

## Otherwise

`go get package "github.com/charmbracelet/bubbletea"`

`go get package "github.com/charmbracelet/lipgloss"`

---

# Init

returns a commmand to run before the first Update()

## `tea.Cmd`

function that returns a `tea.Msg` "later"

good for io, http requests, etc.

---

# Update

receives a message which represents "something happening"

returns

- the updated model
- any commands we want to run

---

# View

returns a string. this IS the ui

---

# main.go

```go
p := tea.NewProgram(ui.New(50, 25), tea.WithAltScreen())
_, err := p.Run()
if err != nil {
  fmt.Printf("L + R, fix your code: %v", err)
  os.Exit(1)
}
```

---

# Customization✨

## Different game rules

- look at update logic in `NextBoard()`

## Change visuals

- change [][]bool to something else?
- use different lipgloss colors
