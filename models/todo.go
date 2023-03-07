package models

const (
	ValidationErr   = "validation err"
	IdNotIntegerErr = "id not integer"
)

type TodoErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
