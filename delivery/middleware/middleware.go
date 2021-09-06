package middleware

import (
	"context"
	"fmt"

	model "github.com/zenkobert/sbtest-2/domain"
	"google.golang.org/grpc"
)

type interceptor struct {
	MovieUsecase model.MovieUsecase
}

func NewInterceptor(usecase model.MovieUsecase) interceptor {
	return interceptor{
		MovieUsecase: usecase,
	}
}

func (in *interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		record := fmt.Sprintf("%s/ %s", info.FullMethod, req)
		in.MovieUsecase.LogToDB(record)

		return handler(ctx, req)
	}
}
