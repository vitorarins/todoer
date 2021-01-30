package todolist

import (
	"errors"

	"github.com/vitorarins/todoer/todo"
)

var (
	ErrEmptyTitle = errors.New("list title is required")
)

type TodoList struct {
	ID    uint32      `json:"id"`
	Title string      `json:"title"`
	Todos []todo.Todo `json:"todos"`
}

func Create(todoList TodoList) error {
	if todoList.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

func GetById(id uint32) (*TodoList, error) {
	return nil, nil
}

func (tl *TodoList) Update(newTodoList *TodoList) error {
	return nil
}

func (tl *TodoList) Delete() error {
	return nil
}
