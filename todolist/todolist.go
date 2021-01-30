package todolist

import (
	"errors"

	"github.com/vitorarins/todoer/persistence"
)

var (
	ErrEmptyTitle = errors.New("list title is required")
)

func Create(todoList persistence.TodoList) error {
	if todoList.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

func GetById(id uint32) (*persistence.TodoList, error) {
	return nil, nil
}

func Update(newTodoList *persistence.TodoList) error {
	return nil
}

func Delete(todoList *persistence.TodoList) error {
	return nil
}
