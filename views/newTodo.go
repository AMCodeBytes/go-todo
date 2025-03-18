package views

import (
	"fmt"
)

// newTodoView: This will return a string for the view
func NewTodoView(s string) string {
	return fmt.Sprintf(
		"What is the new item you wish to enter?\n\n%s\n\n%s",
		s,
		"(Press ctrl+c to quit)",
	) + "\n"
}
