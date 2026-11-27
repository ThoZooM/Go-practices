package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ThoZooM/Go-practices/practice-03/internal/api"
	"github.com/ThoZooM/Go-practices/practice-03/internal/storage"
)

func main() {
	// Создаём хранилище задач
	store := storage.NewMemoryStore()

	// Создаём HTTP-обработчики
	h := api.NewHandlers(store)

	// Создаём маршрутизатор
	mux := http.NewServeMux()

	// Проверка состояния сервера
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Работа со списком задач
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)

	// Получение задачи по ID
	mux.HandleFunc("GET /tasks/", h.GetTask)

	mux.HandleFunc("PATCH /tasks/{id}", h.PatchTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	// Подключаем middleware логирования
	handler := api.Logging(api.CORS(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
