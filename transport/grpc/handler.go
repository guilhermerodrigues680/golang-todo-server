package transportgrpc

import (
	"context"
	"todoapp"
	"todoapp/transport/grpc/pbtodoapp"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoService interface {
	Create(description string) (todoapp.Todo, error)
	Read(id int) (todoapp.Todo, error)
	ReadAll() ([]todoapp.Todo, error)
	Update(id int, description string, done bool) (todoapp.Todo, error)
	Delete(id int) error
}

// FIXME remover pbtodoapp.UnimplementedTodoServiceServer
type TransportGRPC struct {
	pbtodoapp.UnimplementedTodoServiceServer
	service TodoService
	logger  *logrus.Entry
}

func NewTransportGRPC(ts TodoService, logger *logrus.Entry) *TransportGRPC {
	return &TransportGRPC{
		service: ts,
		logger:  logger,
	}
}

func (tg *TransportGRPC) Create(ctx context.Context, todoReq *pbtodoapp.TodoCreateRequest) (*pbtodoapp.Todo, error) {

	todo, err := tg.service.Create(todoReq.Description)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return TodoToPbTodo(todo), nil
	//return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (tg *TransportGRPC) Read(ctx context.Context, idReq *pbtodoapp.Id) (*pbtodoapp.Todo, error) {

	todo, err := tg.service.Read(int(idReq.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return TodoToPbTodo(todo), nil
}

// Adapters

func TodoToPbTodo(todo todoapp.Todo) *pbtodoapp.Todo {
	return &pbtodoapp.Todo{
		Id:          uint64(todo.ID),
		Description: todo.Description,
		Done:        todo.Done,
	}
}
