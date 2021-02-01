// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// TestDirReader is an autogenerated mock type for the dirReader type
type TestDirReader struct {
	mock.Mock
}

// ReadDir provides a mock function with given fields: path
func (_m *TestDirReader) ReadDir(path string) ([]os.FileInfo, error) {
	ret := _m.Called(path)

	var r0 []os.FileInfo
	if rf, ok := ret.Get(0).(func(string) []os.FileInfo); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]os.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
