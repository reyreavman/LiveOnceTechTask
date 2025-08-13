package service

import (
	"context"
	"liveoncetechtask/internal/models"
	"liveoncetechtask/internal/task"
	"liveoncetechtask/pkg/id"
	"time"
)

type TaskService struct {
	taskRepository task.Repository
	idGenerator    id.Generator
}

func NewTaskRepository(taskRepository task.Repository) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, createTaskRequest models.CreateTaskRequest) models.TaskResponse {
	taskToSave := models.ToTask(createTaskRequest, s.idGenerator.Generate(), time.Now())
	savedTask := s.taskRepository.CreateTask(ctx, taskToSave)

	return models.ToTaskResponse(savedTask)
}

func (s *TaskService) GetTaskById(ctx context.Context, id string) (*models.TaskResponse, error) {
	task, err := s.taskRepository.GetTaskById(ctx, id)
	if err != nil {
		return nil, err
	}

	taskResponse := models.ToTaskResponse(*task)

	return &taskResponse, nil
}

func (s *TaskService) GetTasksByStatus(ctx context.Context, status models.Status) []models.TaskResponse {
	if status != "" {
		tasks := s.taskRepository.GetTasksByStatus(ctx, status)
		return models.ToTaskResponses(tasks)
	}

	return models.ToTaskResponses(s.taskRepository.GetAllTasks(ctx))
}
