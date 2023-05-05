package models

const IdNotIntegerErr = "id not integer"
const ValidationErr = "validation err"
const DbErr = "db error"
const NotFoundErr = "todo not found"

type Todo struct {
	Id          uint
	Description string
	IsCompleted bool
	DueDate     string
}

type TodoErr struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (t TodoErr) Error() string {
	return t.Message
}
