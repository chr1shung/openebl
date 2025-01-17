// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/bu_server/business_unit/bu_storage.go

// Package mock_business_unit is a generated GoMock package.
package mock_business_unit

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	business_unit "github.com/openebl/openebl/pkg/bu_server/business_unit"
	model "github.com/openebl/openebl/pkg/bu_server/model"
	storage "github.com/openebl/openebl/pkg/bu_server/storage"
)

// MockBusinessUnitStorage is a mock of BusinessUnitStorage interface.
type MockBusinessUnitStorage struct {
	ctrl     *gomock.Controller
	recorder *MockBusinessUnitStorageMockRecorder
}

// MockBusinessUnitStorageMockRecorder is the mock recorder for MockBusinessUnitStorage.
type MockBusinessUnitStorageMockRecorder struct {
	mock *MockBusinessUnitStorage
}

// NewMockBusinessUnitStorage creates a new mock instance.
func NewMockBusinessUnitStorage(ctrl *gomock.Controller) *MockBusinessUnitStorage {
	mock := &MockBusinessUnitStorage{ctrl: ctrl}
	mock.recorder = &MockBusinessUnitStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBusinessUnitStorage) EXPECT() *MockBusinessUnitStorageMockRecorder {
	return m.recorder
}

// CreateTx mocks base method.
func (m *MockBusinessUnitStorage) CreateTx(ctx context.Context, options ...storage.CreateTxOption) (storage.Tx, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateTx", varargs...)
	ret0, _ := ret[0].(storage.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockBusinessUnitStorageMockRecorder) CreateTx(ctx interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockBusinessUnitStorage)(nil).CreateTx), varargs...)
}

// ListAuthentication mocks base method.
func (m *MockBusinessUnitStorage) ListAuthentication(ctx context.Context, tx storage.Tx, req business_unit.ListAuthenticationRequest) (business_unit.ListAuthenticationResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAuthentication", ctx, tx, req)
	ret0, _ := ret[0].(business_unit.ListAuthenticationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAuthentication indicates an expected call of ListAuthentication.
func (mr *MockBusinessUnitStorageMockRecorder) ListAuthentication(ctx, tx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAuthentication", reflect.TypeOf((*MockBusinessUnitStorage)(nil).ListAuthentication), ctx, tx, req)
}

// ListBusinessUnits mocks base method.
func (m *MockBusinessUnitStorage) ListBusinessUnits(ctx context.Context, tx storage.Tx, req business_unit.ListBusinessUnitsRequest) (business_unit.ListBusinessUnitsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBusinessUnits", ctx, tx, req)
	ret0, _ := ret[0].(business_unit.ListBusinessUnitsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBusinessUnits indicates an expected call of ListBusinessUnits.
func (mr *MockBusinessUnitStorageMockRecorder) ListBusinessUnits(ctx, tx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBusinessUnits", reflect.TypeOf((*MockBusinessUnitStorage)(nil).ListBusinessUnits), ctx, tx, req)
}

// StoreAuthentication mocks base method.
func (m *MockBusinessUnitStorage) StoreAuthentication(ctx context.Context, tx storage.Tx, auth model.BusinessUnitAuthentication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreAuthentication", ctx, tx, auth)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreAuthentication indicates an expected call of StoreAuthentication.
func (mr *MockBusinessUnitStorageMockRecorder) StoreAuthentication(ctx, tx, auth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreAuthentication", reflect.TypeOf((*MockBusinessUnitStorage)(nil).StoreAuthentication), ctx, tx, auth)
}

// StoreBusinessUnit mocks base method.
func (m *MockBusinessUnitStorage) StoreBusinessUnit(ctx context.Context, tx storage.Tx, bu model.BusinessUnit) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBusinessUnit", ctx, tx, bu)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBusinessUnit indicates an expected call of StoreBusinessUnit.
func (mr *MockBusinessUnitStorageMockRecorder) StoreBusinessUnit(ctx, tx, bu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBusinessUnit", reflect.TypeOf((*MockBusinessUnitStorage)(nil).StoreBusinessUnit), ctx, tx, bu)
}
