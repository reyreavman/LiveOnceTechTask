package service

import (
	"context"
	"liveoncetechtask/internal/logger"
	"liveoncetechtask/internal/models"
	"liveoncetechtask/internal/task"
	"liveoncetechtask/pkg/id"
	"time"
)

type TaskService struct {
	taskRepository task.Repository
	idGenerator    id.Generator
	logger         *logger.Logger
}

func NewTaskService(taskRepository task.Repository, logger *logger.Logger) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
		logger:         logger,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, createTaskRequest models.CreateTaskRequest) models.TaskResponse {
	s.logger.Debug("task_service: CreateTask", map[string]interface{}{
		"title":       createTaskRequest.Title,
		"description": createTaskRequest.Description,
	})

	taskToSave := models.ToTask(createTaskRequest, s.idGenerator.Generate(), time.Now())
	savedTask := s.taskRepository.CreateTask(ctx, taskToSave)

	return models.ToTaskResponse(savedTask)
}

func (s *TaskService) GetTaskById(ctx context.Context, id string) (*models.TaskResponse, error) {
	s.logger.Debug("Service: GetTaskById", map[string]interface{}{
		"task_id": id,
	})

	task, err := s.taskRepository.GetTaskById(ctx, id)
	if err != nil {
		return nil, err
	}

	taskResponse := models.ToTaskResponse(*task)

	return &taskResponse, nil
}

func (s *TaskService) GetTasksByStatus(ctx context.Context, status models.Status) []models.TaskResponse {
	s.logger.Debug("Service: GetTasksByStatus", map[string]interface{}{
		"status": status,
	})

	if status != "" {
		tasks := s.taskRepository.GetTasksByStatus(ctx, status)
		return models.ToTaskResponses(tasks)
	}

	return models.ToTaskResponses(s.taskRepository.GetAllTasks(ctx))
}
