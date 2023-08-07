// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mimzeslami/expense_share/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "github.com/mimzeslami/expense_share/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateFellowTravelers mocks base method.
func (m *MockStore) CreateFellowTravelers(arg0 context.Context, arg1 db.CreateFellowTravelersParams) (db.FellowTravelers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFellowTravelers", arg0, arg1)
	ret0, _ := ret[0].(db.FellowTravelers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFellowTravelers indicates an expected call of CreateFellowTravelers.
func (mr *MockStoreMockRecorder) CreateFellowTravelers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFellowTravelers", reflect.TypeOf((*MockStore)(nil).CreateFellowTravelers), arg0, arg1)
}

// CreateTrip mocks base method.
func (m *MockStore) CreateTrip(arg0 context.Context, arg1 db.CreateTripParams) (db.Trips, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrip", arg0, arg1)
	ret0, _ := ret[0].(db.Trips)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTrip indicates an expected call of CreateTrip.
func (mr *MockStoreMockRecorder) CreateTrip(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrip", reflect.TypeOf((*MockStore)(nil).CreateTrip), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// GetFellowTraveler mocks base method.
func (m *MockStore) GetFellowTraveler(arg0 context.Context, arg1 int64) (db.FellowTravelers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFellowTraveler", arg0, arg1)
	ret0, _ := ret[0].(db.FellowTravelers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFellowTraveler indicates an expected call of GetFellowTraveler.
func (mr *MockStoreMockRecorder) GetFellowTraveler(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFellowTraveler", reflect.TypeOf((*MockStore)(nil).GetFellowTraveler), arg0, arg1)
}

// GetTrip mocks base method.
func (m *MockStore) GetTrip(arg0 context.Context, arg1 int64) (db.Trips, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrip", arg0, arg1)
	ret0, _ := ret[0].(db.Trips)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrip indicates an expected call of GetTrip.
func (mr *MockStoreMockRecorder) GetTrip(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrip", reflect.TypeOf((*MockStore)(nil).GetTrip), arg0, arg1)
}

// GetTripFellowTravelers mocks base method.
func (m *MockStore) GetTripFellowTravelers(arg0 context.Context, arg1 int64) ([]db.FellowTravelers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTripFellowTravelers", arg0, arg1)
	ret0, _ := ret[0].([]db.FellowTravelers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTripFellowTravelers indicates an expected call of GetTripFellowTravelers.
func (mr *MockStoreMockRecorder) GetTripFellowTravelers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTripFellowTravelers", reflect.TypeOf((*MockStore)(nil).GetTripFellowTravelers), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 string) (db.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}
