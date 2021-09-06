package repository

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	"github.com/zenkobert/sbtest-2/common/mocks"
	model "github.com/zenkobert/sbtest-2/domain"
)

var (
	searchMovieJsonResponse = `{
		"Search": [
			{
				"Title": "title",
				"Year": "year",
				"imdbID": "imdbid",
				"Type": "type",
				"Poster": "poster"
			}
		],
		"totalResults": "1",
		"Response": "True"
	}`

	searchMovieInvalidJsonResponse = `{
		"Search": [
			{
				"Title": "title",
				"Year": "year",
				"imdbID": "imdbid",
				"Type": "type",
				"Poster": ["poster"]
			}
		],
		"totalResults": "1",
		"Response": "True"
	}`

	errorJsonResponse = `{
		"Response": "False",
		"Error": "Invalid API Key"
	}`

	getMovieDetailJsonResponse = `
		{
			"Title": "title",
			"Year": "year",
			"Rated": "rated",
			"Ratings": [
				{
					"Source": "source",
					"Value": "value"
				}
			]
		}
	`

	getMovieDetailInvalidJsonResponse = `
		{
			"Title": "tile",
			"Year": "year",
			"Rated": "rated",
			"Ratings": "ratings"
		}
	`
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (errReader) Close() (err error) {
	return errors.New("close error")
}

func TestNewMovieRepo(t *testing.T) {
	t.Run("[NewMovieRepo]", func(t *testing.T) {
		apiKey := "randomKey"

		expected := &movieRepo{
			Client: &http.Client{},
			apiKey: apiKey,
		}

		actual := NewMovieRepo(apiKey)
		assert.Equal(t, expected, actual)
	})
}

func TestSearchMovie(t *testing.T) {
	t.Run("[SearchMovies] error response", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		httpClientMock.On("Do", testify.Anything).Return(&http.Response{}, errors.New("error"))

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.SearchMovies("ironman", 1)
		assert.Error(t, err)
	})

	t.Run("[SearchMovies] read body error", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: errReader(0)}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.SearchMovies("ironman", 1)
		if assert.Error(t, err) {
			assert.Equal(t, "read error", err.Error())
		}
	})

	t.Run("[SearchMovies] no error, positive test", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(searchMovieJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 200}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		expectedResult := &model.MovieSearch{
			Search: []model.SearchDetail{
				{"title", "year", "imdbid", "type", "poster"},
			},
			TotalResults: "1",
			Response:     "True",
		}
		result, err := movieRepo.SearchMovies("ironman", 1)
		if assert.Nil(t, err) {
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("[SearchMovies] response body unmarshall error", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(searchMovieInvalidJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 200}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		result, err := movieRepo.SearchMovies("ironman", 1)
		fmt.Println(result)
		assert.Error(t, err)
	})

	t.Run("[SearchMovies] response code >= 400", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(errorJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 401}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.SearchMovies("ironman", 1)
		if assert.Error(t, err) {
			assert.Equal(t, errors.New("oops, something happened"), err)
		}
	})
}

func TestGetMovieDetailByID(t *testing.T) {
	t.Run("[GetMovieDetailByID] error response", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		httpClientMock.On("Do", testify.Anything).Return(&http.Response{}, errors.New("error"))

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.GetMovieDetailByID("id")
		assert.Error(t, err)
	})

	t.Run("[GetMovieDetailByID] read body error", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: errReader(0)}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.GetMovieDetailByID("id")
		if assert.Error(t, err) {
			assert.Equal(t, "read error", err.Error())
		}
	})

	t.Run("[GetMovieDetailByID] no error, positive test", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(getMovieDetailJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 200}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		expectedResult := &model.MovieDetail{
			Title: "title",
			Year:  "year",
			Rated: "rated",
			Ratings: []model.MovieRating{
				{"source", "value"},
			},
		}
		result, err := movieRepo.GetMovieDetailByID("id")
		if assert.Nil(t, err) {
			assert.Equal(t, expectedResult, result)
		}
	})

	t.Run("[GetMovieDetailByID] response body unmarshall error", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(getMovieDetailInvalidJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 200}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		result, err := movieRepo.GetMovieDetailByID("id")
		fmt.Println(result)
		assert.Error(t, err)
	})

	t.Run("[GetMovieDetailByID] response code >= 400", func(t *testing.T) {
		httpClientMock := &mocks.HTTPClient{}
		dummyBody := ioutil.NopCloser(bytes.NewReader([]byte(errorJsonResponse)))

		httpClientMock.On("Do", testify.Anything).Return(&http.Response{Body: dummyBody, StatusCode: 401}, nil)

		movieRepo := &movieRepo{
			Client: httpClientMock,
			apiKey: "abc",
		}

		_, err := movieRepo.GetMovieDetailByID("id")
		if assert.Error(t, err) {
			assert.Equal(t, errors.New("oops, something happened"), err)
		}
	})
}
