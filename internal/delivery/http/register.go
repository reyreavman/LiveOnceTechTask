package http

import (
	"fmt"
	"liveoncetechtask/internal/logger"
	ratelimiter "liveoncetechtask/internal/rate_limiter"
	"liveoncetechtask/internal/task"
	"net/http"
)

func RegisterHTTPEndpoints(mux *http.ServeMux, taskService task.Service, log *logger.Logger, rl *ratelimiter.RateLimiter) {
	h := NewTaskHandler(taskService, log)

	wrapHandler := func(handler http.HandlerFunc) http.Handler {
		var h http.Handler = handler
		h = logger.Middleware(log, h)
		h = ratelimiter.Middleware(rl, h)
		return h
	}

	mux.Handle(fmt.Sprintf("GET %s", TASKS), wrapHandler(h.GetTasks))
	mux.Handle(fmt.Sprintf("GET %s", TASK_BY_ID), wrapHandler(h.GetTaskById))
	mux.Handle(fmt.Sprintf("POST %s", TASKS), wrapHandler(h.CreateTask))
	mux.Handle(fmt.Sprintf("GET %s", TASK_STATUS_LIST), wrapHandler(h.GetStatusList))
}
