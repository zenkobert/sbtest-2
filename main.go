package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	server "github.com/zenkobert/sbtest-2/delivery/grpc"
	repo "github.com/zenkobert/sbtest-2/repository"
	usecase "github.com/zenkobert/sbtest-2/usecase"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	g := errgroup.Group{}
	g.Go(func() error {
		return startGrpcServer()
	})
	g.Go(func() error {
		return startRestServer()
	})

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func startGrpcServer() error {
	movieRepo := repo.NewMovieRepo("faf7e5bb")
	movieUsecase := usecase.NewMovieUsecase(movieRepo)
	movieServer := server.NewMovieServer(movieUsecase)

	grpcServer := grpc.NewServer()
	server.RegisterSearchMovieServer(grpcServer, movieServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "8080"))
	if err != nil {
		log.Println(err)
	}

	log.Println("GRPC Server Started")
	return grpcServer.Serve(listener)
}

func startRestServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := server.RegisterSearchMovieHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:8080",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Println("REST Service Started")
	return http.ListenAndServe(":8081", mux)
}
