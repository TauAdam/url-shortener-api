// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ShortcutSaver is an autogenerated mock type for the ShortcutSaver type
type ShortcutSaver struct {
	mock.Mock
}

// SaveShortcut provides a mock function with given fields: urlText, alias
func (_m *ShortcutSaver) SaveShortcut(urlText string, alias string) (int64, error) {
	ret := _m.Called(urlText, alias)

	if len(ret) == 0 {
		panic("no return value specified for SaveShortcut")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int64, error)); ok {
		return rf(urlText, alias)
	}
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(urlText, alias)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(urlText, alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShortcutSaver creates a new instance of ShortcutSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShortcutSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShortcutSaver {
	mock := &ShortcutSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
