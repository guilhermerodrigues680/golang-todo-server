package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todoserver"
	"todoserver/storage"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type TodoService interface {
	Create(description string) (todoserver.Todo, error)
	Read(id int) (todoserver.Todo, error)
	ReadAll() ([]todoserver.Todo, error)
	Update(id int, description string, done bool) (todoserver.Todo, error)
	Delete(id int) error
}

type TodoHandlers interface {
}

type TransportRest struct {
	Service TodoService
	BaseUrl string
	Logger  *logrus.Entry
	router  *httprouter.Router
	Handler http.Handler
}

func NewTransportRest(ts TodoService, baseUrl string, logger *logrus.Entry) *TransportRest {
	r := httprouter.New()
	m := NewLoggingMiddleware(r, logger)
	tr := TransportRest{
		router:  r,
		Handler: m, // caso não use um Middleware, o Handler deve ser o próprio router
		BaseUrl: baseUrl,
		Service: ts,
		Logger:  logger,
	}
	tr.setHandlers()
	return &tr
}

func (tr *TransportRest) setHandlers() {
	tr.Logger.Trace("Starting handler configuration")
	// Error
	tr.router.GET(tr.BaseUrl+"/error", tr.errorExample)

	// TO DOs
	tr.router.GET(tr.BaseUrl+"/todo", tr.readAllTodos)
	tr.router.POST(tr.BaseUrl+"/todo", tr.createTodo)
	tr.router.GET(tr.BaseUrl+"/todo/:id", tr.readTodo)
	tr.router.PUT(tr.BaseUrl+"/todo/:id", tr.updateTodo)
	tr.router.DELETE(tr.BaseUrl+"/todo/:id", tr.deleteTodo)

	tr.Logger.Trace("Finalized configuration of the manipulators")
}

func (tr *TransportRest) errorExample(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendErrorResponse(http.StatusInternalServerError, "Example of an error", w, r)
}

func (tr *TransportRest) readAllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type TodoResponse struct {
		Todos []todoserver.Todo `json:"todos"`
	}

	todos, err := tr.Service.ReadAll()
	if err != nil {
		tr.Logger.Error(ErrUnknownService, err)
		sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TodoResponse{Todos: todos})
}

func (tr *TransportRest) createTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type TodoRequest struct {
		Description *string `json:"description"`
	}

	var todoReq TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		tr.Logger.Error(ErrDecodeRequestBody, err)
		sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	if todoReq.Description == nil || *todoReq.Description == "" {
		tr.Logger.Error(ErrInvalidRequestBody)
		sendErrorResponse(http.StatusBadRequest, ErrInvalidRequestBody.Error(), w, r)
		return
	}

	newTodo, err := tr.Service.Create(*todoReq.Description)
	if err != nil {
		tr.Logger.Error(ErrUnknownService, err)
		sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newTodo)
}

func (tr *TransportRest) readTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.Logger.Error(ErrDecodeRequestBody, err)
		sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	todo, err := tr.Service.Read(id)
	if errors.Is(err, storage.ErrNotFound) {
		tr.Logger.Error(err)
		sendErrorResponse(http.StatusBadRequest, err.Error(), w, r)
		return
	} else if err != nil {
		tr.Logger.Error(ErrUnknownService, err)
		sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func (tr *TransportRest) updateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type TodoRequest struct {
		Description *string `json:"description"`
		Done        *bool   `json:"done"`
	}

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.Logger.Error(ErrDecodeRequestBody, err)
		sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	var todoReq TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		tr.Logger.Error(ErrDecodeRequestBody, err)
		sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	if todoReq.Description == nil || todoReq.Done == nil || *todoReq.Description == "" {
		tr.Logger.Error(ErrInvalidRequestBody)
		sendErrorResponse(http.StatusBadRequest, ErrInvalidRequestBody.Error(), w, r)
		return
	}

	updatedTodo, err := tr.Service.Update(id, *todoReq.Description, *todoReq.Done)
	if errors.Is(err, storage.ErrNotFound) {
		tr.Logger.Error(err)
		sendErrorResponse(http.StatusBadRequest, err.Error(), w, r)
		return
	} else if err != nil {
		tr.Logger.Error(ErrUnknownService, err)
		sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTodo)
}

func (tr *TransportRest) deleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.Logger.Error(ErrDecodeRequestBody, err)
		sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	err = tr.Service.Delete(id)
	if errors.Is(err, storage.ErrNotFound) {
		tr.Logger.Error(err)
		sendErrorResponse(http.StatusBadRequest, err.Error(), w, r)
		return
	} else if err != nil {
		tr.Logger.Error(ErrUnknownService, err)
		sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// serverErrorResponse representa o padrão de resposta de erro da api
type serverErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func newServerErrorResponse(httpStatus int, message, path string) *serverErrorResponse {
	return &serverErrorResponse{
		Timestamp: time.Now().UTC(),
		Status:    httpStatus,
		Error:     http.StatusText(httpStatus),
		Message:   message,
		Path:      path,
	}
}

func sendErrorResponse(httpStatus int, errMessage string, w http.ResponseWriter, r *http.Request) error {
	errRes := newServerErrorResponse(httpStatus, errMessage, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(errRes)
}
