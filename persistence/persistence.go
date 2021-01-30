package persistence

import (
	"errors"
	"time"
)

var (
	ErrTodoListNotFound = errors.New("todo list not found")
)

type TodoList struct {
	ID    uint32
	Title string
}

type Todo struct {
	ID          uint32
	ListID      uint32
	Description string
	Comments    string
	DueDate     time.Time
	Labels      []string
	Done        bool
}

type Persistence interface {
	InsertTodoList(todoList *TodoList) error
	InsertTodo(todo *Todo) error
	GetTodoListByID(id uint32) (*TodoList, error)
	GetTodoByID(id uint32) (*Todo, error)
	GetTodosByListID(listID uint32) ([]Todo, error)
	UpdateTodoList(todoList *TodoList) (*TodoList, error)
	UpdateTodo(todo *Todo) (*Todo, error)
	DeleteTodoList(todoList *TodoList) error
	DeleteTodo(todo *Todo) error
}
