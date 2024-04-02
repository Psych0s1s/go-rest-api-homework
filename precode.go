package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Printf("Ошибка при маршалинге задач: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := tasks[newTask.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	if newTask.ID == "" {
		http.Error(w, "Отсутствует ID для создания задачи", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest) //имхо, при такой проверке, статус: http.StatusNotFound (404) выглядел бы логичней
		return
	}

	resp, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //имхо, а здесь лучше бы http.StatusInternalServerError (500), наверное
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; !ok {
		w.WriteHeader(http.StatusBadRequest) //имхо, при такой проверке, статус: http.StatusNotFound (404) выглядел бы логичней
		return
	}

	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	// Запуск HTTP сервера на порту 8080 с обработкой ошибок
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err)
	}
}
