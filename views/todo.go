package views

import (
	"fmt"
	"go-todo/models"
)

func TodoView(choices []models.Todo, pos int) string {
	// Header
	s := "What do I need to do? (alt+h for help)\n\n"

	for i, choice := range choices {
		cursor := " "

		if pos == i {
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
