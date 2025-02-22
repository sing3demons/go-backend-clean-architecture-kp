package route

import (
	"time"

	"github.com/sing3demons/go-backend-clean-architecture/api/handler"
	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/mongo"
	"github.com/sing3demons/go-backend-clean-architecture/repository"
	"github.com/sing3demons/go-backend-clean-architecture/usecase"
)

func NewTaskRoute(db mongo.Database, collection string, router bootstrap.IApplication) {
	timeout := time.Duration(2) * time.Second
	repo := repository.NewTaskRepository(db, collection)
	service := usecase.NewTaskUsecase(repo, timeout)
	handler := handler.NewTaskHandler(service)

	router.Get("/task", handler.GetTask)

	router.Post("/task", handler.CreateTask)
}
