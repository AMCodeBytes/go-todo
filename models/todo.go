package models

import (
	"encoding/json"
	"io"
	"os"
)

type Todos struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

func GetTodos() []Todo {
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

	return todos.Todos
}

func SaveTodo(todos []Todo) {
	b, err := json.Marshal(Todos{todos})
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
