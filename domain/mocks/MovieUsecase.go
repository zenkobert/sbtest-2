// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	model "github.com/zenkobert/sbtest-2/domain"
)

// MovieUsecase is an autogenerated mock type for the MovieUsecase type
type MovieUsecase struct {
	mock.Mock
}

// GetMovieDetailByID provides a mock function with given fields: id
func (_m *MovieUsecase) GetMovieDetailByID(id string) (*model.MovieDetail, error) {
	ret := _m.Called(id)

	var r0 *model.MovieDetail
	if rf, ok := ret.Get(0).(func(string) *model.MovieDetail); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MovieDetail)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchMovies provides a mock function with given fields: title, page
func (_m *MovieUsecase) SearchMovies(title string, page uint32) (*model.MovieSearch, error) {
	ret := _m.Called(title, page)

	var r0 *model.MovieSearch
	if rf, ok := ret.Get(0).(func(string, uint32) *model.MovieSearch); ok {
		r0 = rf(title, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MovieSearch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint32) error); ok {
		r1 = rf(title, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}