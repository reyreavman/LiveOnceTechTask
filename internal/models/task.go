package models

import "time"

type Task struct {
	Id          string
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
}

type CreateTaskRequest struct {
	Title       string
	Description string
	Status      Status
}

type TaskResponse struct {
	Id          string
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
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
