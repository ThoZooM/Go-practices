package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ThoZooM/Go-practices/practice-03/internal/storage"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

// GET /tasks
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()

	// Фильтрация по названию: ?q=text
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	if q != "" {
		filtered := tasks[:0]

		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}

		tasks = filtered
	}

	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

// POST /tasks
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" &&
		!strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}

	req.Title = strings.TrimSpace(req.Title)

	// Проверка обязательного поля
	if req.Title == "" {
		BadRequest(w, "title is required")
		return
	}

	// Проверка длины названия задачи
	titleLength := len([]rune(req.Title))

	if titleLength < 3 || titleLength > 140 {
		JSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Error: "title must contain 3 to 140 characters",
		})
		return
	}

	// Создание задачи
	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)

}

// GET /tasks/{id}
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 2 {
		NotFound(w, "invalid path")
		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)

	if err != nil {
		BadRequest(w, "invalid id")
		return
	}

	t, err := h.Store.Get(id)

	if err != nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusOK, t)
}

// PATCH /tasks/{id}
func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		BadRequest(w, "invalid id")
		return
	}

	task, err := h.Store.MarkDone(id)
	if err != nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusOK, task)
}

// DELETE /tasks/{id}
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		BadRequest(w, "invalid id")
		return
	}

	if err := h.Store.Delete(id); err != nil {
		NotFound(w, "task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
