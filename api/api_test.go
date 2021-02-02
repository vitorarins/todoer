package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vitorarins/todoer/repository"
)

// Todo List

func TestTodoListCreation(t *testing.T) {
	type Test struct {
		name           string
		requestBody    []byte
		method         string
		injectResponse repository.TodoList
		injectErr      error
		wantStatusCode int
		want           repository.TodoList
	}

	tests := []Test{
		{
			name:        "SuccessCreatingTodoList",
			requestBody: validTodoListRequestBody(t),
			injectResponse: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			want: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForDelete",
			method:         "DELETE",
			requestBody:    validTodoListRequestBody(t),
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestIfRepoInsertReturnsErrEmptyTitle",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      repository.ErrEmptyTitle,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRequestBodyIsEmpty",
			requestBody:    []byte{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRequestBodyIsNotValidJSON",
			requestBody:    []byte("{notvalidjson]"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "InternalServerErrorInsertTodoListError",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoList = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodPost
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath
			request := newRequest(t, method, testURL, test.requestBody)
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := repository.TodoList{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodoListGetAll(t *testing.T) {
	type Test struct {
		name           string
		method         string
		injectResponse []repository.TodoList
		injectErr      error
		wantStatusCode int
		want           []repository.TodoList
	}

	tests := []Test{
		{
			name: "SuccessRetrievingTodoLists",
			injectResponse: []repository.TodoList{
				repository.TodoList{
					ID:    1,
					Title: "Routine",
				},
			},
			want: []repository.TodoList{
				repository.TodoList{
					ID:    1,
					Title: "Routine",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForDelete",
			method:         "DELETE",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "InternalServerErrorGetTodoListsError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoListSlice = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodGet
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath
			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := []repository.TodoList{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodoListGetByID(t *testing.T) {
	type Test struct {
		name           string
		method         string
		idPath         string
		injectResponse repository.TodoList
		injectErr      error
		wantStatusCode int
		want           repository.TodoList
	}

	tests := []Test{
		{
			name: "SuccessRetrievingTodoLists",
			injectResponse: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			want: repository.TodoList{
				ID:    1,
				Title: "Routine",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForPost",
			method:         "POST",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestWrongPath",
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoUpdateReturnsErrTodoListNotFound",
			injectErr:      repository.ErrTodoListNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "InternalServerErrorGetTodoListsError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoList = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodGet
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/1"
			if test.idPath != "" {
				testURL = server.URL + TodoListPath + test.idPath
			}

			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := repository.TodoList{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodoListUpdate(t *testing.T) {
	type Test struct {
		name           string
		requestBody    []byte
		method         string
		idPath         string
		injectErr      error
		wantStatusCode int
	}

	tests := []Test{
		{
			name:           "SuccessUpdatingTodoList",
			requestBody:    validTodoListRequestBody(t),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForPost",
			method:         "POST",
			requestBody:    validTodoListRequestBody(t),
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestWrongPath",
			requestBody:    validTodoListRequestBody(t),
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoUpdateReturnsErrEmptyTitle",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      repository.ErrEmptyTitle,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoUpdateReturnsErrTodoListNotFound",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      repository.ErrTodoListNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "BadRequestIfRequestBodyIsEmpty",
			requestBody:    []byte{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRequestBodyIsNotValidJSON",
			requestBody:    []byte("{notvalidjson]"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "InternalServerErrorUpdateTodoListError",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodPut
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/1"
			if test.idPath != "" {
				testURL = server.URL + TodoListPath + test.idPath
			}

			request := newRequest(t, method, testURL, test.requestBody)
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}
		})
	}
}

func TestTodoListDelete(t *testing.T) {
	type Test struct {
		name           string
		method         string
		idPath         string
		injectErr      error
		wantStatusCode int
	}

	tests := []Test{
		{
			name:           "SuccessDeletingTodoList",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForPost",
			method:         "POST",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestWrongPath",
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoDeleteReturnsErrTodoListNotFound",
			injectErr:      repository.ErrTodoListNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "InternalServerErrorDeleteTodoListError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodDelete
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/1"
			if test.idPath != "" {
				testURL = server.URL + TodoListPath + test.idPath
			}

			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}
		})
	}
}

// Todo

func TestTodoCreation(t *testing.T) {
	type Test struct {
		name           string
		requestBody    []byte
		method         string
		listIDPath     string
		injectResponse repository.Todo
		injectErr      error
		wantStatusCode int
		want           repository.Todo
	}

	tests := []Test{
		{
			name:        "SuccessCreatingTodo",
			requestBody: validTodoRequestBody(t),
			injectResponse: repository.Todo{
				ID:          1,
				Description: "Make the bed",
			},
			want: repository.Todo{
				ID:          1,
				Description: "Make the bed",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "MethodNotAllowedForDelete",
			method:         "DELETE",
			requestBody:    validTodoRequestBody(t),
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestWrongPath",
			requestBody:    validTodoRequestBody(t),
			listIDPath:     "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoInsertReturnsErrEmptyDescription",
			requestBody:    validTodoRequestBody(t),
			injectErr:      repository.ErrEmptyDescription,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "NotFoundIfRepoInsertReturnsErrTodoListNotFound",
			requestBody:    validTodoRequestBody(t),
			injectErr:      repository.ErrTodoListNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "BadRequestIfRequestBodyIsEmpty",
			requestBody:    []byte{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRequestBodyIsNotValidJSON",
			requestBody:    []byte("{notvalidjson]"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "InternalServerErrorInsertTodoListError",
			requestBody:    validTodoRequestBody(t),
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodo = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodPost
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/0/todo"
			if test.listIDPath != "" {
				testURL = server.URL + TodoListPath + test.listIDPath + "/todo"
			}

			request := newRequest(t, method, testURL, test.requestBody)
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := repository.Todo{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodosGetByListID(t *testing.T) {
	type Test struct {
		name           string
		method         string
		listIDPath     string
		injectResponse []repository.Todo
		injectErr      error
		wantStatusCode int
		want           []repository.Todo
	}

	tests := []Test{
		{
			name: "SuccessRetrievingTodos",
			injectResponse: []repository.Todo{
				repository.Todo{
					ID:          1,
					ListID:      0,
					Description: "Make the bed",
				},
			},
			want: []repository.Todo{
				repository.Todo{
					ID:          1,
					ListID:      0,
					Description: "Make the bed",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "BadRequestWrongPath",
			listIDPath:     "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "MethodNotAllowedForDelete",
			method:         "DELETE",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "InternalServerErrorGetTodosError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodoSlice = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodGet
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/0/todo"
			if test.listIDPath != "" {
				testURL = server.URL + TodoListPath + test.listIDPath + "/todo"
			}

			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := []repository.Todo{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodoGetByID(t *testing.T) {
	type Test struct {
		name           string
		method         string
		idPath         string
		injectResponse repository.Todo
		injectErr      error
		wantStatusCode int
		want           repository.Todo
	}

	tests := []Test{
		{
			name: "SuccessRetrievingTodoByID",
			injectResponse: repository.Todo{
				ID:          1,
				ListID:      0,
				Description: "Make the bed",
			},
			want: repository.Todo{
				ID:          1,
				ListID:      0,
				Description: "Make the bed",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "BadRequestWrongIDPath",
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "MethodNotAllowedForPost",
			method:         "POST",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "NotFoundIfRepoGetTodoByIDReturnsTodoNotFound",
			injectErr:      repository.ErrTodoNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "InternalServerErrorGetTodosError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeTodo = test.injectResponse
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodGet
			if test.method != "" {
				method = test.method
			}

			testURL := server.URL + TodoListPath + "/0/todo/0"
			if test.idPath != "" {
				testURL = server.URL + TodoListPath + "/0/todo" + test.idPath
			}

			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			got := repository.Todo{}
			helperFromJSON(t, res.Body, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("api: POST %s mismatch (-want +got):\n%s", TodoListPath, diff)
			}

		})
	}
}

func TestTodoUpdate(t *testing.T) {
	type Test struct {
		name           string
		method         string
		listIDPath     string
		idPath         string
		requestBody    []byte
		injectErr      error
		wantStatusCode int
	}

	tests := []Test{
		{
			name:           "SuccessUpdatingTodo",
			requestBody:    validTodoRequestBody(t),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "BadRequestWrongIDPath",
			requestBody:    validTodoRequestBody(t),
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "MethodNotAllowedForPost",
			requestBody:    validTodoRequestBody(t),
			method:         "POST",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "BadRequestIfRepoInsertReturnsErrEmptyDescription",
			requestBody:    validTodoRequestBody(t),
			injectErr:      repository.ErrEmptyDescription,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRepoUpdateReturnsErrTodoNotFound",
			requestBody:    validTodoListRequestBody(t),
			injectErr:      repository.ErrTodoNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "BadRequestIfRequestBodyIsEmpty",
			requestBody:    []byte{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestIfRequestBodyIsNotValidJSON",
			requestBody:    []byte("{notvalidjson]"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "InternalServerErrorGetTodosError",
			requestBody:    validTodoRequestBody(t),
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodPut
			if test.method != "" {
				method = test.method
			}

			testBaseURL := server.URL + TodoListPath + "/0/todo"
			if test.listIDPath != "" {
				testBaseURL = server.URL + TodoListPath + test.listIDPath + "/todo"
			}

			testURL := testBaseURL + "/0"
			if test.idPath != "" {
				testURL = testBaseURL + test.idPath
			}

			request := newRequest(t, method, testURL, test.requestBody)
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}
		})
	}
}

func TestTodoDelete(t *testing.T) {
	type Test struct {
		name           string
		method         string
		listIDPath     string
		idPath         string
		injectErr      error
		wantStatusCode int
	}

	tests := []Test{
		{
			name:           "SuccessDeletingTodo",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "BadRequestWrongIDPath",
			idPath:         "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "BadRequestWrongListIDPath",
			listIDPath:     "/wrongpath",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "MethodNotAllowedForPost",
			method:         "POST",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "NotFoundIfRepoDeleteReturnsErrTodoListNotFound",
			injectErr:      repository.ErrTodoListNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "NotFoundIfRepoDeleteReturnsErrTodoNotFound",
			injectErr:      repository.ErrTodoNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "InternalServerErrorGetTodosError",
			injectErr:      errors.New("injected generic error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewFakeStorage()
			repo.FakeError = test.injectErr
			api := NewApi(repo)
			service := api.RegisterRoutes()
			server := httptest.NewServer(service)
			defer server.Close()

			method := http.MethodDelete
			if test.method != "" {
				method = test.method
			}

			testBaseURL := server.URL + TodoListPath + "/0/todo"
			if test.listIDPath != "" {
				testBaseURL = server.URL + TodoListPath + test.listIDPath + "/todo"
			}

			testURL := testBaseURL + "/0"
			if test.idPath != "" {
				testURL = testBaseURL + test.idPath
			}

			request := newRequest(t, method, testURL, []byte{})
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusOK {
				wantErr := ErrorResponse{}
				helperFromJSON(t, res.Body, &wantErr)

				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}
		})
	}
}

type FakeStorage struct {
	FakeTodoList      repository.TodoList
	FakeTodoListSlice []repository.TodoList
	FakeTodo          repository.Todo
	FakeTodoSlice     []repository.Todo
	FakeError         error
}

func NewFakeStorage() *FakeStorage {
	return &FakeStorage{
		FakeTodoList:      repository.TodoList{},
		FakeTodoListSlice: []repository.TodoList{},
		FakeTodo:          repository.Todo{},
		FakeTodoSlice:     []repository.Todo{},
		FakeError:         nil,
	}
}

func (fs *FakeStorage) InsertTodoList(todoList repository.TodoList) (*repository.TodoList, error) {
	return &fs.FakeTodoList, fs.FakeError
}
func (fs *FakeStorage) GetAllTodoLists() ([]repository.TodoList, error) {
	return fs.FakeTodoListSlice, fs.FakeError
}
func (fs *FakeStorage) GetTodoListByID(id uint32) (*repository.TodoList, error) {
	return &fs.FakeTodoList, fs.FakeError
}
func (fs *FakeStorage) UpdateTodoList(todoList repository.TodoList) error {
	return fs.FakeError
}
func (fs *FakeStorage) DeleteTodoListByID(id uint32) error {
	return fs.FakeError
}

func (fs *FakeStorage) InsertTodo(todo repository.Todo) (*repository.Todo, error) {
	return &fs.FakeTodo, fs.FakeError
}
func (fs *FakeStorage) GetTodoByID(id uint32) (*repository.Todo, error) {
	return &fs.FakeTodo, fs.FakeError
}
func (fs *FakeStorage) GetTodosByListID(listID uint32) ([]repository.Todo, error) {
	return fs.FakeTodoSlice, fs.FakeError
}
func (fs *FakeStorage) UpdateTodo(todo repository.Todo) error {
	return fs.FakeError
}
func (fs *FakeStorage) DeleteTodo(todo repository.Todo) error {
	return fs.FakeError
}

func helperFromJSON(t *testing.T, data io.Reader, v interface{}) {
	t.Helper()

	dec := json.NewDecoder(data)
	err := dec.Decode(&v)
	if err != nil {
		t.Fatal(err)
	}
}

func helperToJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	j, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func newRequest(t *testing.T, method string, url string, body []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func validTodoListRequestBody(t *testing.T) []byte {
	return helperToJSON(t, repository.TodoList{
		Title: "Routine",
	})
}

func validTodoRequestBody(t *testing.T) []byte {
	return helperToJSON(t, repository.Todo{
		Description: "Make the bed",
	})
}
