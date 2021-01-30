package todolist

import (
	"errors"
	"testing"

	"github.com/vitorarins/todoer/persistence"
)

func TestCreate(t *testing.T) {
	type Test struct {
		name    string
		title   string
		wantErr error
	}

	tests := []Test{
		{
			name:    "SuccessRoutine",
			title:   "Routine",
			wantErr: nil,
		},
		{
			name:    "ErrEmptyTitle",
			title:   "",
			wantErr: ErrEmptyTitle,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todoList := persistence.TodoList{
				Title: test.title,
			}

			err := Create(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}
		})
	}
}
