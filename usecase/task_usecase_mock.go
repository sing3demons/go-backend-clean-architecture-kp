package usecase

import (
	"context"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/stretchr/testify/mock"
)

// MockTaskUsecase is a mock for the TaskUsecase interface
type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) Create(c context.Context, task *domain.Task) error {
	args := m.Called(c, task)
	return args.Error(0)
}

func (m *MockTaskUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	args := m.Called(c, userID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	args := m.Called(c, taskID)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) FetchAll(c context.Context) ([]domain.Task, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.Task), args.Error(1)
}
