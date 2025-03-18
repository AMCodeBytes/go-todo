package views

import (
	"go-todo/models"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func HelpView() string {
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

	m := models.Model{Table: t}
	return baseStyle.Render(m.Table.View()) + "\n"
}
