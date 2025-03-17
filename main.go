package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Todos struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

type model struct {
	NewItem   bool
	textInput textinput.Model
	choices   []Todo           // items on the list
	cursor    int              // which item our cursor is pointing at
	selected  map[int]struct{} // which items are selected
	help      bool             // display commands when this is true
}

func initialModel() model {
	file, err := os.OpenFile("todo.json", os.O_RDWR|os.O_CREATE, 0644)

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

	defer file.Close()

	ti := textinput.New()
	ti.Placeholder = "Enter todo item here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Blink = true

	return model{
		textInput: ti,
		choices:   todos.Todos,
		selected:  make(map[int]struct{}),
		NewItem:   false,
		help:      false,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func save(m model) {
	b, err := json.Marshal(Todos{m.choices})

	if err != nil {
		panic("Failed to marshal JSON data")
	}

	file, err := os.OpenFile("todo.json", os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		// Failed to open the file, thus create a new blank file
		panic("Failed to open the file")
	}

	err = file.Truncate(0)

	if err != nil {
		panic("Failed to remove contents from the file")
	}

	_, err = file.Write(b)

	if err != nil {
		panic("Failed to write to the file")
	}

	defer file.Close()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save the file
			save(m)
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "delete":
			m.choices = append(m.choices[:m.cursor], m.choices[m.cursor+1:]...)
		case "alt+h":
			m.help = !m.help
		case "ctrl+n":
			m.NewItem = !m.NewItem
		case "tab":
			// ok := m.choices[m.cursor].Completed
			if m.choices[m.cursor].Completed {
				m.choices[m.cursor].Completed = false
			} else {
				m.choices[m.cursor].Completed = true
			}
		case "enter":
			if m.NewItem {
				var todo Todo
				todo.Item = m.textInput.Value()
				todo.Completed = false

				m.choices = append(m.choices, todo)
				m.NewItem = !m.NewItem

				m.textInput.SetValue("")
			}
		}
	}

	if m.NewItem {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.NewItem {
		return fmt.Sprintf(
			"What is the new item you wish to enter?\n\n%s\n\n%s",
			m.textInput.View(),
			"(Press ctrl+c to quit)",
		) + "\n"
	} else if m.help {
		s := "What do I need to do? (alt+h to close help)\n\n"
		s += "ctrl+s | Save the todo list\n"
		s += "ctrl+c | Quit the app\n"
		s += "up | Move up the list\n"
		s += "down | Move down the list\n"
		s += "delete | Delete item from the todo list\n"
		s += "alt+h | Toggle the help commands\n"
		s += "ctrl+n | Toggle create new todo item input\n"
		s += "tab | Complete a todo item\n"
		s += "enter | Submit the text input\n"
		return s
	} else {
		// The View function is also where you could use Lip Gloss to style the view
		// Header
		s := "What do I need to do? (ctrl+h for help)\n\n"

		for i, choice := range m.choices {
			cursor := " "

			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if ok := choice.Completed; ok {
				checked = "x"
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Item)
		}

		// Footer
		s += "\n (Press ctrl+c to quit)\n"

		// Send UI to be rendered
		return s
	}
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}
