package usecase

import (
	"context"
	"time"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (u *taskUsecase) Create(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.taskRepository.Create(ctx, task)
}

func (u *taskUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.taskRepository.FetchByUserID(ctx, userID)
}

func (u *taskUsecase) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.taskRepository.FetchByTaskID(ctx, taskID)
}

func (u *taskUsecase) FetchAll(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.taskRepository.FetchAll(ctx)
}
