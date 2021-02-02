package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/vitorarins/todoer/repository"
)

const (
	dateLayout     = time.RFC3339
	TodoListPath   = "/todolist"
	TodoListIDPath = TodoListPath + "/{id}"
	TodoPath       = TodoListPath + "/{list_id}/todo"
	TodoIDPath     = TodoPath + "/{id}"
)

var (
	ErrInvalidDueDate = errors.New("due_date is invalid")
)

type TodoListTransport struct {
	ID    uint32 `json:"id"`
	Title string `json:"title"`
}

type TodoTransport struct {
	ID          uint32   `json:"id"`
	ListID      uint32   `json:"list_id"`
	Description string   `json:"description"`
	Comments    string   `json:"comments"`
	DueDate     string   `json:"due_date"`
	Labels      []string `json:"labels"`
	Done        bool     `json:"done"`
}

type Error struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Api struct {
	repo repository.Repository
}

func NewApi(repo repository.Repository) Api {
	return Api{
		repo: repo,
	}
}

func (a *Api) RegisterRoutes() http.Handler {
	handler := mux.NewRouter()
	handler.HandleFunc(TodoListPath, a.TodoList)
	handler.HandleFunc(TodoListIDPath, a.TodoListByID)
	handler.HandleFunc(TodoPath, a.Todo)
	handler.HandleFunc(TodoIDPath, a.TodoByID)
	return handler
}

// Todo List

func (a *Api) TodoList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"path": TodoListPath})

	switch req.Method {
	case http.MethodPost:
		a.CreateTodoList(res, req)
	case http.MethodGet:
		a.GetAllTodoLists(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("method %q is not allowed", req.Method)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("method not allowed")
		return
	}

}

func (a *Api) CreateTodoList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "CreateTodoList"})

	dec := json.NewDecoder(req.Body)
	todoListReq := TodoListTransport{}

	err := dec.Decode(&todoListReq)
	if err != nil {
		msg := fmt.Sprintf("cant parse request body as JSON:%v", err)
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("invalid request body")
		return
	}

	newTodoList, err := a.repo.InsertTodoList(fromTransportToTodoList(todoListReq))
	if err != nil {
		if errors.Is(err, repository.ErrEmptyTitle) {
			res.WriteHeader(http.StatusBadRequest)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("bad request error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todoListRes := toTransportTodoList(*newTodoList)
	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todoListRes))
}

func (a *Api) GetAllTodoLists(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "GetAllTodoLists"})

	todoLists, err := a.repo.GetAllTodoLists()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todoListsRes := []TodoListTransport{}
	for _, tl := range todoLists {
		todoListsRes = append(todoListsRes, toTransportTodoList(tl))
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todoListsRes))
}

func (a *Api) TodoListByID(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"path": TodoListIDPath})

	switch req.Method {
	case http.MethodGet:
		a.GetTodoList(res, req)
	case http.MethodPut:
		a.UpdateTodoList(res, req)
	case http.MethodDelete:
		a.DeleteTodoList(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("method %q is not allowed", req.Method)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("method not allowed")
		return
	}

}

func (a *Api) GetTodoList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "GetTodoList"})

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}

	todoList, err := a.repo.GetTodoListByID(uint32(id))
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todoListRes := toTransportTodoList(*todoList)
	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todoListRes))
}

func (a *Api) UpdateTodoList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "UpdateTodoList"})

	dec := json.NewDecoder(req.Body)
	todoListReq := TodoListTransport{}

	err := dec.Decode(&todoListReq)
	if err != nil {
		msg := fmt.Sprintf("cant parse request body as JSON:%v", err)
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("invalid request body")
		return
	}

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}
	todoListReq.ID = uint32(id)

	err = a.repo.UpdateTodoList(fromTransportToTodoList(todoListReq))
	if err != nil {
		if errors.Is(err, repository.ErrEmptyTitle) {
			res.WriteHeader(http.StatusBadRequest)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("bad request error")
			return
		}
		if errors.Is(err, repository.ErrTodoListNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, []byte{})
}

func (a *Api) DeleteTodoList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "DeleteTodoList"})

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}

	err = a.repo.DeleteTodoListByID(uint32(id))
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, []byte{})
}

// Todo

func (a *Api) Todo(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"path": TodoPath})

	switch req.Method {
	case http.MethodPost:
		a.CreateTodo(res, req)
	case http.MethodGet:
		a.GetTodosByList(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("method %q is not allowed", req.Method)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("method not allowed")
		return
	}

}

func (a *Api) CreateTodo(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "CreateTodo"})

	dec := json.NewDecoder(req.Body)
	todoReq := TodoTransport{}

	err := dec.Decode(&todoReq)
	if err != nil {
		msg := fmt.Sprintf("cant parse request body as JSON:%v", err)
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("invalid request body")
		return
	}

	vars := mux.Vars(req)
	listID, err := strconv.ParseUint(vars["list_id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "list_id", err)
		return
	}
	todoReq.ListID = uint32(listID)

	todoForInsert, err := fromTransportToTodo(todoReq)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
		logger.WithError(err).Warning("bad request error")
		return
	}

	newTodo, err := a.repo.InsertTodo(todoForInsert)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyDescription) {
			res.WriteHeader(http.StatusBadRequest)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("bad request error")
			return
		}
		if errors.Is(err, repository.ErrTodoListNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todoRes := toTransportTodo(*newTodo)
	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todoRes))
}

