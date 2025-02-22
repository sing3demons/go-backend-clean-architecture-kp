package mocks

import (
	"context"

	"github.com/sing3demons/go-backend-clean-architecture/domain"
	mock "github.com/stretchr/testify/mock"
)

type TaskRepository struct {
	mock.Mock
}

func (_m *TaskRepository) Create(c context.Context, task *domain.Task) error {
	ret := _m.Called(c, task)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Task) error); ok {
		r0 = rf(c, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *TaskRepository) FetchAll(c context.Context) ([]domain.Task, error) {
	ret := _m.Called(c)

	var r0 []domain.Task
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Task); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)

	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TaskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	ret := _m.Called(c, userID)

	var r0 []domain.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.Task); ok {
		r0 = rf(c, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TaskRepository) FetchByTaskID(c context.Context, taskID string) (domain.Task, error) {
	ret := _m.Called(c, taskID)

	var r0 domain.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Task); ok {
		r0 = rf(c, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTaskRepository interface {
	mock.TestingT
	Cleanup(func())
}

func NewTaskRepository() *TaskRepository {
	m := &TaskRepository{}
	m.On("Create", mock.Anything, mock.Anything).Return(nil)
	m.On("FetchAll", mock.Anything).Return([]domain.Task{}, nil)
	m.On("FetchByUserID", mock.Anything, mock.Anything).Return([]domain.Task{}, nil)
	m.On("FetchByTaskID", mock.Anything, mock.Anything).Return(domain.Task{}, nil)

	// mock.Mock.Test(t)

	// t.Cleanup(func() { mock.AssertExpectations(t) })

	return m
}
