package repository

import (
	"context"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/sing3demons/go-backend-clean-architecture/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (r *taskRepository) Create(c context.Context, task *domain.Task) error {
	col := r.database.Collection(r.collection)

	task.ID = primitive.NewObjectID()

	_, err := col.InsertOne(c, task)

	return err
}

func (r *taskRepository) FetchAll(c context.Context) ([]domain.Task, error) {
	col := r.database.Collection(r.collection)
	tasks := []domain.Task{}

	cursor, err := col.Find(context.Background(), bson.M{})
	if err != nil {
		return tasks, err
	}

	defer cursor.Close(c)

	if err := cursor.All(c, &tasks); err != nil {
		return nil, err
	}

	return tasks, err
}

func (r *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	col := r.database.Collection(r.collection)
	tasks := []domain.Task{}
	idHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return tasks, err
	}

	cursor, err := col.Find(c, domain.Task{UserID: idHex})
	if err != nil {
		return tasks, err
	}

	defer cursor.Close(c)

	if err := cursor.All(c, &tasks); err != nil {
		return nil, err
	}

	return tasks, err
}

func (r *taskRepository) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	col := r.database.Collection(r.collection)
	task := domain.Task{}
	idHex, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return task, err
	}

	result := col.FindOne(c, domain.Task{ID: idHex})
	if err := result.Decode(&task); err != nil {
		return task, err
	}
	return task, nil
}
