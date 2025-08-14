package server

import (
	"context"
	thttp "liveoncetechtask/internal/delivery/http"
	"liveoncetechtask/internal/logger"
	ratelimiter "liveoncetechtask/internal/rate_limiter"
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
	logger  *logger.Logger
}

func NewApp() *App {
	log := logger.NewLogger()
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository)

	return &App{
		service: taskService,
		logger:  log,
	}
}

func (a *App) Run(port string) error {
	rl := ratelimiter.NewRateLimiter(100, time.Second)
	defer rl.Stop()

	router := http.NewServeMux()

	thttp.RegisterHTTPEndpoints(router, a.service, a.logger, rl)

	a.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		a.logger.Info("Server started", "", map[string]interface{}{
			"port": 8080,
		})
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
