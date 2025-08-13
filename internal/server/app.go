package server

import (
	"context"
	thttp "liveoncetechtask/internal/delivery/http"
	"liveoncetechtask/internal/task"
	"liveoncetechtask/internal/task/repository"
	"liveoncetechtask/internal/task/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	service task.Service
}

func NewApp() *App {
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository)

	return &App{
		service: taskService,
	}
}

func (a *App) Run(port string) error {
	router := http.NewServeMux()

	thttp.RegisterHTTPEndpoints(router, a.service)

	a.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server started")
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to listen and serve: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
