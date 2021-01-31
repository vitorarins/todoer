package persistence

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests for TodoList

func TestInsertTodoList(t *testing.T) {
	type Test struct {
		name                      string
		todoListTable             map[uint32]TodoList
		todoListAutoincrement     uint32
		wantTodoListAutoincrement uint32
		want                      *TodoList
		wantErr                   error
	}

	localStorage := NewLocalStorage()

	todoList := TodoList{
		Title: "Routine",
	}

	tests := []Test{
		{
			name:                      "SuccessInsertRoutine",
			todoListTable:             map[uint32]TodoList{},
			todoListAutoincrement:     0,
			wantTodoListAutoincrement: 1,
			want: &TodoList{
				ID:    0,
				Title: "Routine",
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable
			localStorage.TodoListAutoincrement = test.todoListAutoincrement

			got, err := localStorage.InsertTodoList(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoListAutoincrement, localStorage.TodoListAutoincrement); diff != "" {
				t.Errorf("InsertTodoList() mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("InsertTodoList() mismatch (-want +got):\n%s", diff)
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

func TestGetAllTodoList(t *testing.T) {
	type Test struct {
		name          string
		todoListTable map[uint32]TodoList
		want          []TodoList
		wantErr       error
	}

	localStorage := NewLocalStorage()

	tests := []Test{
		{
			name: "SuccessGetAllTodoListOneItem",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			want: []TodoList{
				TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			wantErr: nil,
		},
		{
			name: "SuccessGetAllTodoListTwoItems",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
				1: TodoList{
					ID:    1,
					Title: "Work",
				},
			},
			want: []TodoList{
				TodoList{
					ID:    0,
					Title: "Routine",
				},
				1: TodoList{
					ID:    1,
					Title: "Work",
				},
			},
			wantErr: nil,
		},
		{
			name:          "SuccessEmptyTodoListTable",
			todoListTable: map[uint32]TodoList{},
			want:          []TodoList{},
			wantErr:       ErrEmptyTodoListTable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable

			got, err := localStorage.GetAllTodoLists()

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetAllTodoLists() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateTodoList(t *testing.T) {
	type Test struct {
		name              string
		todoListTable     map[uint32]TodoList
		wantTodoListTable map[uint32]TodoList
		wantErr           error
	}

	localStorage := NewLocalStorage()

	todoList := TodoList{
		ID:    0,
		Title: "Routine",
	}

	tests := []Test{
		{
			name: "SuccessUpdateRoutine",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Rot",
				},
			},
			wantTodoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			wantErr: nil,
		},
		{
			name:              "ErrUpdateTodoListNotFound",
			todoListTable:     map[uint32]TodoList{},
			wantTodoListTable: map[uint32]TodoList{},
			wantErr:           ErrTodoListNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable

			err := localStorage.UpdateTodoList(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoListTable, localStorage.TodoListTable); diff != "" {
				t.Errorf("UpdateTodoList() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteTodoList(t *testing.T) {
	type Test struct {
		name              string
		todoListTable     map[uint32]TodoList
		wantTodoListTable map[uint32]TodoList
		wantErr           error
	}

	localStorage := NewLocalStorage()

	todoList := TodoList{
		ID:    0,
		Title: "Routine",
	}

	tests := []Test{
		{
			name: "SuccessDeleteRoutine",
			todoListTable: map[uint32]TodoList{
				0: TodoList{
					ID:    0,
					Title: "Routine",
				},
			},
			wantTodoListTable: map[uint32]TodoList{},
			wantErr:           nil,
		},
		{
			name:              "ErrDeleteTodoListNotFound",
			todoListTable:     map[uint32]TodoList{},
			wantTodoListTable: map[uint32]TodoList{},
			wantErr:           ErrTodoListNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable

			err := localStorage.DeleteTodoList(todoList)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoListTable, localStorage.TodoListTable); diff != "" {
				t.Errorf("DeleteTodoList() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Tests for Todo

func TestInsertTodo(t *testing.T) {

	type Test struct {
		name                  string
		todoListTable         map[uint32]TodoList
		todoAutoincrement     uint32
		wantTodoAutoincrement uint32
		want                  *Todo
		wantErr               error
	}

	localStorage := NewLocalStorage()

	todo := Todo{
		ListID:      0,
		Description: "Make the bed.",
		Comments:    "",
		Labels:      []string{"", ""},
	}

	tests := []Test{
		{
			name: "SuccessInsertMakeTheBed",
			todoListTable: map[uint32]TodoList{
				0: TodoList{},
			},
			todoAutoincrement:     0,
			wantTodoAutoincrement: 1,
			want: &Todo{
				ID:          0,
				ListID:      0,
				Description: "Make the bed.",
				Comments:    "",
				Labels:      []string{"", ""},
			},
			wantErr: nil,
		},
		{
			name:                  "ErrInsertTodoListNotFound",
			todoListTable:         map[uint32]TodoList{},
			todoAutoincrement:     0,
			wantTodoAutoincrement: 0,
			want:                  nil,
			wantErr:               ErrTodoListNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoListTable = test.todoListTable
			localStorage.TodoAutoincrement = test.todoAutoincrement

			got, err := localStorage.InsertTodo(todo)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoAutoincrement, localStorage.TodoAutoincrement); diff != "" {
				t.Errorf("InsertTodo() mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("InsertTodo() mismatch (-want +got):\n%s", diff)
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
			name: "ErrGetEmptyTodoList",
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
			want:                 []Todo{},
			wantErr:              ErrEmptyTodoList,
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
				t.Errorf("GetTodosByListID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	type Test struct {
		name          string
		todoTable     map[uint32]Todo
		wantTodoTable map[uint32]Todo
		wantErr       error
	}

	localStorage := NewLocalStorage()

	todo := Todo{
		ID:          0,
		ListID:      0,
		Description: "Make the bed.",
		Comments:    "It was hard",
		Labels:      []string{"bed", "bedroom"},
		Done:        true,
	}

	tests := []Test{
		{
			name: "SuccessUpdateMakeTheBed",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the",
					Comments:    "",
					Labels:      []string{"", ""},
				},
			},
			wantTodoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the bed.",
					Comments:    "It was hard",
					Labels:      []string{"bed", "bedroom"},
					Done:        true,
				},
			},
			wantErr: nil,
		},
		{
			name:          "ErrUpdateTodoNotFound",
			todoTable:     map[uint32]Todo{},
			wantTodoTable: map[uint32]Todo{},
			wantErr:       ErrTodoNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoTable = test.todoTable

			err := localStorage.UpdateTodo(todo)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoTable, localStorage.TodoTable); diff != "" {
				t.Errorf("UpdateTodo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	type Test struct {
		name                     string
		todoTable                map[uint32]Todo
		todoListRelationship     map[uint32][]uint32
		wantTodoTable            map[uint32]Todo
		wantTodoListRelationship map[uint32][]uint32
		wantErr                  error
	}

	localStorage := NewLocalStorage()

	todo := Todo{
		ID:          0,
		ListID:      0,
		Description: "Make the bed.",
		Comments:    "It was hard",
		Labels:      []string{"bed", "bedroom"},
		Done:        true,
	}

	tests := []Test{
		{
			name: "SuccessDeleteMakeTheBed",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the bed.",
					Comments:    "It was hard",
					Labels:      []string{"bed", "bedroom"},
					Done:        true,
				},
			},
			todoListRelationship: map[uint32][]uint32{
				0: []uint32{0},
			},
			wantTodoTable:            map[uint32]Todo{},
			wantTodoListRelationship: map[uint32][]uint32{},
			wantErr:                  nil,
		},
		{
			name:                     "ErrDeleteTodoNotFound",
			todoTable:                map[uint32]Todo{},
			todoListRelationship:     map[uint32][]uint32{},
			wantTodoTable:            map[uint32]Todo{},
			wantTodoListRelationship: map[uint32][]uint32{},
			wantErr:                  ErrTodoNotFound,
		},
		{
			name: "SuccessDeleteTodoListRelationshipEmpty",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the bed.",
					Comments:    "It was hard",
					Labels:      []string{"bed", "bedroom"},
					Done:        true,
				},
			},
			todoListRelationship:     map[uint32][]uint32{},
			wantTodoTable:            map[uint32]Todo{},
			wantTodoListRelationship: map[uint32][]uint32{},
			wantErr:                  nil,
		},
		{
			name: "SuccessDeleteTodoListRelationshipEmpty",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the bed.",
					Comments:    "It was hard",
					Labels:      []string{"bed", "bedroom"},
					Done:        true,
				},
			},
			todoListRelationship:     map[uint32][]uint32{},
			wantTodoTable:            map[uint32]Todo{},
			wantTodoListRelationship: map[uint32][]uint32{},
			wantErr:                  nil,
		},
		{
			name: "SuccessDeleteTodoListRelationshipBiggerThanOne",
			todoTable: map[uint32]Todo{
				0: Todo{
					ID:          0,
					ListID:      0,
					Description: "Make the bed.",
					Comments:    "It was hard",
					Labels:      []string{"bed", "bedroom"},
					Done:        true,
				},
				1: Todo{
					ID:          1,
					ListID:      0,
					Description: "Sweep the floor.",
					Comments:    "It was ok",
					Labels:      []string{"bedroom"},
					Done:        true,
				},
			},
			todoListRelationship: map[uint32][]uint32{
				0: []uint32{0, 1},
			},
			wantTodoTable: map[uint32]Todo{
				1: Todo{
					ID:          1,
					ListID:      0,
					Description: "Sweep the floor.",
					Comments:    "It was ok",
					Labels:      []string{"bedroom"},
					Done:        true,
				},
			},
			wantTodoListRelationship: map[uint32][]uint32{
				0: []uint32{1},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localStorage.TodoTable = test.todoTable
			localStorage.TodoListRelationship = test.todoListRelationship

			err := localStorage.DeleteTodo(todo)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}

			if diff := cmp.Diff(test.wantTodoTable, localStorage.TodoTable); diff != "" {
				t.Errorf("DeleteTodo() todo table mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.wantTodoListRelationship, localStorage.TodoListRelationship); diff != "" {
				t.Errorf("DeleteTodo() todo list relation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
