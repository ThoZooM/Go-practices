package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/ThoZooM/Go-practices/practice-04/internal/task"
	myMW "github.com/ThoZooM/Go-practices/practice-04/pkg/middleware"
)

func main() {
	// Создаём хранилище задач

	repo := task.NewRepo()

	// Загружаем сохранённые задачи при запуске
	if err := repo.Load(); err != nil {
		log.Fatalf("failed to load tasks: %v", err)
	}

	h := task.NewHandler(repo)

	// Создаём роутер chi
	r := chi.NewRouter()

	// Подключаем middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(myMW.Logger)
	r.Use(myMW.SimpleCORS)

	// Проверка состояния сервера
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Группа маршрутов API

	// Версионирование API
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/tasks", h.Routes())
		})
	})

	// Запуск сервера
	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
