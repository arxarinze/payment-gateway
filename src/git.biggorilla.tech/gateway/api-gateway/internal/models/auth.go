package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GenericResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
