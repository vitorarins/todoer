package persistence

import (
	"errors"
	"testing"
	"time"
)

func TestInsertTodoList(t *testing.T) {
	type Test struct {
		name    string
		title   string
		wantErr error
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name:    "SuccessInsertRoutine",
			title:   "Routine",
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoList := TodoList{
				Title: test.title,
			}

			err := localStorage.InsertTodoList(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}
		})
	}
}

func TestInsertTodo(t *testing.T) {

	type Test struct {
		name          string
		todoListTable map[uint32]TodoList
		wantErr       error
	}

	todo := Todo{
		ListID:      0,
		Description: "Make the bed.",
		Comments:    "",
		DueDate:     time.Now(),
		Labels:      []string{"", ""},
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name: "SuccessInsertMakeTheBed",
			todoListTable: map[uint32]TodoList{
				0: TodoList{},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable

			err := localStorage.InsertTodo(todo)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}
		})
	}
}
