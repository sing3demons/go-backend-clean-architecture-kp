package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/sing3demons/go-backend-clean-architecture/domain/mocks"
	"github.com/sing3demons/go-backend-clean-architecture/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskHandler(t *testing.T) {
	t.Run("Create Task", func(t *testing.T) {
		timeout := time.Duration(2) * time.Second
		repo := mocks.NewTaskRepository()

		id, _ := primitive.ObjectIDFromHex("67b998e4d5b0121df1966470")

		body := domain.Task{
			ID:    id,
			Title: "title",
		}

		repo.Create(context.TODO(), &body)
		service := usecase.NewTaskUsecase(repo, timeout)

		handler := NewTaskHandler(service)

		c := bootstrap.NewMockMuxContext(http.MethodPost, "/test?query=value", body)

		if err := handler.CreateTask(c); err != nil {
			t.Error("Error")
		}

		actual := domain.Task{}
		err := c.Body(&actual)
		assert.NoError(t, err)
		assert.Equal(t, body, actual)
		assert.Equal(t, 200, c.Res.Code)
	})

	t.Run("ReadInput Task Fail", func(t *testing.T) {
		timeout := time.Duration(2) * time.Second
		repo := mocks.NewTaskRepository()

		service := usecase.NewTaskUsecase(repo, timeout)

		handler := NewTaskHandler(service)

		c := bootstrap.NewMockMuxContext(http.MethodPost, "/test?query=value", nil)

		if err := handler.CreateTask(c); err != nil {
			t.Error("Error")
		}
		actual := strings.TrimSpace(c.Res.Body.String())
		expected := "\"EOF\""
		assert.Equal(t, 400, c.Res.Code)
		assert.Equal(t, expected, actual)
	})

	t.Run("Create Task Fail", func(t *testing.T) {
		expectedError := errors.New("failed to create task")
		id, _ := primitive.ObjectIDFromHex("67b998e4d5b0121df1966470")

		modelsTask := domain.Task{
			ID:    id,
			Title: "title",
		}

		service := new(usecase.MockTaskUsecase)
		service.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(expectedError)

		handler := NewTaskHandler(service)

		c := bootstrap.NewMockMuxContext(http.MethodPost, "/test", modelsTask)

		if err := handler.CreateTask(c); err != nil {
			t.Error("Error")
		}
		actual := strings.TrimSpace(c.Res.Body.String())
		expected := "\"failed to create task\""
		assert.Equal(t, 500, c.Res.Code)
		assert.Equal(t, expected, actual)
	})

	t.Run("Get Task", func(t *testing.T) {
		timeout := time.Duration(2) * time.Second
		repo := mocks.NewTaskRepository()

		id, _ := primitive.ObjectIDFromHex("67b998e4d5b0121df1966470")

		body := domain.Task{
			ID:    id,
			Title: "title",
		}

		repo.Create(context.TODO(), &body)
		service := usecase.NewTaskUsecase(repo, timeout)

		handler := NewTaskHandler(service)

		c := bootstrap.NewMockMuxContext(http.MethodGet, "/test", nil)

		if err := handler.GetTask(c); err != nil {
			t.Error("Error")
		}

		actual := []domain.Task{}
		err := c.Body(&actual)
		assert.NoError(t, err)
		assert.Equal(t, 200, c.Res.Code)
	})

	t.Run("Get Task Fail", func(t *testing.T) {
		expectedError := errors.New("failed to fetch task")

		service := new(usecase.MockTaskUsecase)
		service.On("FetchAll", mock.Anything).Return([]domain.Task{}, expectedError)

		handler := NewTaskHandler(service)

		c := bootstrap.NewMockMuxContext(http.MethodGet, "/test", nil)

		if err := handler.GetTask(c); err != nil {
			t.Error("Error")
		}
		actual := strings.TrimSpace(c.Res.Body.String())
		expected := "\"failed to fetch task\""
		assert.Equal(t, 500, c.Res.Code)
		assert.Equal(t, expected, actual)
	})

}
