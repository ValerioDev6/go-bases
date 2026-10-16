package role

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"web/internal/role/dtos"
	httpresponse "web/internal/transport/http/response"
)

const maxRequestBody = 1 << 20

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/roles", h.handleCollection)
	mux.HandleFunc("/roles/", h.handleItem)
}

func (h *Handler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		roles, err := h.service.GetAll(r.Context())
		if err != nil {
			h.internalError(w, r, err)
			return
		}
		h.writeSuccess(w, http.StatusOK, roles)
	case http.MethodPost:
		request, err := h.decodeRole(w, r)
		if err != nil {
			h.writeFailure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
			return
		}
		created, err := h.service.Create(r.Context(), Role{Name: request.Name})
		if err != nil {
			h.serviceError(w, r, err)
			return
		}
		h.writeSuccess(w, http.StatusCreated, created)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		h.writeFailure(w, http.StatusMethodNotAllowed, httpresponse.GeneralError("method not allowed"))
	}
}

func (h *Handler) handleItem(w http.ResponseWriter, r *http.Request) {
	id, ok := h.idFromPath(w, r.URL.Path)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		role, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			h.serviceError(w, r, err)
			return
		}
		h.writeSuccess(w, http.StatusOK, role)
	case http.MethodPut:
		request, err := h.decodeRole(w, r)
		if err != nil {
			h.writeFailure(w, http.StatusBadRequest, httpresponse.FieldError("body", "invalid request body"))
			return
		}
		updated, err := h.service.Update(r.Context(), id, Role{Name: request.Name})
		if err != nil {
			h.serviceError(w, r, err)
			return
		}
		h.writeSuccess(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), id); err != nil {
			h.serviceError(w, r, err)
			return
		}
		h.writeSuccess(w, http.StatusOK, nil)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPut+", "+http.MethodDelete)
		h.writeFailure(w, http.StatusMethodNotAllowed, httpresponse.GeneralError("method not allowed"))
	}
}

func (h *Handler) idFromPath(w http.ResponseWriter, path string) (int64, bool) {
	value := strings.TrimPrefix(path, "/roles/")
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 || strings.Contains(value, "/") {
		h.writeFailure(w, http.StatusBadRequest, httpresponse.FieldError("id", "invalid role id"))
		return 0, false
	}
	return id, true
}

func (h *Handler) decodeRole(w http.ResponseWriter, r *http.Request) (dtos.RoleRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request dtos.RoleRequest
	if err := decoder.Decode(&request); err != nil {
		return dtos.RoleRequest{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return dtos.RoleRequest{}, errors.New("request body must contain one JSON object")
	}
	return request, nil
}

func (h *Handler) serviceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNameRequired):
		h.writeFailure(w, http.StatusBadRequest, httpresponse.FieldError("name", err.Error()))
	case errors.Is(err, ErrNotFound):
		h.writeFailure(w, http.StatusNotFound, httpresponse.GeneralError(err.Error()))
	case errors.Is(err, ErrConflict):
		h.writeFailure(w, http.StatusConflict, httpresponse.FieldError("name", err.Error()))
	default:
		h.internalError(w, r, err)
	}
}

func (h *Handler) internalError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("role request failed", "method", r.Method, "path", r.URL.Path, "error", err)
	h.writeFailure(w, http.StatusInternalServerError, httpresponse.GeneralError("internal server error"))
}

func (h *Handler) writeSuccess(w http.ResponseWriter, status int, data any) {
	if err := httpresponse.Success(w, status, data); err != nil {
		h.logger.Error("encode role response", "error", err)
	}
}

func (h *Handler) writeFailure(w http.ResponseWriter, status int, responseErrors ...httpresponse.Error) {
	if err := httpresponse.Failure(w, status, responseErrors...); err != nil {
		h.logger.Error("encode role response", "error", err)
	}
}
