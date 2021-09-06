package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
	"github.com/zenkobert/sbtest-2/domain/mocks"
	"google.golang.org/grpc"
)

func TestNewInterceptor(t *testing.T) {
	t.Run("[NewInterceptor]", func(t *testing.T) {
		expected := interceptor{&mocks.MovieUsecase{}}
		actual := NewInterceptor(&mocks.MovieUsecase{})

		assert.Equal(t, expected, actual)
	})
}

func TestUnary(t *testing.T) {
	t.Run("[Unary]", func(t *testing.T) {
		movieUsecase := mocks.MovieUsecase{}
		movieUsecase.On("LogToDB", testify.Anything).Return(nil)
		in := NewInterceptor(&movieUsecase)

		var handler = func(context.Context, interface{}) (interface{}, error) {
			return "abc", nil
		}

		result, err := in.Unary(context.TODO(), "interface{}", &grpc.UnaryServerInfo{}, handler)
		assert.Nil(t, err)
		assert.Equal(t, "abc", result)
	})
}

func TestLogToDB(t *testing.T) {
	t.Run("[logToDB] movieUsecase return error", func(t *testing.T) {
		movieUsecase := mocks.MovieUsecase{}
		movieUsecase.On("LogToDB", testify.Anything).Return(errors.New("error"))
		in := NewInterceptor(&movieUsecase)

		err := in.logToDB("")
		assert.Error(t, err)
	})

	t.Run("[logToDB] no error, positive test", func(t *testing.T) {
		movieUsecase := mocks.MovieUsecase{}
		movieUsecase.On("LogToDB", testify.Anything).Return(nil)
		in := NewInterceptor(&movieUsecase)

		err := in.logToDB("")
		assert.Nil(t, err)
	})
}
