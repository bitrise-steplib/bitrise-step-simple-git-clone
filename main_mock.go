package main

import "github.com/stretchr/testify/mock"

type mockGitCommand struct {
	mock.Mock
}

func (_m *mockGitCommand) addRemote(name string, url string) error {
	ret := _m.Called(name, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *mockGitCommand) checkout(arg string) error {
	ret := _m.Called(arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *mockGitCommand) fetchWithRetry(opts ...string) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...string) error); ok {
		r0 = rf(opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *mockGitCommand) init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *mockGitCommand) merge(arg string) error {
	ret := _m.Called(arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConfigParser struct {
	mock.Mock
}

func (_m *mockConfigParser) parse(conf *config) error {
	ret := _m.Called(conf)

	var r0 error
	if rf, ok := ret.Get(0).(func(*config) error); ok {
		r0 = rf(conf)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockGitCommandFactory struct {
	mock.Mock
}

func (_m *mockGitCommandFactory) new(dir string) (gitCommand, error) {
	ret := _m.Called(dir)

	var r0 gitCommand
	if rf, ok := ret.Get(0).(func(string) gitCommand); ok {
		r0 = rf(dir)
	} else {
		if ret.Get(0) != nil {
			r0, ok = ret.Get(0).(gitCommand)
			if ok {
			}
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
