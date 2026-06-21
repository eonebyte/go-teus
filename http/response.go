package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data interface{},
) {
	render.Status(r, status)

	render.JSON(w, r, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Created(
	w http.ResponseWriter,
	data interface{},
) {
	render.Status(nil, http.StatusCreated)

	render.JSON(w, nil, APIResponse{
		Success: true,
		Message: "created successfully",
		Data:    data,
	})
}

func NoContent(
	w http.ResponseWriter,
) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(
	w http.ResponseWriter,
	err error,
) {
	render.Status(nil, http.StatusInternalServerError)

	render.JSON(w, nil, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func BadRequest(
	w http.ResponseWriter,
	message string,
) {
	render.Status(nil, http.StatusBadRequest)

	render.JSON(w, nil, APIResponse{
		Success: false,
		Error:   message,
	})
}

func NotFound(
	w http.ResponseWriter,
	message string,
) {
	render.Status(nil, http.StatusNotFound)

	render.JSON(w, nil, APIResponse{
		Success: false,
		Error:   message,
	})
}

func Unauthorized(
	w http.ResponseWriter,
	message string,
) {
	render.Status(nil, http.StatusUnauthorized)

	render.JSON(w, nil, APIResponse{
		Success: false,
		Error:   message,
	})
}
