package todolist

import (
	"errors"
	"testing"
	"time"

	"github.com/vitorarins/todoer/todo"
)

func TestCreate(t *testing.T) {
	type Test struct {
		name    string
		title   string
		todos   []todo.Todo
		wantErr error
	}

	tests := []Test{
		{
			name:    "SuccessRoutine",
			title:   "Routine",
			todos:   []todo.Todo{},
			wantErr: nil,
		},
		{
			name:    "ErrEmptyTitle",
			title:   "",
			todos:   []todo.Todo{},
			wantErr: ErrEmptyTitle,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoList := TodoList{
				Title: test.title,
				Todos: test.todos,
			}

			err := Create(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}
		})
	}
}

func parseTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
