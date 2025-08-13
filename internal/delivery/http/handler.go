package http

import (
	"context"
	"encoding/json"
	"errors"
	"liveoncetechtask/internal/models"
	"liveoncetechtask/internal/task"
	"net/http"
	"strings"
	"time"
)

type TaskHandler struct {
	service task.Service
	headers map[string]string
}

func NewTaskHandler(service task.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
		headers: map[string]string{
		"Content-Type": "application/json",
		},
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var creatTaskRequest models.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&creatTaskRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	taskResponse := h.service.CreateTask(ctx, creatTaskRequest)

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(taskResponse)

	constructResponse(h.headers, http.StatusCreated, w, taskResponse)
}

func (h *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	taskId := strings.TrimPrefix(r.URL.Path, TASK_BY_ID)

	taskResponse, err := h.service.GetTaskById(ctx, taskId)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(taskResponse)

	constructResponse(h.headers, http.StatusOK, w, taskResponse)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status := models.Status(r.URL.Query().Get("status"))
	taskResponse := h.service.GetTasksByStatus(ctx, status)
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(taskResponse)

	constructResponse(h.headers, http.StatusOK, w, taskResponse)
}

func constructResponse(headers map[string]string, status int, w http.ResponseWriter, responseBody any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseBody)

	for k, v := range headers {
		w.Header().Set(k ,v)
	}
}