package persistence

import (
	"errors"
	"time"
)

var (
	ErrTodoListNotFound   = errors.New("todo list not found")
	ErrTodoNotFound       = errors.New("todo not found")
	ErrEmptyTodoList      = errors.New("todo list is empty")
	ErrEmptyTodoListTable = errors.New("todo list table is empty")
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
	InsertTodoList(todoList TodoList) (*TodoList, error)
	GetAllTodoLists() ([]TodoList, error)
	GetTodoListByID(id uint32) (*TodoList, error)
	UpdateTodoList(todoList TodoList) error
	DeleteTodoList(todoList TodoList) error
	InsertTodo(todo Todo) (*Todo, error)
	GetTodoByID(id uint32) (*Todo, error)
	GetTodosByListID(listID uint32) ([]Todo, error)
	UpdateTodo(todo Todo) error
	DeleteTodo(todo Todo) error
}
