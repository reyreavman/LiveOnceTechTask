package task

import (
	"context"
	"liveoncetechtask/internal/models"
)

type Repository interface {
	CreateTask(ctx context.Context, task models.Task) models.Task
	GetTaskById(ctx context.Context, id string) (*models.Task, error)
	GetTasksByStatus(ctx context.Context, status models.Status) []models.Task
	GetAllTasks(ctx context.Context) []models.Task
}
