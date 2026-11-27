package task

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu       sync.RWMutex
	seq      int64
	items    map[int64]*Task
	filePath string
}

func NewRepo() *Repo {
	return &Repo{
		items:    make(map[int64]*Task),
		filePath: "tasks.json",
	}
}

// Загружаем задачи из файла при запуске
func (r *Repo) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	r.items = make(map[int64]*Task)
	r.seq = 0

	for _, t := range tasks {
		if t == nil {
			return errors.New("invalid task in JSON file")
		}

		r.items[t.ID] = t

		if t.ID > r.seq {
			r.seq = t.ID
		}
	}

	return nil
}

// Сохраняем задачи в JSON-файл.
// Вызывать только при удержании блокировки r.mu.
func (r *Repo) saveLocked() error {
	tasks := make([]*Task, 0, len(r.items))

	for _, t := range r.items {
		tasks = append(tasks, t)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0600)
}

// Получение всех задач
func (r *Repo) List() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*Task, 0, len(r.items))

	for _, t := range r.items {
		copyTask := *t
		out = append(out, &copyTask)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out
}

// Получение задачи по ID
func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	copyTask := *t
	return &copyTask, nil
}

// Создание новой задачи
func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	now := time.Now()

	t := &Task{
		ID:        r.seq,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.items[t.ID] = t

	if err := r.saveLocked(); err != nil {
		log.Printf("failed to save tasks: %v", err)
	}

	copyTask := *t
	return &copyTask
}

// Обновление задачи
func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()

	if err := r.saveLocked(); err != nil {
		log.Printf("failed to save tasks: %v", err)
	}

	copyTask := *t
	return &copyTask, nil
}

// Удаление задачи
func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}

	delete(r.items, id)

	if err := r.saveLocked(); err != nil {
		log.Printf("failed to save tasks: %v", err)
	}

	return nil
}
