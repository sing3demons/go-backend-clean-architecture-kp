package handler

import (
	"fmt"

	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/domain"
)

type TaskHandler struct {
	TaskService domain.TaskUsecase
}

func NewTaskHandler(taskService domain.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
	}
}

func (handler *TaskHandler) CreateTask(ctx bootstrap.IContext) error {
	fmt.Println("Create Task")
	var task domain.Task
	if err := ctx.ReadInput(&task); err != nil {
		return ctx.Response(400, err.Error())
	}

	if err := handler.TaskService.Create(ctx.Context(), &task); err != nil {
		return ctx.Response(500, err.Error())
	}

	return ctx.Response(200, task)
}

func (h *TaskHandler) GetTask(ctx bootstrap.IContext) error {
	tasks, err := h.TaskService.FetchAll(ctx.Context())

	if err != nil {
		return ctx.Response(500, err.Error())
	}

	return ctx.Response(200, tasks)
}
