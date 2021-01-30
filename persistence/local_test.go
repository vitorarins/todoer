package persistence

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Tests for TodoList

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

func TestGetTodoList(t *testing.T) {
	type Test struct {
		name          string
		todoListTable map[uint32]TodoList
		id            uint32
		want          *TodoList
		wantErr       error
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name: "SuccessGetRoutine",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			id: 0,
			want: &TodoList{
				ID:    0,
				Title: "Routine",
			},
			wantErr: nil,
		},
		{
			name:          "ErrGetTodoListNotFound",
			todoListTable: map[uint32]TodoList{},
			id:            0,
			want:          nil,
			wantErr:       ErrTodoListNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable

			got, err := localStorage.GetTodoListByID(test.id)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTodoListByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Tests for Todo

func TestInsertTodo(t *testing.T) {

	type Test struct {
		name          string
		todoListTable map[uint32]TodoList
		wantErr       error
	}

	localStorage := NewLocalStorage()

	todo := Todo{
		ListID:      0,
		Description: "Make the bed.",
		Comments:    "",
		DueDate:     time.Now(),
		Labels:      []string{"", ""},
	}

	tests := []Test{
		{
			name: "SuccessInsertMakeTheBed",
			todoListTable: map[uint32]TodoList{
				0: TodoList{},
			},
			wantErr: nil,
		},
		{
			name:          "ErrInsertTodoListNotFound",
			todoListTable: map[uint32]TodoList{},
			wantErr:       ErrTodoListNotFound,
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

func TestGetTodo(t *testing.T) {
	type Test struct {
		name      string
		todoTable map[uint32]Todo
		id        uint32
		want      *Todo
		wantErr   error
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name: "SuccessGetMakeTheBed",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					Description: "Make the bed",
				},
			},
			id: 0,
			want: &Todo{
				ID:          0,
				Description: "Make the bed",
			},
			wantErr: nil,
		},
		{
			name:      "ErrGetTodoNotFound",
			todoTable: map[uint32]Todo{},
			id:        0,
			want:      nil,
			wantErr:   ErrTodoNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoTable = test.todoTable

			got, err := localStorage.GetTodoByID(test.id)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTodoListByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetTodosByListID(t *testing.T) {
	type Test struct {
		name                 string
		todoListTable        map[uint32]TodoList
		todoTable            map[uint32]Todo
		todoListRelationship map[uint32][]uint32
		listID               uint32
		want                 []Todo
		wantErr              error
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name: "SuccessGetTodosFromRoutine",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					Description: "Make the bed",
				},
			},
			todoListRelationship: map[uint32][]uint32{
				0: []uint32{0},
			},
			listID: 0,
			want: []Todo{
				Todo{
					ID:          0,
					Description: "Make the bed",
				},
			},
			wantErr: nil,
		},
		{
			name: "SuccessGetEmptyTodosFromRoutine",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					Description: "Make the bed",
				},
			},
			todoListRelationship: map[uint32][]uint32{},
			listID:               0,
			want:                 nil,
			wantErr:              nil,
		},
		{
			name:                 "ErrGetTodoListNotFound",
			todoListTable:        map[uint32]TodoList{},
			todoTable:            map[uint32]Todo{},
			todoListRelationship: map[uint32][]uint32{},
			listID:               0,
			want:                 nil,
			wantErr:              ErrTodoListNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable
			localStorage.TodoTable = test.todoTable
			localStorage.TodoListRelationship = test.todoListRelationship

			got, err := localStorage.GetTodosByListID(test.listID)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTodoListByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
