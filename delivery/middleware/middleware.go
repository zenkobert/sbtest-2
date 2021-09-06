package middleware

import (
	"context"
	"fmt"
	"log"

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

func (in *interceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	record := fmt.Sprintf("%s/ %s", info.FullMethod, req)

	// don't need to wait until logging finish
	// client need to be served asap
	go in.logToDB(record)

	return handler(ctx, req)
}

func (in *interceptor) logToDB(record string) error {
	err := in.MovieUsecase.LogToDB(record)
	if err != nil {
		log.Println(err)
	}

	return err
}
