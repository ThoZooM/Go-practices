package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ThoZooM/Go-practices/practice-03/internal/storage"
)

// Проверка получения пустого списка задач
func TestListTasks(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	h.ListTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("Expected empty list, got %s", rec.Body.String())
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected JSON Content-Type")
	}
}

// Проверка создания задачи
func TestCreateTask(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	body := `{"title":"Test task"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Test task") {
		t.Errorf("Task was not created: %s", rec.Body.String())
	}

	tasks := store.List()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}
