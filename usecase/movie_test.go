package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	commonMock "github.com/zenkobert/sbtest-2/common/mocks"
	model "github.com/zenkobert/sbtest-2/domain"
	"github.com/zenkobert/sbtest-2/domain/mocks"
)

func TestNewMovieUsecase(t *testing.T) {
	t.Run("[NewMovieUsecase]", func(t *testing.T) {
		expected := &movieUsecase{
			MovieRepo: &mocks.MovieRepository{},
			MovieDB:   &commonMock.DummyDB{},
		}

		actual := NewMovieUsecase(&mocks.MovieRepository{}, &commonMock.DummyDB{})
		assert.Equal(t, expected, actual)
	})
}

func TestSearchMovies(t *testing.T) {
	t.Run("[SearchMovies] movieRepo returns error", func(t *testing.T) {
		movieRepoMock := &mocks.MovieRepository{}
		movieDBMock := &commonMock.DummyDB{}
		movieRepoMock.On("SearchMovies", testify.Anything, testify.Anything).Return(&model.MovieSearch{}, errors.New("error"))
		movieDBMock.On("Log", testify.Anything).Return(nil)

		usecase := NewMovieUsecase(movieRepoMock, movieDBMock)
		_, err := usecase.SearchMovies("test", 1)
		if assert.Error(t, err) {
			assert.Equal(t, "error", err.Error())
		}
	})

	t.Run("[SearchMovies] positive test, no error", func(t *testing.T) {
		expectedResult := &model.MovieSearch{
			Search: []model.SearchDetail{
				{"title", "year", "imdbid", "type", "poster"},
			},
			TotalResults: "1",
			Response:     "True",
		}

		movieRepoMock := &mocks.MovieRepository{}
		movieDBMock := &commonMock.DummyDB{}
		movieRepoMock.On("SearchMovies", testify.Anything, testify.Anything).Return(expectedResult, nil)
		movieDBMock.On("Log", testify.Anything).Return(nil)

		usecase := NewMovieUsecase(movieRepoMock, movieDBMock)
		result, err := usecase.SearchMovies("test", 1)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedResult, result)
		}
	})
}

func TestGetMovieDetailByID(t *testing.T) {
	t.Run("[GetMovieDetailByID] movieRepo returns error", func(t *testing.T) {
		movieRepoMock := &mocks.MovieRepository{}
		movieDBMock := &commonMock.DummyDB{}
		movieRepoMock.On("GetMovieDetailByID", testify.Anything).Return(&model.MovieDetail{}, errors.New("error"))
		movieDBMock.On("Log", testify.Anything).Return(nil)

		usecase := NewMovieUsecase(movieRepoMock, movieDBMock)
		_, err := usecase.GetMovieDetailByID("id")
		if assert.Error(t, err) {
			assert.Equal(t, "error", err.Error())
		}
	})

	t.Run("[GetMovieDetailByID] positive test, no error", func(t *testing.T) {
		expectedResult := &model.MovieDetail{
			Title: "title",
			Year:  "year",
			Rated: "rated",
			Ratings: []model.MovieRating{
				{"source", "value"},
			},
		}

		movieRepoMock := &mocks.MovieRepository{}
		movieDBMock := &commonMock.DummyDB{}
		movieRepoMock.On("GetMovieDetailByID", testify.Anything).Return(expectedResult, nil)
		movieDBMock.On("Log", testify.Anything).Return(nil)

		usecase := NewMovieUsecase(movieRepoMock, movieDBMock)
		result, err := usecase.GetMovieDetailByID("id")
		if assert.Nil(t, err) {
			assert.Equal(t, expectedResult, result)
		}
	})
}

func TestLogToDB(t *testing.T) {
	t.Run("[LogToDB] return error", func(t *testing.T) {
		movieRepoMock := &mocks.MovieRepository{}
		movieDBMock := &commonMock.DummyDB{}
		movieDBMock.On("Log", testify.Anything).Return(errors.New("error"))

		usecase := movieUsecase{movieRepoMock, movieDBMock}
		err := usecase.LogToDB("")
		assert.Error(t, err)
	})
}
