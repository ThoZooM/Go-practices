package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Создаём тестовый HTTP-роутер
// Каждый тест использует собственный временный JSON-файл
func setupTestRouter(t *testing.T) *chi.Mux {
	t.Helper()

	repo := NewRepo()
	repo.filePath = t.TempDir() + "/tasks.json"

	handler := NewHandler(repo)

	r := chi.NewRouter()
	r.Mount("/api/v1/tasks", handler.Routes())

	return r
}

// Тест получения пустого списка задач
func TestListTasks(t *testing.T) {
	r := setupTestRouter(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks",
		nil,
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tasks []*Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected empty list, got %d tasks", len(tasks))
	}
}

// Тест создания задачи
func TestCreateTask(t *testing.T) {
	r := setupTestRouter(t)

	body := []byte(`{"title":"Test task"}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var result Task
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 || result.Title != "Test task" || result.Done {
		t.Fatalf("unexpected task: %+v", result)
	}
}

// Тест получения задачи по ID
func TestGetTask(t *testing.T) {
	r := setupTestRouter(t)

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		bytes.NewBufferString(`{"title":"Created task"}`),
	)
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("create failed: %d", createW.Code)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/1",
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result Task
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 || result.Title != "Created task" {
		t.Fatalf("unexpected task: %+v", result)
	}
}

// Тест обновления задачи
func TestUpdateTask(t *testing.T) {
	r := setupTestRouter(t)

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		bytes.NewBufferString(`{"title":"Original task"}`),
	)
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("create failed: %d", createW.Code)
	}

	body := []byte(`{"title":"Updated task","done":true}`)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/1",
		bytes.NewReader(body),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result Task
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result.Title != "Updated task" || !result.Done {
		t.Fatalf("unexpected updated task: %+v", result)
	}
}

// Тест удаления задачи
func TestDeleteTask(t *testing.T) {
	r := setupTestRouter(t)

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		bytes.NewBufferString(`{"title":"Delete task"}`),
	)
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("create failed: %d", createW.Code)
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/1",
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	getReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/1",
		nil,
	)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deletion, got %d", getW.Code)
	}
}

// Тест проверки длины title
func TestInvalidTitle(t *testing.T) {
	r := setupTestRouter(t)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		bytes.NewBufferString(`{"title":"Hi"}`),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// Тест некорректного ID
func TestInvalidID(t *testing.T) {
	r := setupTestRouter(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/abc",
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
