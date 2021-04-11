package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todoapp"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type TodoService interface {
	Create(description string) (todoapp.Todo, error)
	Read(id int) (todoapp.Todo, error)
	ReadAll() ([]todoapp.Todo, error)
	Update(id int, description string, done bool) (todoapp.Todo, error)
	Delete(id int) error
}

type TransportRest struct {
	Handler http.Handler
	BaseUrl string
	service TodoService
	logger  *logrus.Entry
	router  *httprouter.Router
}

// serverErrorResponse representa o padrão de resposta de erro da api
type serverErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

func NewTransportRest(ts TodoService, baseUrl string, logger *logrus.Entry) *TransportRest {
	r := httprouter.New()
	m := NewLoggingMiddleware(r, logger)
	tr := TransportRest{
		router:  r,
		Handler: m, // caso não use um Middleware, o Handler deve ser o próprio router
		BaseUrl: baseUrl,
		service: ts,
		logger:  logger,
	}
	tr.setHandlers()
	return &tr
}

func (tr *TransportRest) setHandlers() {
	tr.logger.Trace("Starting handler configuration")
	// Error
	tr.router.GET(tr.BaseUrl+"/error", tr.errorExample)

	// TO-DO
	tr.router.GET(tr.BaseUrl+"/todo", tr.readAllTodos)
	tr.router.POST(tr.BaseUrl+"/todo", tr.createTodo)
	tr.router.GET(tr.BaseUrl+"/todo/:id", tr.readTodo)
	tr.router.PUT(tr.BaseUrl+"/todo/:id", tr.updateTodo)
	tr.router.DELETE(tr.BaseUrl+"/todo/:id", tr.deleteTodo)

	tr.logger.Trace("Finalized configuration of the manipulators")
}

func (tr *TransportRest) errorExample(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tr.sendErrorResponse(http.StatusInternalServerError, "Example of an error", w, r)
}

func (tr *TransportRest) readAllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type TodoResponse struct {
		Todos []todoapp.Todo `json:"todos"`
	}

	todos, err := tr.service.ReadAll()
	if err != nil {
		tr.sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	tr.sendJsonResponse(http.StatusOK, TodoResponse{Todos: todos}, w, r)
}

func (tr *TransportRest) createTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type TodoRequest struct {
		Description *string `json:"description"`
	}

	var todoReq TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		tr.sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	if todoReq.Description == nil || *todoReq.Description == "" {
		tr.sendErrorResponse(http.StatusBadRequest, ErrInvalidRequestBody.Error(), w, r)
		return
	}

	newTodo, err := tr.service.Create(*todoReq.Description)
	if err != nil {
		tr.sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	tr.sendJsonResponse(http.StatusOK, newTodo, w, r)
}

func (tr *TransportRest) readTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	todo, err := tr.service.Read(id)
	if errors.Is(err, todoapp.ErrNotFound) {
		tr.sendErrorResponse(http.StatusNotFound, err.Error(), w, r)
		return
	} else if err != nil {
		tr.sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	tr.sendJsonResponse(http.StatusOK, todo, w, r)
}

func (tr *TransportRest) updateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type TodoRequest struct {
		Description *string `json:"description"`
		Done        *bool   `json:"done"`
	}

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	var todoReq TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		tr.sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	if todoReq.Description == nil || todoReq.Done == nil || *todoReq.Description == "" {
		tr.sendErrorResponse(http.StatusBadRequest, ErrInvalidRequestBody.Error(), w, r)
		return
	}

	updatedTodo, err := tr.service.Update(id, *todoReq.Description, *todoReq.Done)
	if errors.Is(err, todoapp.ErrNotFound) {
		tr.sendErrorResponse(http.StatusNotFound, err.Error(), w, r)
		return
	} else if err != nil {
		tr.sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	tr.sendJsonResponse(http.StatusOK, updatedTodo, w, r)
}

func (tr *TransportRest) deleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		tr.sendErrorResponse(http.StatusBadRequest, fmt.Sprintf("%s : %s", ErrDecodeRequestBody, err), w, r)
		return
	}

	err = tr.service.Delete(id)
	if errors.Is(err, todoapp.ErrNotFound) {
		tr.sendErrorResponse(http.StatusNotFound, err.Error(), w, r)
		return
	} else if err != nil {
		tr.sendErrorResponse(http.StatusInternalServerError, fmt.Sprintf("%s : %s", ErrUnknownService, err), w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (tr *TransportRest) sendJsonResponse(httpStatus int, v interface{}, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(v)
}

func (tr *TransportRest) sendErrorResponse(httpStatus int, errMessage string, w http.ResponseWriter, r *http.Request) error {
	tr.logger.Error(errMessage)
	errRes := serverErrorResponse{
		Timestamp: time.Now().UTC(),
		Status:    httpStatus,
		Error:     http.StatusText(httpStatus),
		Message:   errMessage,
		Path:      r.RequestURI,
	}
	return tr.sendJsonResponse(httpStatus, errRes, w, r)
}
