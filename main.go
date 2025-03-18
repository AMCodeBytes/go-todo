package main

import (
	"fmt"
	"go-todo/models"
	"go-todo/views"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// type Todos struct {
// 	Todos []Todo `json:"todos"`
// }

// type Todo struct {
// 	Item      string `json:"item"`
// 	Completed bool   `json:"completed"`
// }

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	NewItem   bool
	table     table.Model
	textInput textinput.Model
	choices   []models.Todo    // items on the list
	cursor    int              // which item our cursor is pointing at
	selected  map[int]struct{} // which items are selected
	help      bool             // display commands when this is true
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter todo item here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.Cursor.Blink = true

	return model{
		textInput: ti,
		choices:   models.GetTodos(),
		selected:  make(map[int]struct{}),
		NewItem:   false,
		help:      false,
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
			models.SaveTodo(m.choices)
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

			if m.help {
				m.table.Focus()
			} else {
				m.table.Blur()
			}
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
				var todo models.Todo
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

	if m.help {
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.NewItem {
		return views.NewTodoView(m.textInput.View())
	} else if m.help {
		columns := []table.Column{
			{Title: "Commands", Width: 10},
			{Title: "Description", Width: 35},
		}

		rows := []table.Row{
			{"ctrl+s", "Save the todo list"},
			{"ctrl+c", "Quit the app"},
			{"ctrl+n", "Toggle create new todo item input"},
			{"alt+h", "Toggle the help commands"},
			{"up", "Move up the list"},
			{"down", "Move down the list"},
			{"tab", "Complete a todo item"},
			{"delete", "Delete item from the todo list"},
			{"enter", "Submit the text input"},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		)

		style := table.DefaultStyles()

		style.Header = style.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)

		style.Selected = style.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(style)

		m := model{table: t}
		s := baseStyle.Render(m.table.View()) + "\n"
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
