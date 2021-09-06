package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	model "github.com/zenkobert/sbtest-2/domain"
	mock "github.com/zenkobert/sbtest-2/domain/mocks"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var todoContext = context.TODO()

func TestNewMovieServer(t *testing.T) {
	t.Run("[NewMovieServer]", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}

		expected := &movieServer{
			MovieUsecase: movieUsecaseMock,
		}

		actual := NewMovieServer(movieUsecaseMock)
		assert.Equal(t, expected, actual)
	})
}

func TestSearchMovie(t *testing.T) {
	t.Run("[SearchMovie] ensure page always > 0 and searchword is URL encoded", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("SearchMovies", testify.Anything, testify.Anything).Return(&model.MovieSearch{}, nil)

		serv := &movieServer{movieUsecaseMock}
		req := &SearchMovieRequest{
			Pagination: -1,
			Searchword: "iron man",
		}

		serv.SearchMovie(todoContext, req)
		assert.Equal(t, int32(1), req.Pagination)
		assert.Equal(t, "iron+man", req.Searchword)
	})

	t.Run("[SearchMovie] IF searchword is null, RETURN InvalidArgument error", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("SearchMovies", testify.Anything, testify.Anything).Return(&model.MovieSearch{}, nil)

		serv := &movieServer{movieUsecaseMock}
		req := &SearchMovieRequest{
			Pagination: 1,
			Searchword: "",
		}

		_, err := serv.SearchMovie(todoContext, req)
		if assert.Error(t, err) {
			st, _ := status.FromError(err)
			assert.Equal(t, st.Code(), codes.InvalidArgument)
		}

	})

	t.Run("[SearchMovie] MovieUsecase return an error", func(t *testing.T) {
		errMsg := "Oops, something happened"

		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("SearchMovies", testify.Anything, testify.Anything).Return(&model.MovieSearch{}, errors.New(errMsg))

		serv := &movieServer{movieUsecaseMock}
		req := &SearchMovieRequest{Searchword: "ironman"}

		expectedErr := status.Error(codes.Internal, errMsg)
		_, actualErr := serv.SearchMovie(todoContext, req)
		if assert.Error(t, actualErr) {
			assert.Equal(t, expectedErr, actualErr)
		}
	})

	t.Run("[SearchMovie] movieSearch has an error", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("SearchMovies", testify.Anything, testify.Anything).Return(&model.MovieSearch{Error: "error"}, nil)

		serv := &movieServer{movieUsecaseMock}
		req := &SearchMovieRequest{Searchword: "ironman"}

		_, actualErr := serv.SearchMovie(todoContext, req)
		if assert.Error(t, actualErr) {
			assert.Equal(t, movieNotFoundError, actualErr)
		}
	})

	t.Run("[SearchMovie] no error, movieSearch has result", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieSearchResult := &model.MovieSearch{
			Search: []model.SearchDetail{
				{"Iron Man", "2008", "id1", "movie", "poster1"},
				{"Captain America", "2011", "id2", "movie", "poster2"},
			},
		}
		movieUsecaseMock.On("SearchMovies", testify.Anything, testify.Anything).Return(movieSearchResult, nil)

		serv := &movieServer{movieUsecaseMock}

		actualResult, _ := serv.SearchMovie(todoContext, &SearchMovieRequest{Searchword: "ironman"})
		assert.Equal(t, serv.convertMovieSearchToRPCResponse(movieSearchResult), actualResult)
	})
}

func TestGetMovieDetail(t *testing.T) {
	t.Run("[GetMovieDetail] malformed imdb id", func(t *testing.T) {
		testCases := []string{
			"aasdlkjasdlkasjd",
			"tt515a542",
			"ttt5154542",
			"123456",
		}

		movieUsecaseMock := &mock.MovieUsecase{}

		serv := &movieServer{movieUsecaseMock}

		for _, testCase := range testCases {
			_, err := serv.GetMovieDetail(todoContext, &GetMovieDetailRequest{Id: testCase})
			assert.Equal(t, incorrectImdbIDError, err)
		}
	})

	t.Run("[GetMovieDetail] valid imdb id but movieUsecase return an error", func(t *testing.T) {
		errMsg := "Oops, something happened"

		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("GetMovieDetailByID", testify.Anything).Return(&model.MovieDetail{}, errors.New(errMsg))

		serv := &movieServer{movieUsecaseMock}
		req := &GetMovieDetailRequest{Id: "tt1234567"}

		expectedErr := status.Error(codes.Internal, errMsg)
		_, actualErr := serv.GetMovieDetail(todoContext, req)
		if assert.Error(t, actualErr) {
			assert.Equal(t, expectedErr, actualErr)
		}
	})

	t.Run("[GetMovieDetail] detail has an error", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieUsecaseMock.On("GetMovieDetailByID", testify.Anything).Return(&model.MovieDetail{Error: "error"}, nil)

		serv := &movieServer{movieUsecaseMock}
		req := &GetMovieDetailRequest{Id: "tt1234567"}

		_, actualErr := serv.GetMovieDetail(todoContext, req)
		if assert.Error(t, actualErr) {
			assert.Equal(t, movieNotFoundError, actualErr)
		}
	})

	t.Run("[GetMovieDetail] no error, getMovieDetail has result", func(t *testing.T) {
		movieUsecaseMock := &mock.MovieUsecase{}
		movieDetailResult := &model.MovieDetail{
			"title", "year", "rated", "released", "runtime", "genre", "director", "writer", "actors",
			"plot", "language", "country", "awards", "poster", []model.MovieRating{{"source", "value"}}, "metascore",
			"imdbrating", "imdbvotes", "imdbid", "type", "dvd", "boxoffice", "production", "website", "true", "",
		}

		movieUsecaseMock.On("GetMovieDetailByID", testify.Anything).Return(movieDetailResult, nil)

		serv := &movieServer{movieUsecaseMock}
		req := &GetMovieDetailRequest{Id: "tt1234567"}

		actualResult, _ := serv.GetMovieDetail(todoContext, req)
		assert.Equal(t, serv.convertMovieDetailToRPCResponse(movieDetailResult), actualResult)
	})
}
