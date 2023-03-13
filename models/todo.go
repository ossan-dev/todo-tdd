package models

const (
	ValidationErr   = "validation err"
	IdNotIntegerErr = "id not integer"
	TodoNotFoundErr = "unknown todo"
	DbErr           = "database error"
)

type TodoErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (t TodoErr) Error() string {
	return t.Message
}

type Todo struct {
	ID          int
	Description string
	IsCompleted bool
	DueDate     string
}

func NewTodo(id int, description string, isCompleted bool, dueDate string) Todo {
	return Todo{
		ID:          id,
		Description: description,
		IsCompleted: isCompleted,
		DueDate:     dueDate,
	}
}
