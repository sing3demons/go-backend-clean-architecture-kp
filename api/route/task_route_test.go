package route

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/domain"
	"github.com/sing3demons/go-backend-clean-architecture/mongo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type DB struct {
	collection *mocks.Collection
}

func newMockDatabase() *DB {
	return &DB{}
}

func (d *DB) collectionName() string {
	return "task"
}

func (d *DB) DatabaseSuccess() *mocks.Database {
	database := &mocks.Database{}
	collectionName := d.collectionName()

	database.On("Collection", collectionName).Return(d.collection).Once()
	return database
}

func (d *DB) GetTypeString(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func (d *DB) Find(documents interface{}) *mocks.Database {
	docSlice := reflect.ValueOf(documents)
	var docInterfaces []interface{}
	for i := 0; i < docSlice.Len(); i++ {
		docInterfaces = append(docInterfaces, docSlice.Index(i).Interface())
	}

	d.collection = &mocks.Collection{}

	d.collection.On("All", mock.Anything, mock.AnythingOfType(d.GetTypeString(documents))).Return(nil).Once()

	mockCursor, err := mongo.NewCursorFromDocuments(docInterfaces, nil, nil)
	d.collection.On("Find", mock.Anything, mock.Anything).Return(mockCursor, err).Once()

	return d.DatabaseSuccess()
}

func (d *DB) Create(document interface{}) *mocks.Database {
	d.collection = &mocks.Collection{}
	mockTaskTD := primitive.NewObjectID()
	d.collection.On("InsertOne", mock.Anything, mock.AnythingOfType(d.GetTypeString(document))).Return(mockTaskTD, nil).Once()

	return d.DatabaseSuccess()
}

type MockContext struct {
	path   string
	method string
	body   io.Reader
}

func NewMockContext() *MockContext {
	return &MockContext{
		body: nil,
	}
}

func (m *MockContext) Get(path string) {
	m.path = path
	m.method = http.MethodGet
}

func (m *MockContext) Post(path string, body io.Reader) {
	m.path = path
	m.method = http.MethodPost
	m.body = body
}

func (m MockContext) NewContext() (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(m.method, m.method, m.body)
	rec := httptest.NewRecorder()
	return req, rec
}

func (m *MockContext) Request() *http.Request {
	req, _ := http.NewRequest(m.method, m.path, m.body)
	return req
}

func (m *MockContext) Response() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestTaskApi(t *testing.T) {

	t.Run("Test Task API", func(t *testing.T) {
		server := bootstrap.NewApplication(&bootstrap.Config{
			AppConfig: bootstrap.AppConfig{
				Port: "3000",
			},
		}, bootstrap.NewZapLogger(zap.NewNop()))

		databaseHelper := newMockDatabase()

		id, _ := primitive.ObjectIDFromHex("67b998e4d5b0121df1966470")

		tasks := []domain.Task{
			{
				ID:    id,
				Title: "title",
			},
		}

		db := databaseHelper.Find(tasks)

		router := Setup(db, databaseHelper.collectionName(), server)

		c := NewMockContext()
		c.Get("/task")

		rec := c.Response()

		router.ServeHTTP(rec, c.Request())
		assert.Equal(t, http.StatusOK, rec.Code)

		expected := tasks

		actual := []domain.Task{}
		err := json.Unmarshal(rec.Body.Bytes(), &actual)
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("TASK POST", func(t *testing.T) {
		server := bootstrap.NewApplication(&bootstrap.Config{
			AppConfig: bootstrap.AppConfig{
				Port: "3000",
			},
		}, bootstrap.NewZapLogger(zap.NewNop()))
		databaseHelper := newMockDatabase()

		id, _ := primitive.ObjectIDFromHex("67b998e4d5b0121df1966470")

		task := domain.Task{
			ID:    id,
			Title: "title",
		}

		body := domain.Task{
			Title: "title",
		}
		jsonData, _ := json.Marshal(&body)
		db := databaseHelper.Create(&task)

		router := Setup(db, databaseHelper.collectionName(), server)

		c := NewMockContext()
		c.Post("/task", bytes.NewBuffer(jsonData))

		rec := c.Response()

		router.ServeHTTP(rec, c.Request())
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := domain.Task{}
		err := json.Unmarshal(rec.Body.Bytes(), &actual)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

}
