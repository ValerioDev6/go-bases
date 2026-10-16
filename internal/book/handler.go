package book

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const maxRequestBody = 1 << 20

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/books", h.handleBooks)
	mux.HandleFunc("/books/", h.handleBookByID)
}

func (h *Handler) handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleBookByID(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 || strings.Contains(idText, "/") {
		h.writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPut+", "+http.MethodDelete)
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAll(r.Context())
	if err != nil {
		h.internalError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, books)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	book, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.serviceError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, book)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	book, err := h.decodeBook(w, r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	created, err := h.service.Create(r.Context(), book)
	if err != nil {
		h.serviceError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, id int) {
	book, err := h.decodeBook(w, r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	updated, err := h.service.Update(r.Context(), id, book)
	if err != nil {
		h.serviceError(w, r, err)
		return
	}
	h.writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.serviceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) decodeBook(w http.ResponseWriter, r *http.Request) (Book, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var book Book
	if err := decoder.Decode(&book); err != nil {
		return Book{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Book{}, errors.New("request body must contain one JSON object")
	}
	return book, nil
}

func (h *Handler) serviceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrTitleRequired):
		h.writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		h.writeError(w, http.StatusNotFound, err.Error())
	default:
		h.internalError(w, r, err)
	}
}

func (h *Handler) internalError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("book request failed", "method", r.Method, "path", r.URL.Path, "error", err)
	h.writeError(w, http.StatusInternalServerError, "internal server error")
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		h.logger.Error("encode response", "error", err)
	}
}
