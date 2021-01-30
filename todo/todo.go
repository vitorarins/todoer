package todo

import (
	"errors"
	"time"
)

var (
	ErrEmptyDescription = errors.New("todo description is required")
)

type Todo struct {
	ID          uint32    `json:"id"`
	ListID      uint32    `json:"list_id"`
	Description string    `json:"description"`
	Comments    string    `json:"comments"`
	DueDate     time.Time `json:"due_date"`
	Labels      []string  `json:"labels"`
	Done        bool      `json:"done"`
}

func Create(todo Todo) error {
	if todo.Description == "" {
		return ErrEmptyDescription
	}
	return nil
}

func GetById(id uint32) (*Todo, error) {
	return nil, nil
}

func GetByList(listID uint32) ([]Todo, error) {
	return nil, nil
}

func (t *Todo) Update(newTodo *Todo) error {
	return nil
}

func (tl *Todo) Delete() error {
	return nil
}
