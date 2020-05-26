// +build debug

// Code generated by MockGen. DO NOT EDIT.
// Source: flash_bag.go

package cookie

import (
	reflect "reflect"

	"github.com/cjtoolkit/ctx/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockFlashBag is a mock of FlashBag interface
type MockFlashBag struct {
	ctrl     *gomock.Controller
	recorder *MockFlashBagMockRecorder
}

// MockFlashBagMockRecorder is the mock recorder for MockFlashBag
type MockFlashBagMockRecorder struct {
	mock *MockFlashBag
}

// NewMockFlashBag creates a new mock instance
func NewMockFlashBag(ctrl *gomock.Controller) *MockFlashBag {
	mock := &MockFlashBag{ctrl: ctrl}
	mock.recorder = &MockFlashBagMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFlashBag) EXPECT() *MockFlashBagMockRecorder {
	return m.recorder
}

// GetFlashBag mocks base method
func (m *MockFlashBag) GetFlashBag(context ctx.Context) FlashBagValues {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlashBag", context)
	ret0, _ := ret[0].(FlashBagValues)
	return ret0
}

// GetFlashBag indicates an expected call of GetFlashBag
func (mr *MockFlashBagMockRecorder) GetFlashBag(context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlashBag", reflect.TypeOf((*MockFlashBag)(nil).GetFlashBag), context)
}

// SaveFlashBagToSession mocks base method
func (m *MockFlashBag) SaveFlashBagToSession(context ctx.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveFlashBagToSession", context)
}

// SaveFlashBagToSession indicates an expected call of SaveFlashBagToSession
func (mr *MockFlashBagMockRecorder) SaveFlashBagToSession(context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFlashBagToSession", reflect.TypeOf((*MockFlashBag)(nil).SaveFlashBagToSession), context)
}
