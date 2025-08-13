package task

import (
	"context"
	"liveoncetechtask/internal/models"
)

type Service interface {
	CreateTask(ctx context.Context, createTaskRequest models.CreateTaskRequest) models.TaskResponse
	GetTaskById(ctx context.Context, id string) (*models.TaskResponse, error)
	GetTasksByStatus(ctx context.Context, status *models.Status) []models.Task
}
