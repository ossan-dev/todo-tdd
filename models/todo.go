package models

const (
	ValidationErr   = "validation err"
	IdNotIntegerErr = "id not integer"
	TodoNotFoundErr = "unknown todo"
)

type TodoErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
