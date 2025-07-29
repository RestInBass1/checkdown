package httpHandlers

import (
	"checkdown/apiService/internal/pkg/logger"
	"checkdown/apiService/internal/usecase"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// вход POST /tasks
type createReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ответ POST /tasks
type createResp struct {
	ID int64 `json:"id"`
}

// строка → int64
func atoi(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

// единообразно пишем JSON‑ответ
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json)")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

type Handler struct {
	svc usecase.TaskService
}

func New(svc usecase.TaskService) *Handler { return &Handler{svc} }

func (h *Handler) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(Logging)
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Put("/{id}/done", h.done)
		r.Delete("/{id}", h.delete)
	})
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in createReq
	_ = json.NewDecoder(r.Body).Decode(&in)
	logger.Log.Debugw("create request", "title", in.Title)
	id, err := h.svc.Create(r.Context(), usecase.Task{
		Title:       in.Title,
		Description: in.Description,
	})
	if err != nil {
		logger.Log.Errorw("create task failed", "err", err)
		writeJSON(w, http.StatusInternalServerError, err)
	}
	logger.Log.Infow("task created", "id", id)
	writeJSON(w, http.StatusCreated, createResp{id})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err = h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) done(w http.ResponseWriter, r *http.Request) {
	id, err := atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err = h.svc.MarkDone(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
