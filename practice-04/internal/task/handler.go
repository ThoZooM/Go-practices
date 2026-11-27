package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Регистрация CRUD-маршрутов
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}

// GET /tasks — список задач

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	// Значения пагинации по умолчанию
	page := 1
	limit := 10

	// Получаем номер страницы
	if value := r.URL.Query().Get("page"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			httpError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = n
	}

	// Получаем количество задач на странице
	if value := r.URL.Query().Get("limit"); value != "" {
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 || n > 100 {
			httpError(w, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
		limit = n
	}

	// Получаем все задачи (отсортированные по ID)
	tasks := h.repo.List()

	// Фильтрация по статусу done
	if values, exists := r.URL.Query()["done"]; exists {
		if len(values) != 1 {
			httpError(w, http.StatusBadRequest, "invalid done parameter")
			return
		}

		done, err := strconv.ParseBool(values[0])
		if err != nil {
			httpError(w, http.StatusBadRequest, "done must be true or false")
			return
		}

		filtered := make([]*Task, 0)
		for _, t := range tasks {
			if t.Done == done {
				filtered = append(filtered, t)
			}
		}

		tasks = filtered
	}

	// Пагинация после фильтрации
	start := (page - 1) * limit

	if start >= len(tasks) || start < 0 {
		writeJSON(w, http.StatusOK, []*Task{})
		return
	}

	end := start + limit
	if end > len(tasks) {
		end = len(tasks)
	}

	writeJSON(w, http.StatusOK, tasks[start:end])
}

// GET /tasks/{id} — получение задачи
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}

	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func validateTitle(title string) (string, bool) {
	title = strings.TrimSpace(title)

	length := len([]rune(title))

	if length < 3 || length > 100 {
		return title, false
	}

	return title, true
}

// POST /tasks — создание задачи

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json")
		return
	}

	title, valid := validateTitle(req.Title)
	if !valid {
		httpError(w, http.StatusUnprocessableEntity,
			"title must contain 3 to 100 characters")
		return
	}

	t := h.repo.Create(title)
	writeJSON(w, http.StatusCreated, t)
}

type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// PUT /tasks/{id} — обновление задачи

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}

	var req updateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json")
		return
	}

	title, valid := validateTitle(req.Title)
	if !valid {
		httpError(w, http.StatusUnprocessableEntity,
			"title must contain 3 to 100 characters")
		return
	}

	t, err := h.repo.Update(id, title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// DELETE /tasks/{id} — удаление задачи
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Получение ID задачи из URL
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, true
	}

	return id, false
}

// Формирование JSON-ответа
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// Формирование JSON-ошибки
func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
