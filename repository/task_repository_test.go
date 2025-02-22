package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/sing3demons/go-backend-clean-architecture/mongo/mocks"
	"github.com/sing3demons/go-backend-clean-architecture/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const title = "test title"

func TestTaskRepositoryCreate(t *testing.T) {
	databaseHelper := &mocks.Database{}
	collectionHelper := &mocks.Collection{}

	collectionName := domain.CollectionTask

	mockTask := &domain.Task{
		ID:    primitive.NewObjectID(),
		Title: title,
	}

	mockEmptyTask := &domain.Task{}
	mockTaskTD := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(mockTaskTD, nil).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		err := repo.Create(nil, mockTask)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(mockTaskTD, assert.AnError).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		err := repo.Create(nil, mockEmptyTask)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})
}

type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

func TestTaskRepositoryFetchByUserID(t *testing.T) {
	databaseHelper := &mocks.Database{}
	collectionHelper := &mocks.Collection{}

	collectionName := domain.CollectionTask

	mockTask := &domain.Task{
		ID:    primitive.NewObjectID(),
		Title: title,
	}

	taskMode := "domain.Task"

	t.Run("success", func(t *testing.T) {
		collectionHelper.On("All", mock.Anything, mock.AnythingOfType("*[]domain.Task")).Return(nil).Once()

		tasks := []domain.Task{
			{
				ID:    primitive.NewObjectID(),
				Title: title,
			},
		}

		var documents []any

		for _, task := range tasks {
			documents = append(documents, task)
		}

		mockCursor, err := mongo.NewCursorFromDocuments(documents, nil, nil)
		collectionHelper.On("Find", mock.Anything, mock.AnythingOfType(taskMode)).Return(mockCursor, err).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err = repo.FetchByUserID(context.TODO(), mockTask.UserID.Hex())

		assert.NoError(t, err)
	})

	t.Run("error primitive.ObjectIDFromHex", func(t *testing.T) {
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByUserID(context.TODO(), "invalid")

		assert.Error(t, err)
	})

	t.Run("error collection.Find", func(t *testing.T) {
		collectionHelper.On("Find", mock.Anything, mock.AnythingOfType(taskMode)).Return(nil, assert.AnError).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByUserID(context.TODO(), mockTask.UserID.Hex())

		assert.Error(t, err)
	})

	t.Run("error cursor.All", func(t *testing.T) {
		msg := errors.New("error cursor.All")
		ctx := context.TODO()

		// Create a mock cursor
		mockCursor := &mocks.Cursor{}
		mockCursor.On("All", ctx, mock.Anything).Return(msg)
		mockCursor.On("Close", ctx).Return(nil)

		collectionHelper.On("Find", mock.Anything, mock.AnythingOfType(taskMode)).Return(mockCursor, nil)
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByUserID(ctx, mockTask.UserID.Hex())

		assert.Error(t, err)
		assert.Equal(t, msg, err)
	})
}

func TestTaskRepositoryFetchByTaskID(t *testing.T) {
	databaseHelper := &mocks.Database{}
	collectionHelper := &mocks.Collection{}

	collectionName := domain.CollectionTask

	mockTask := &domain.Task{
		ID:     primitive.NewObjectID(),
		Title:  title,
		UserID: primitive.NewObjectID(),
	}

	taskMode := "domain.Task"

	t.Run("success", func(t *testing.T) {
		mockSingleResult := mongo.NewSingleResultFromDocument(bson.M{}, nil, nil)
		collectionHelper.On("FindOne", mock.Anything, mock.AnythingOfType(taskMode)).Return(mockSingleResult).Once()

		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByTaskID(context.TODO(), mockTask.ID.Hex())

		assert.NoError(t, err)
	})

	t.Run("error primitive.ObjectIDFromHex", func(t *testing.T) {
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByTaskID(context.TODO(), "invalid")

		assert.Error(t, err)
	})

	t.Run("error collection.FindOne", func(t *testing.T) {
		mockSingleResult := &mocks.SingleResult{}
		mockSingleResult.On("Decode", mock.AnythingOfType("*domain.Task")).Return(errors.New("error")).Once()
		collectionHelper.On("FindOne", mock.Anything, mock.AnythingOfType(taskMode)).Return(mockSingleResult).Once()
		databaseHelper.On("Collection", collectionName).Return(collectionHelper).Once()

		repo := repository.NewTaskRepository(databaseHelper, collectionName)
		_, err := repo.FetchByTaskID(context.TODO(), mockTask.ID.Hex())

		assert.Error(t, err)
	})

}
