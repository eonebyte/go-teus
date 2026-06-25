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
	r *http.Request,
	data interface{},
) {
	render.Status(r, http.StatusCreated)

	render.JSON(w, r, APIResponse{
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
	r *http.Request,
	err error,
) {
	render.Status(r, http.StatusInternalServerError)

	render.JSON(w, r, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func BadRequest(
	w http.ResponseWriter,
	r *http.Request,
	message string,
) {
	render.Status(r, http.StatusBadRequest)

	render.JSON(w, r, APIResponse{
		Success: false,
		Error:   message,
	})
}

func NotFound(
	w http.ResponseWriter,
	r *http.Request,
	message string,
) {
	render.Status(r, http.StatusNotFound)

	render.JSON(w, r, APIResponse{
		Success: false,
		Error:   message,
	})
}

func Unauthorized(
	w http.ResponseWriter,
	r *http.Request,
	message string,
) {
	render.Status(r, http.StatusUnauthorized)

	render.JSON(w, r, APIResponse{
		Success: false,
		Error:   message,
	})
}
