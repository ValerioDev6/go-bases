package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"web/internal/auth/dtos"
	"web/internal/middleware"
	httpresponse "web/internal/transport/http/response"
)

type Handler struct {
	authService *Service
}

func NewAuthHandler(authService *Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.Handle("GET /auth/me", middleware.RequireAuthentication(h.authService, http.HandlerFunc(h.Me)))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var request dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
		return
	}

	createdUser, err := h.authService.Register(r.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrNameRequired):
			httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("name", err.Error()))
		case errors.Is(err, ErrEmailRequired):
			httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("email", err.Error()))
		case errors.Is(err, ErrPasswordRequired), errors.Is(err, ErrPasswordTooShort), errors.Is(err, ErrPasswordTooLong):
			httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("password", err.Error()))
		case errors.Is(err, ErrEmailInUse):
			httpresponse.Failure(w, http.StatusConflict, httpresponse.FieldError("email", err.Error()))
		default:
			httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
		}
		return
	}

	response := dtos.RegisterResponse{
		ID:    createdUser.ID,
		Email: createdUser.Email,
		Name:  createdUser.Name,
	}
	httpresponse.Success(w, http.StatusCreated, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request dtos.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpresponse.Failure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
		return
	}

	token, err := h.authService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpresponse.Failure(w, http.StatusUnauthorized, httpresponse.GeneralError(err.Error()))
			return
		}
		httpresponse.Failure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
		return
	}

	httpresponse.Success(w, http.StatusOK, dtos.LoginResponse{Token: token})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.AuthenticatedUserID(r.Context())
	if err != nil {
		httpresponse.Failure(w, http.StatusUnauthorized, httpresponse.GeneralError(ErrInvalidToken.Error()))
		return
	}
	authenticatedUser, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		httpresponse.Failure(w, http.StatusUnauthorized, httpresponse.GeneralError(ErrInvalidToken.Error()))
		return
	}
	httpresponse.Success(w, http.StatusOK, authenticatedUser)
}
