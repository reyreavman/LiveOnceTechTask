package http

import (
	"context"
	"encoding/json"
	"errors"
	"liveoncetechtask/internal/logger"
	"liveoncetechtask/internal/models"
	"liveoncetechtask/internal/task"
	"net/http"
	"strings"
	"time"
)

type TaskHandler struct {
	service task.Service
	logger  *logger.Logger

	headers map[string]string
}

func NewTaskHandler(service task.Service, logger *logger.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	h.logger.Info("create_task started", "", nil)

	var createTaskRequest models.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&createTaskRequest); err != nil {
		h.logger.Error("create_task: decode failed", "", err, nil)

		error := models.Error{
			Type:    "invalid request body",
			Message: err.Error(),
		}
		http.Error(
			w,
			models.ErrorToJSON(error),
			http.StatusBadRequest,
		)
		return
	}

	if err := models.ValidateCreateTaskRequest(createTaskRequest); err != nil {
		h.logger.Error("create_task: validation failed", "", err, map[string]interface{}{
			"request": createTaskRequest,
		})

		error := models.Error{
			Type:    "invalid request body",
			Message: err.Error(),
		}
		http.Error(
			w,
			models.ErrorToJSON(error),
			http.StatusBadRequest,
		)
		return
	}
	taskResponse := h.service.CreateTask(ctx, createTaskRequest)

	h.logger.Info("create_task completed", "", map[string]interface{}{
		"task_id": taskResponse.Id,
		"status":  taskResponse.Status,
	})

	constructResponse(h.headers, http.StatusCreated, w, taskResponse)
}

func (h *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	taskId := strings.TrimPrefix(r.URL.Path, PATTERN_TASK_BY_ID)

	taskResponse, err := h.service.GetTaskById(ctx, taskId)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			h.logger.Warning("get_task_by_id: task not found", "", map[string]interface{}{
				"task_id": taskId,
			})
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			h.logger.Error("get_task_by_id: failed", "", err, map[string]interface{}{
				"task_id": taskId,
			})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info("get_task_by_id completed", "", map[string]interface{}{
		"task_id": taskId,
	})

	constructResponse(h.headers, http.StatusOK, w, taskResponse)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status := models.Status(r.URL.Query().Get("status"))
	taskResponse := h.service.GetTasksByStatus(ctx, status)
	h.logger.Info("get_tasks completed", "", map[string]interface{}{
		"count":  len(taskResponse),
		"status": status,
	})
	
	constructResponse(h.headers, http.StatusOK, w, taskResponse)
}

func constructResponse(headers map[string]string, status int, w http.ResponseWriter, responseBody any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseBody)

	for k, v := range headers {
		w.Header().Set(k, v)
	}
}
