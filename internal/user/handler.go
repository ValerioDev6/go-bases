package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	httpresponse "web/internal/transport/http/response"
	"web/internal/user/dtos"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.HandleUsers)
	mux.HandleFunc("/users/", h.HandleUserByID)
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.service.GetAllUsers(r.Context())
		if err != nil {
			httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
			return
		}
		httpresponse.Success(w, http.StatusOK, users)

	case http.MethodPost:
		var request dtos.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
			return
		}

		user := &User{
			RoleID: request.RoleID,
			Name:   request.Name,
			Email:  request.Email,
		}
		createdUser, err := h.service.CreateUser(r.Context(), user)
		if err != nil {
			switch {
			case errors.Is(err, ErrNameRequired):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("name", err.Error()))
			case errors.Is(err, ErrEmailRequired):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("email", err.Error()))
			case errors.Is(err, ErrRoleRequired), errors.Is(err, ErrRoleNotFound):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("role_id", err.Error()))
			case errors.Is(err, ErrConflict):
				httpresponse.Failure(w, http.StatusConflict, httpresponse.FieldError("email", err.Error()))
			default:
				httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
			}
			return
		}
		httpresponse.Success(w, http.StatusCreated, createdUser)

	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		httpresponse.Failure(w, http.StatusMethodNotAllowed, httpresponse.GeneralError("method not allowed"))
	}
}

func (h *Handler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("id", "invalid user id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.service.GetUserByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				httpresponse.Failure(w, http.StatusNotFound, httpresponse.GeneralError(err.Error()))
				return
			}
			httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
			return
		}
		httpresponse.Success(w, http.StatusOK, user)

	case http.MethodPut:
		var request dtos.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
			return
		}

		user := &User{
			RoleID: request.RoleID,
			Name:   request.Name,
			Email:  request.Email,
		}
		updatedUser, err := h.service.UpdateUser(r.Context(), id, user)
		if err != nil {
			switch {
			case errors.Is(err, ErrNameRequired):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("name", err.Error()))
			case errors.Is(err, ErrEmailRequired):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("email", err.Error()))
			case errors.Is(err, ErrRoleRequired), errors.Is(err, ErrRoleNotFound):
				httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("role_id", err.Error()))
			case errors.Is(err, ErrNotFound):
				httpresponse.Failure(w, http.StatusNotFound, httpresponse.GeneralError(err.Error()))
			case errors.Is(err, ErrConflict):
				httpresponse.Failure(w, http.StatusConflict, httpresponse.FieldError("email", err.Error()))
			default:
				httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
			}
			return
		}
		httpresponse.Success(w, http.StatusOK, updatedUser)

	case http.MethodDelete:
		err := h.service.DeleteUser(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				httpresponse.Failure(w, http.StatusNotFound, httpresponse.GeneralError(err.Error()))
				return
			}
			httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
			return
		}
		httpresponse.Success(w, http.StatusOK, nil)

	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPut+", "+http.MethodDelete)
		httpresponse.Failure(w, http.StatusMethodNotAllowed, httpresponse.GeneralError("method not allowed"))
	}
}
