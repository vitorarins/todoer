package api

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vitorarins/todoer/pb"
	"github.com/vitorarins/todoer/repository"
)

var ctx = context.Background()

// TodoList

func TestTodoListGrpcApiCreation(t *testing.T) {
	type Test struct {
		name              string
		createTodoListReq *pb.CreateTodoListRequest
		injectResponse    repository.TodoList
		injectErr         error
		want              *pb.CreateTodoListReply
	}

	tests := []Test{
		{
			name:              "SuccessCreatingTodoList",
			createTodoListReq: validCreateTodoListRequest(t),
			injectResponse: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			want: &pb.CreateTodoListReply{
				TodoList: &pb.TodoList{
					Id:    1,
					Title: "Routine",
				},
			},
		},
		{
			name:              "BadRequestIfRepoInsertReturnsErrEmptyTitle",
			createTodoListReq: validCreateTodoListRequest(t),
			injectErr:         repository.ErrEmptyTitle,
		},
		{
			name:              "InternalServerErrorInsertTodoListError",
			createTodoListReq: validCreateTodoListRequest(t),
			injectErr:         errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoList = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.CreateTodoList(ctx, test.createTodoListReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.CreateTodoListReply{}),
				cmpopts.IgnoreUnexported(pb.TodoList{})); diff != "" {

				t.Errorf("grpc_api: CreateTodoList mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoListGrpcApiGetAll(t *testing.T) {
	type Test struct {
		name           string
		injectResponse []repository.TodoList
		injectErr      error
		want           *pb.GetAllTodoListsReply
	}

	tests := []Test{
		{
			name: "SuccessGetAllTodoLists",
			injectResponse: []repository.TodoList{
				repository.TodoList{
					ID:    1,
					Title: "Work",
				},
			},
			want: &pb.GetAllTodoListsReply{
				TodoLists: []*pb.TodoList{
					&pb.TodoList{
						Id:    1,
						Title: "Work",
					},
				},
			},
		},
		{
			name:      "InternalServerErrorInsertTodoListError",
			injectErr: errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoListSlice = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.GetAllTodoLists(ctx, &pb.Empty{})
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.GetAllTodoListsReply{}),
				cmpopts.IgnoreUnexported(pb.TodoList{})); diff != "" {

				t.Errorf("grpc_api: GetAllTodoLists mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoListGrpcApiGetByID(t *testing.T) {
	type Test struct {
		name           string
		getTodoListReq *pb.GetTodoListRequest
		injectResponse repository.TodoList
		injectErr      error
		want           *pb.GetTodoListReply
	}

	tests := []Test{
		{
			name:           "SuccessGetTodoList",
			getTodoListReq: validGetTodoListRequest(t),
			injectResponse: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			want: &pb.GetTodoListReply{
				TodoList: &pb.TodoList{
					Id:    1,
					Title: "Routine",
				},
			},
		},
		{
			name:           "ErrIfTodoListNotFound",
			getTodoListReq: validGetTodoListRequest(t),
			injectErr:      repository.ErrTodoListNotFound,
		},
		{
			name:           "InternalServerErrorInsertTodoListError",
			getTodoListReq: validGetTodoListRequest(t),
			injectErr:      errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoList = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.GetTodoList(ctx, test.getTodoListReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.GetTodoListReply{}),
				cmpopts.IgnoreUnexported(pb.TodoList{})); diff != "" {

				t.Errorf("grpc_api: GetTodoList mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoListGrpcApiUpdate(t *testing.T) {
	type Test struct {
		name              string
		updateTodoListReq *pb.UpdateTodoListRequest
		injectResponse    repository.TodoList
		injectErr         error
	}

	tests := []Test{
		{
			name:              "SuccessUpdatingTodoList",
			updateTodoListReq: validUpdateTodoListRequest(t),
			injectResponse: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
		},
		{
			name:              "BadRequestIfRepoUpdateReturnsErrEmptyTitle",
			updateTodoListReq: validUpdateTodoListRequest(t),
			injectErr:         repository.ErrEmptyTitle,
		},
		{
			name:              "ErrIfTodoListNotFound",
			updateTodoListReq: validUpdateTodoListRequest(t),
			injectErr:         repository.ErrTodoListNotFound,
		},
		{
			name:              "InternalServerErrorInsertTodoListError",
			updateTodoListReq: validUpdateTodoListRequest(t),
			injectErr:         errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoList = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			_, err := grpcApi.UpdateTodoList(ctx, test.updateTodoListReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTodoListGrpcApiDelete(t *testing.T) {
	type Test struct {
		name              string
		deleteTodoListReq *pb.DeleteTodoListRequest
		injectErr         error
	}

	tests := []Test{
		{
			name:              "SuccessDeletingTodoList",
			deleteTodoListReq: validDeleteTodoListRequest(t),
		},
		{
			name:              "ErrIfTodoListNotFound",
			deleteTodoListReq: validDeleteTodoListRequest(t),
			injectErr:         repository.ErrTodoListNotFound,
		},
		{
			name:              "InternalServerErrorInsertTodoListError",
			deleteTodoListReq: validDeleteTodoListRequest(t),
			injectErr:         errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			_, err := grpcApi.DeleteTodoList(ctx, test.deleteTodoListReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}
		})
	}
}

func validCreateTodoListRequest(t *testing.T) *pb.CreateTodoListRequest {
	return &pb.CreateTodoListRequest{
		Title: "Routine",
	}
}

func validGetTodoListRequest(t *testing.T) *pb.GetTodoListRequest {
	return &pb.GetTodoListRequest{
		Id: 0,
	}
}

func validUpdateTodoListRequest(t *testing.T) *pb.UpdateTodoListRequest {
	return &pb.UpdateTodoListRequest{
		TodoList: &pb.TodoList{
			Id:    0,
			Title: "Routine",
		},
	}
}

func validDeleteTodoListRequest(t *testing.T) *pb.DeleteTodoListRequest {
	return &pb.DeleteTodoListRequest{
		Id: 0,
	}
}

// Todo

func TestTodoGrpcApiCreation(t *testing.T) {
	type Test struct {
		name           string
		createTodoReq  *pb.CreateTodoRequest
		injectResponse repository.Todo
		injectErr      error
		want           *pb.CreateTodoReply
		wantErr        bool
	}

	tests := []Test{
		{
			name:          "SuccessCreatingTodo",
			createTodoReq: validCreateTodoRequest(t),
			injectResponse: repository.Todo{
				ID:          1,
				Description: "Routine",
				DueDate:     parseTime(t, "2021-01-01T00:00:00Z"),
			},
			want: &pb.CreateTodoReply{
				Todo: &pb.Todo{
					Id:          1,
					Description: "Routine",
					DueDate:     "2021-01-01T00:00:00Z",
				},
			},
		},
		{
			name: "ErrCreatingTodoInvalidDueDate",
			createTodoReq: &pb.CreateTodoRequest{
				Description: "Routine",
				DueDate:     "00010101",
			},
			injectResponse: repository.Todo{
				ID:          1,
				Description: "Routine",
			},
			wantErr: true,
		},
		{
			name:          "BadRequestIfRepoInsertReturnsErrEmptyDescription",
			createTodoReq: validCreateTodoRequest(t),
			injectErr:     repository.ErrEmptyDescription,
			wantErr:       true,
		},
		{
			name:          "InternalServerErrorInsertTodoError",
			createTodoReq: validCreateTodoRequest(t),
			injectErr:     errors.New("injected generic error"),
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodo = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.CreateTodo(ctx, test.createTodoReq)
			if err != nil && !test.wantErr {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.CreateTodoReply{}),
				cmpopts.IgnoreUnexported(pb.Todo{})); diff != "" {

				t.Errorf("grpc_api: CreateTodo mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoGrpcApiGetTodosByList(t *testing.T) {
	type Test struct {
		name              string
		getTodosByListReq *pb.GetTodosByListRequest
		injectResponse    []repository.Todo
		injectErr         error
		want              *pb.GetTodosByListReply
	}

	tests := []Test{
		{
			name:              "SuccessGetTodosByList",
			getTodosByListReq: validGetTodosByListRequest(t),
			injectResponse: []repository.Todo{
				repository.Todo{
					ID:          1,
					Description: "Make the bed",
				},
			},
			want: &pb.GetTodosByListReply{
				Todos: []*pb.Todo{
					&pb.Todo{
						Id:          1,
						Description: "Make the bed",
					},
				},
			},
		},
		{
			name:              "ErrIfTodoListNotFound",
			getTodosByListReq: validGetTodosByListRequest(t),
			injectErr:         repository.ErrTodoListNotFound,
		},
		{
			name:              "InternalServerErrorInsertTodoError",
			getTodosByListReq: validGetTodosByListRequest(t),
			injectErr:         errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoSlice = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.GetTodosByList(ctx, test.getTodosByListReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.GetTodosByListReply{}),
				cmpopts.IgnoreUnexported(pb.Todo{})); diff != "" {

				t.Errorf("grpc_api: GetTodosByList mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoGrpcApiGetByID(t *testing.T) {
	type Test struct {
		name           string
		getTodoReq     *pb.GetTodoRequest
		injectResponse repository.Todo
		injectErr      error
		want           *pb.GetTodoReply
	}

	tests := []Test{
		{
			name:       "SuccessGetTodo",
			getTodoReq: validGetTodoRequest(t),
			injectResponse: repository.Todo{
				ID:          1,
				Description: "Make the bed",
			},
			want: &pb.GetTodoReply{
				Todo: &pb.Todo{
					Id:          1,
					Description: "Make the bed",
				},
			},
		},
		{
			name:       "InternalServerErrorInsertTodoError",
			getTodoReq: validGetTodoRequest(t),
			injectErr:  errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodo = test.injectResponse
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			got, err := grpcApi.GetTodo(ctx, test.getTodoReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got,
				cmpopts.IgnoreUnexported(pb.GetTodoReply{}),
				cmpopts.IgnoreUnexported(pb.Todo{})); diff != "" {

				t.Errorf("grpc_api: GetTodo mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestTodoGrpcApiUpdate(t *testing.T) {
	type Test struct {
		name          string
		updateTodoReq *pb.UpdateTodoRequest
		injectErr     error
		wantErr       bool
	}

	tests := []Test{
		{
			name:          "SuccessUpdatingTodo",
			updateTodoReq: validUpdateTodoRequest(t),
		},
		{
			name: "ErrUpdatingTodoInvalidDueDate",
			updateTodoReq: &pb.UpdateTodoRequest{
				Todo: &pb.Todo{
					Description: "Routine",
					DueDate:     "00010101",
				},
			},
			wantErr: true,
		},
		{
			name:          "BadRequestIfRepoUpdateReturnsErrEmptyDescription",
			updateTodoReq: validUpdateTodoRequest(t),
			injectErr:     repository.ErrEmptyDescription,
			wantErr:       true,
		},
		{
			name:          "InternalServerErrorInsertTodoError",
			updateTodoReq: validUpdateTodoRequest(t),
			injectErr:     errors.New("injected generic error"),
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			_, err := grpcApi.UpdateTodo(ctx, test.updateTodoReq)
			if err != nil && !test.wantErr {
				t.Fatal(err)
			}
		})
	}
}

func TestTodoGrpcApiDelete(t *testing.T) {
	type Test struct {
		name          string
		deleteTodoReq *pb.DeleteTodoRequest
		injectErr     error
	}

	tests := []Test{
		{
			name:          "SuccessDeletingTodo",
			deleteTodoReq: validDeleteTodoRequest(t),
		},
		{
			name:          "ErrIfTodoNotFound",
			deleteTodoReq: validDeleteTodoRequest(t),
			injectErr:     repository.ErrTodoNotFound,
		},
		{
			name:          "InternalServerErrorInsertTodoError",
			deleteTodoReq: validDeleteTodoRequest(t),
			injectErr:     errors.New("injected generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			grpcApi := NewGrpcApi(repo)

			_, err := grpcApi.DeleteTodo(ctx, test.deleteTodoReq)
			if err != nil && test.injectErr == nil {
				t.Fatal(err)
			}
		})
	}
}

func validCreateTodoRequest(t *testing.T) *pb.CreateTodoRequest {
	return &pb.CreateTodoRequest{
		Description: "Routine",
		DueDate:     "2021-01-01T00:00:00Z",
	}
}

func validGetTodoRequest(t *testing.T) *pb.GetTodoRequest {
	return &pb.GetTodoRequest{
		Id: 0,
	}
}

func validGetTodosByListRequest(t *testing.T) *pb.GetTodosByListRequest {
	return &pb.GetTodosByListRequest{
		ListId: 0,
	}
}

func validUpdateTodoRequest(t *testing.T) *pb.UpdateTodoRequest {
	return &pb.UpdateTodoRequest{
		Todo: &pb.Todo{
			Id:          0,
			Description: "Routine",
			DueDate:     "2021-01-01T00:00:00Z",
		},
	}
}

func validDeleteTodoRequest(t *testing.T) *pb.DeleteTodoRequest {
	return &pb.DeleteTodoRequest{
		Id: 0,
	}
}
