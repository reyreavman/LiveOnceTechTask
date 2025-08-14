package repository

import (
	"context"
	"liveoncetechtask/internal/models"
	"liveoncetechtask/internal/task"
	"sync"
)

type TaskRepository struct {
	mu     sync.RWMutex
	store  map[string]models.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		store:  make(map[string]models.Task),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task models.Task) models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[task.Id] = task

	return task
}

func (r *TaskRepository) GetTaskById(ctx context.Context, id string) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskToReturn, exists := r.store[id]
	if exists {
		return &taskToReturn, nil
	} else {
		return nil, task.ErrTaskNotFound
	}
}

func (r *TaskRepository) GetTasksByStatus(ctx context.Context, status models.Status) []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Task, 0, len(r.store))
	for _, task := range r.store {
		if task.Status == status {
			result = append(result, task)
		}
	}

	return result
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Task, 0, len(r.store))
	for _, task := range r.store {
		result = append(result, task)
	}

	return result
}
