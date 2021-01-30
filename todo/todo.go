package todo

import (
	"errors"

	"github.com/vitorarins/todoer/persistence"
)

var (
	ErrEmptyDescription = errors.New("todo description is required")
)

func Create(todo persistence.Todo) error {
	if todo.Description == "" {
		return ErrEmptyDescription
	}
	return nil
}

func GetById(id uint32) (*persistence.Todo, error) {
	return nil, nil
}

func GetByList(todoList *persistence.TodoList) ([]persistence.Todo, error) {
	return nil, nil
}

func Update(newTodo *persistence.Todo) error {
	return nil
}

func Delete(todo *persistence.Todo) error {
	return nil
}
