package models

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	Table     table.Model      // the table model
	TextInput textinput.Model  // the text input model
	Choices   []Todo           // items on the list
	Cursor    int              // which item our cursor is pointing at
	Selected  map[int]struct{} // which items are selected
	NewItem   bool             // display the text input to create a new todo item
	Help      bool             // display commands when this is true
}
