package repository

import (
	"errors"
	"time"
)

var (
	ErrTodoListNotFound = errors.New("todo list not found")
	ErrTodoNotFound     = errors.New("todo not found")
	ErrEmptyTodoList    = errors.New("todo list is empty")
	ErrEmptyTitle       = errors.New("todo list title is empty")
	ErrEmptyDescription = errors.New("todo item description is empty")
)

type TodoList struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
	Todos []Todo `json:"todos"`
}

type Todo struct {
	ID          uint32    `json:"id"`
	ListID      uint32    `json:"list_id"`
	Description string    `json:"description"`
	Comments    string    `json:"comments"`
	DueDate     time.Time `json:"due_date"`
	Labels      []string  `json:"labels"`
	Done        bool      `json:"done"`
}

type Repository interface {
	InsertTodoList(todoList TodoList) (*TodoList, error)
	GetAllTodoLists() ([]TodoList, error)
	GetTodoListByID(id uint32) (*TodoList, error)
	UpdateTodoList(todoList TodoList) error
	DeleteTodoListByID(id uint32) error
	InsertTodo(todo Todo) (*Todo, error)
	GetTodoByID(id uint32) (*Todo, error)
	GetTodosByListID(listID uint32) ([]Todo, error)
	UpdateTodo(todo Todo) error
	DeleteTodo(todo Todo) error
}
