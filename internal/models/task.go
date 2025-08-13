package models

import (
	"fmt"
	"time"
)

type Task struct {
	Id          string
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type TaskResponse struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func ValidateCreateTaskRequest(createTaskRequest CreateTaskRequest) error {
	if createTaskRequest.Title == "" {
		return fmt.Errorf("title is required")
	}
	if createTaskRequest.Description == "" {
		return fmt.Errorf("description is required")
	}
	if createTaskRequest.Status == "" {
		return fmt.Errorf("status is required")
	}

	return nil
}

func ToTask(createTaskRequest CreateTaskRequest, id string, createdAt time.Time) Task {
	return Task{
		Id:          id,
		Title:       createTaskRequest.Title,
		Description: createTaskRequest.Description,
		Status:      createTaskRequest.Status,
		CreatedAt:   createdAt,
	}
}

func ToTaskResponse(task Task) TaskResponse {
	return TaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
	}
}

func ToTaskResponses(tasks []Task) []TaskResponse {
	result := make([]TaskResponse, 0, len(tasks))

	for _, task := range tasks {
		result = append(result, ToTaskResponse(task))
	}

	return result
}
