package http

import (
	"fmt"
	"liveoncetechtask/internal/task"
	"net/http"
)

func RegisterHTTPEndpoints(mux *http.ServeMux, taskService task.Service) {
	h := NewTaskHandler(taskService)

	mux.HandleFunc(fmt.Sprintf("GET %s", TASKS), h.GetTasks)
	mux.HandleFunc(fmt.Sprintf("GET %s", TASK_BY_ID), h.GetTaskById)
	mux.HandleFunc(fmt.Sprintf("POST %s", TASKS), h.CreateTask)
}
