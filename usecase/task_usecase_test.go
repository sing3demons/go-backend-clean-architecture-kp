package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/sing3demons/go-backend-clean-architecture/repository"
	"github.com/sing3demons/go-backend-clean-architecture/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFetchByUserId(t *testing.T) {
	mockTaskRepository := new(repository.MockTaskRepository)
	userObjectID := primitive.NewObjectID()
	userID := userObjectID.Hex()

	t.Run("success", func(t *testing.T) {

		mockTask := domain.Task{
			ID:     primitive.NewObjectID(),
			Title:  "Test_Title",
			UserID: userObjectID,
		}

		mockListTask := make([]domain.Task, 0)
		mockListTask = append(mockListTask, mockTask)

		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(mockListTask, nil).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		list, err := u.FetchByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.Len(t, list, len(mockListTask))

		mockTaskRepository.AssertExpectations(t)
	})
	t.Run("error", func(t *testing.T) {
		mockTaskRepository.On("FetchByUserID", mock.Anything, userID).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		list, err := u.FetchByUserID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, list)

		mockTaskRepository.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	mockTaskRepository := new(repository.MockTaskRepository)

	t.Run("success", func(t *testing.T) {
		mockTask := domain.Task{
			ID:     primitive.NewObjectID(),
			Title:  "Test Title",
			UserID: primitive.NewObjectID(),
		}

		mockTaskRepository.On("Create", mock.Anything, &mockTask).Return(nil).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		err := u.Create(context.Background(), &mockTask)

		assert.NoError(t, err)

		mockTaskRepository.AssertExpectations(t)
	})
	t.Run("error", func(t *testing.T) {
		mockTask := domain.Task{
			ID:     primitive.NewObjectID(),
			Title:  "Test Title",
			UserID: primitive.NewObjectID(),
		}

		mockTaskRepository.On("Create", mock.Anything, &mockTask).Return(errors.New("Unexpected")).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		err := u.Create(context.Background(), &mockTask)

		assert.Error(t, err)

		mockTaskRepository.AssertExpectations(t)
	})
}

func TestFetchByTaskID(t *testing.T) {
	mockTaskRepository := new(repository.MockTaskRepository)
	taskObjectID := primitive.NewObjectID()
	taskID := taskObjectID.Hex()

	t.Run("success", func(t *testing.T) {

		mockTask := domain.Task{
			ID:     taskObjectID,
			Title:  "Test_Title",
			UserID: primitive.NewObjectID(),
		}

		mockTaskRepository.On("FetchByTaskID", mock.Anything, taskID).Return(mockTask, nil).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		task, err := u.FetchByTaskID(context.Background(), taskID)

		assert.NoError(t, err)
		assert.NotNil(t, task)

		mockTaskRepository.AssertExpectations(t)
	})
	t.Run("error", func(t *testing.T) {
		mockTaskRepository.On("FetchByTaskID", mock.Anything, taskID).Return(domain.Task{}, errors.New("Unexpected")).Once()

		u := usecase.NewTaskUsecase(mockTaskRepository, time.Second*2)

		task, err := u.FetchByTaskID(context.Background(), taskID)

		assert.Error(t, err)
		assert.Equal(t, domain.Task{}, task)

		mockTaskRepository.AssertExpectations(t)
	})
}
