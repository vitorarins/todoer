package api

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vitorarins/todoer/pb"
	"github.com/vitorarins/todoer/repository"
)

var (
	ErrInvalidDueDate = errors.New("due_date is invalid")
)

type GrpcApi struct {
	repo repository.Repository
	pb.UnimplementedTodoerServer
}

func NewGrpcApi(repo repository.Repository) *GrpcApi {
	return &GrpcApi{
		repo: repo,
	}
}

func (ga *GrpcApi) CreateTodoList(ctx context.Context, req *pb.CreateTodoListRequest) (*pb.CreateTodoListReply, error) {
	logger := log.WithFields(log.Fields{"action": "CreateTodoList"})

	todoListReq := repository.TodoList{
		Title: req.Title,
	}

	newTodoList, err := ga.repo.InsertTodoList(todoListReq)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyTitle) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	reply := &pb.CreateTodoListReply{
		TodoList: toProtoTodoList(*newTodoList),
	}
	return reply, nil
}

func (ga *GrpcApi) GetAllTodoLists(ctx context.Context, req *pb.Empty) (*pb.GetAllTodoListsReply, error) {
	logger := log.WithFields(log.Fields{"action": "GetAllTodoLists"})

	todoLists, err := ga.repo.GetAllTodoLists()
	if err != nil {
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	todoListsReply := []*pb.TodoList{}
	for _, tl := range todoLists {
		tlReply := toProtoTodoList(tl)
		todoListsReply = append(todoListsReply, tlReply)
	}

	reply := &pb.GetAllTodoListsReply{
		TodoLists: todoListsReply,
	}
	return reply, nil
}

func (ga *GrpcApi) GetTodoList(ctx context.Context, req *pb.GetTodoListRequest) (*pb.GetTodoListReply, error) {
	logger := log.WithFields(log.Fields{"action": "GetTodoList"})

	todoList, err := ga.repo.GetTodoListByID(req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	todoListReply := toProtoTodoList(*todoList)

	reply := &pb.GetTodoListReply{
		TodoList: todoListReply,
	}
	return reply, nil
}

func (ga *GrpcApi) UpdateTodoList(ctx context.Context, req *pb.UpdateTodoListRequest) (*pb.Empty, error) {
	logger := log.WithFields(log.Fields{"action": "UpdateTodoList"})

	todoListReq := fromProtoTodoList(req.TodoList)

	err := ga.repo.UpdateTodoList(todoListReq)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyTitle) ||
			errors.Is(err, repository.ErrTodoListNotFound) {

			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (ga *GrpcApi) DeleteTodoList(ctx context.Context, req *pb.DeleteTodoListRequest) (*pb.Empty, error) {
	logger := log.WithFields(log.Fields{"action": "DeleteTodoList"})

	err := ga.repo.DeleteTodoListByID(req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	return &pb.Empty{}, nil
}

// Todo

func (ga *GrpcApi) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	logger := log.WithFields(log.Fields{"action": "CreateTodo"})

	todoReq := repository.Todo{
		ListID:      req.ListId,
		Description: req.Description,
		Comments:    req.Comments,
		Labels:      req.Labels,
		Done:        req.Done,
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			return nil, ErrInvalidDueDate
		}
		todoReq.DueDate = dueDate
	}

	newTodo, err := ga.repo.InsertTodo(todoReq)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyDescription) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	reply := &pb.CreateTodoReply{
		Todo: toProtoTodo(*newTodo),
	}
	return reply, nil
}

func (ga *GrpcApi) GetTodosByList(ctx context.Context, req *pb.GetTodosByListRequest) (*pb.GetTodosByListReply, error) {
	logger := log.WithFields(log.Fields{"action": "GetTodosByList"})

	todos, err := ga.repo.GetTodosByListID(req.ListId)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	todosReply := []*pb.Todo{}
	for _, t := range todos {
		tReply := toProtoTodo(t)
		todosReply = append(todosReply, tReply)
	}

	reply := &pb.GetTodosByListReply{
		Todos: todosReply,
	}
	return reply, nil
}

func (ga *GrpcApi) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	logger := log.WithFields(log.Fields{"action": "GetTodo"})

	todo, err := ga.repo.GetTodoByID(req.Id)
	if err != nil {
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	todoReply := toProtoTodo(*todo)

	reply := &pb.GetTodoReply{
		Todo: todoReply,
	}
	return reply, nil
}

func (ga *GrpcApi) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.Empty, error) {
	logger := log.WithFields(log.Fields{"action": "UpdateTodo"})

	todoReq, err := fromProtoTodo(req.Todo)
	if err != nil {
		return nil, err
	}

	err = ga.repo.UpdateTodo(todoReq)
	if err != nil {
		if errors.Is(err, repository.ErrEmptyDescription) ||
			errors.Is(err, repository.ErrTodoNotFound) {

			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (ga *GrpcApi) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.Empty, error) {
	logger := log.WithFields(log.Fields{"action": "DeleteTodo"})

	todoReq := repository.Todo{
		ID:     req.Id,
		ListID: req.ListId,
	}
	err := ga.repo.DeleteTodo(todoReq)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			logger.WithError(err).Warning("bad request error")
			return nil, err
		}
		logger.WithError(err).Error("internal server error")
		return nil, err
	}

	return &pb.Empty{}, nil
}

func fromProtoTodoList(ptl *pb.TodoList) repository.TodoList {
	return repository.TodoList{
		ID:    ptl.Id,
		Title: ptl.Title,
	}
}

func fromProtoTodo(pt *pb.Todo) (repository.Todo, error) {

	todo := repository.Todo{
		ID:          pt.Id,
		ListID:      pt.ListId,
		Description: pt.Description,
		Comments:    pt.Comments,
		Labels:      pt.Labels,
		Done:        pt.Done,
	}

	if pt.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, pt.DueDate)
		if err != nil {
			return repository.Todo{}, ErrInvalidDueDate
		}
		todo.DueDate = dueDate
	}

	return todo, nil
}

func toProtoTodoList(tl repository.TodoList) *pb.TodoList {
	return &pb.TodoList{
		Id:    tl.ID,
		Title: tl.Title,
	}
}

func toProtoTodo(todo repository.Todo) *pb.Todo {
	return &pb.Todo{
		Id:          todo.ID,
		ListId:      todo.ListID,
		Description: todo.Description,
		Comments:    todo.Comments,
		Labels:      todo.Labels,
		DueDate:     todo.DueDate.String(),
		Done:        todo.Done,
	}
}
