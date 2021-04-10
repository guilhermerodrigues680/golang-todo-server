package todoserver

// Representa um TO DO no sistema
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TodoStorage interface {
	Create(description string) (Todo, error)
	Read(id int) (Todo, error)
	ReadAll() ([]Todo, error)
	Update(id int, description string, done bool) (Todo, error)
	Delete(id int) error
}

type TodoService struct {
	s TodoStorage
}

func NewTodoService(s TodoStorage) *TodoService {
	return &TodoService{s: s}
}

func (ts *TodoService) Create(description string) (Todo, error) {
	return ts.s.Create(description)
}

func (ts *TodoService) Read(id int) (Todo, error) {
	return ts.s.Read(id)
}

func (ts *TodoService) ReadAll() ([]Todo, error) {
	return ts.s.ReadAll()
}

func (ts *TodoService) Update(id int, description string, done bool) (Todo, error) {
	return ts.s.Update(id, description, done)
}

func (ts *TodoService) Delete(id int) error {
	return ts.s.Delete(id)
}
