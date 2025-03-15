package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Todos struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Item      string `json:"item"`
	Completed int    `json:"completed"`
}

type model struct {
	choices  []Todo           // items on the list
	cursor   int              // which item our cursor is pointing at
	selected map[int]struct{} // which items are selected
}

func initialModel() model {
	file, err := os.Open("todo.json")

	if err != nil {
		// Failed to open the file, thus create a new blank file
		panic("Failed to open the file")
	}

	byteValue, err := io.ReadAll(file)

	if err != nil {
		// Issue reading the file
		panic("Failed to read the file")
	}

	var todos Todos

	json.Unmarshal(byteValue, &todos)

	return model{
		choices:  todos.Todos,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "w":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "s":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The View function is also where you could use Lip Gloss to style the view
	// Header
	s := "What do I need to do?\n\n"

	for i, choice := range m.choices {
		cursor := " "

		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Item)
	}

	// Footer
	s += "\n Press q to quit.\n"

	// Send UI to be rendered
	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}