func (a *Api) GetTodosByList(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "GetTodosByList"})

	vars := mux.Vars(req)
	listID, err := strconv.ParseUint(vars["list_id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "list_id", err)
		return
	}

	todos, err := a.repo.GetTodosByListID(uint32(listID))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todosRes := []TodoTransport{}
	for _, t := range todos {
		todosRes = append(todosRes, toTransportTodo(t))
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todosRes))
}

func (a *Api) TodoByID(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"path": TodoIDPath})

	switch req.Method {
	case http.MethodGet:
		a.GetTodo(res, req)
	case http.MethodPut:
		a.UpdateTodo(res, req)
	case http.MethodDelete:
		a.DeleteTodo(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
		msg := fmt.Sprintf("method %q is not allowed", req.Method)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("method not allowed")
		return
	}

}

func (a *Api) GetTodo(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "GetTodo"})

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}

	todo, err := a.repo.GetTodoByID(uint32(id))
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	todoRes := toTransportTodo(*todo)
	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, toJSON(logger, todoRes))
}

func (a *Api) UpdateTodo(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "UpdateTodo"})

	dec := json.NewDecoder(req.Body)
	todoReq := TodoTransport{}

	err := dec.Decode(&todoReq)
	if err != nil {
		msg := fmt.Sprintf("cant parse request body as JSON:%v", err)
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
		logger.WithFields(log.Fields{"error": msg}).Warning("invalid request body")
		return
	}

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}
	todoReq.ID = uint32(id)

	todoForUpdate, err := fromTransportToTodo(todoReq)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
		logger.WithError(err).Warning("bad request error")
		return
	}

	err = a.repo.UpdateTodo(todoForUpdate)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyDescription) {
			res.WriteHeader(http.StatusBadRequest)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("bad request error")
			return
		}
		if errors.Is(err, repository.ErrTodoNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, []byte{})
}

func (a *Api) DeleteTodo(res http.ResponseWriter, req *http.Request) {
	logger := log.WithFields(log.Fields{"action": "DeleteTodo"})

	vars := mux.Vars(req)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "id", err)
		return
	}

	listID, err := strconv.ParseUint(vars["list_id"], 10, 32)
	if err != nil {
		handleFieldParsingError(logger, res, "list_id", err)
		return
	}

	todoForDelete := repository.Todo{
		ID:     uint32(id),
		ListID: uint32(listID),
	}

	err = a.repo.DeleteTodo(todoForDelete)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		if errors.Is(err, repository.ErrTodoNotFound) {
			res.WriteHeader(http.StatusNotFound)
			logResponseBodyWrite(logger, res, newErrorResponse(logger, err.Error()))
			logger.WithError(err).Warning("not found error")
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		logResponseBodyWrite(logger, res, newErrorResponse(logger, "internal server error"))
		logger.WithError(err).Error("internal server error")
		return
	}

	res.WriteHeader(http.StatusOK)
	logResponseBodyWrite(logger, res, []byte{})
}

func logResponseBodyWrite(logger *log.Entry, w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warning("writing response body")
	}
}

func newErrorResponse(logger *log.Entry, message string) []byte {
	return toJSON(logger, ErrorResponse{
		Error: Error{Message: message},
	})
}

func toJSON(logger *log.Entry, v interface{}) []byte {
	res, err := json.Marshal(v)
	if err != nil {
		logger.WithError(err).Warning("unable to marshal as JSON")
	}
	return res
}

func handleFieldParsingError(logger *log.Entry, res http.ResponseWriter, fieldName string, err error) {
	msg := fmt.Sprintf("can't parse %q from request:%v", fieldName, err)
	res.WriteHeader(http.StatusBadRequest)
	logResponseBodyWrite(logger, res, newErrorResponse(logger, msg))
	logger.WithError(err).WithFields(log.Fields{"field": fieldName}).Warning("invalid field on request")
}

func fromTransportToTodoList(ttl TodoListTransport) repository.TodoList {
	return repository.TodoList{
		ID:    ttl.ID,
		Title: ttl.Title,
	}
}

func toTransportTodoList(tl repository.TodoList) TodoListTransport {
	return TodoListTransport{
		ID:    tl.ID,
		Title: tl.Title,
	}
}

func fromTransportToTodo(tt TodoTransport) (repository.Todo, error) {
	todo := repository.Todo{
		ID:          tt.ID,
		ListID:      tt.ListID,
		Description: tt.Description,
		Comments:    tt.Comments,
		Labels:      tt.Labels,
		Done:        tt.Done,
	}

	if tt.DueDate != "" {
		dueDate, err := time.Parse(dateLayout, tt.DueDate)
		if err != nil {
			return repository.Todo{}, ErrInvalidDueDate
		}
		todo.DueDate = dueDate
	}

	return todo, nil
}

func toTransportTodo(t repository.Todo) TodoTransport {
	todoTransport := TodoTransport{
		ID:          t.ID,
		ListID:      t.ListID,
		Description: t.Description,
		Comments:    t.Comments,
		Labels:      t.Labels,
		Done:        t.Done,
	}

	if !t.DueDate.IsZero() {
		todoTransport.DueDate = t.DueDate.Format(dateLayout)
	}

	return todoTransport
}
