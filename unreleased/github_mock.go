// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.go

package unreleased

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of FetcherInterface interface
type MockFetcherInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockFetcherInterfaceRecorder
}

// Recorder for MockFetcherInterface (not exported)
type _MockFetcherInterfaceRecorder struct {
	mock *MockFetcherInterface
}

func NewMockFetcherInterface(ctrl *gomock.Controller) *MockFetcherInterface {
	mock := &MockFetcherInterface{ctrl: ctrl}
	mock.recorder = &_MockFetcherInterfaceRecorder{mock}
	return mock
}

func (_m *MockFetcherInterface) EXPECT() *_MockFetcherInterfaceRecorder {
	return _m.recorder
}

func (_m *MockFetcherInterface) FetchTags(slug string) (Tags, error) {
	ret := _m.ctrl.Call(_m, "FetchTags", slug)
	ret0, _ := ret[0].(Tags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFetcherInterfaceRecorder) FetchTags(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FetchTags", arg0)
}

func (_m *MockFetcherInterface) FetchRepos(affiliation string, perPage int) (Repositories, error) {
	ret := _m.ctrl.Call(_m, "FetchRepos", affiliation, perPage)
	ret0, _ := ret[0].(Repositories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFetcherInterfaceRecorder) FetchRepos(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FetchRepos", arg0, arg1)
}

func (_m *MockFetcherInterface) CompareCommits(slug string, base string, head string) ([]Commit, error) {
	ret := _m.ctrl.Call(_m, "CompareCommits", slug, base, head)
	ret0, _ := ret[0].([]Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFetcherInterfaceRecorder) CompareCommits(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CompareCommits", arg0, arg1, arg2)
}
