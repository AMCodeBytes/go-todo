package main

import (
	"fmt"
	"go-todo/models"
	"go-todo/views"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model models.Model

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter todo item here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Blink = true

	return model{
		TextInput: ti,
		Choices:   models.GetTodos(),
		Selected:  make(map[int]struct{}),
		NewItem:   false,
		Help:      false,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save the file
			models.SaveTodo(m.Choices)
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "delete":
			m.Choices = append(m.Choices[:m.Cursor], m.Choices[m.Cursor+1:]...)
		case "alt+h":
			m.Help = !m.Help

			if m.Help {
				m.Table.Focus()
			} else {
				m.Table.Blur()
			}
		case "ctrl+n":
			m.NewItem = !m.NewItem
		case "tab":
			// ok := m.choices[m.cursor].Completed
			if m.Choices[m.Cursor].Completed {
				m.Choices[m.Cursor].Completed = false
			} else {
				m.Choices[m.Cursor].Completed = true
			}
		case "enter":
			if m.NewItem {
				var todo models.Todo
				todo.Item = m.TextInput.Value()
				todo.Completed = false

				m.Choices = append(m.Choices, todo)
				m.NewItem = !m.NewItem

				m.TextInput.SetValue("")
			}
		}
	}

	if m.NewItem {
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	if m.Help {
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.NewItem {
		return views.NewTodoView(m.TextInput.View())
	} else if m.Help {
		return views.HelpView()
	} else {
		return views.TodoView(m.Choices, m.Cursor)
	}
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error: %v", err)
		os.Exit(1)
	}
}
