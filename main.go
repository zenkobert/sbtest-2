package main

import (
	"fmt"
	"log"
	"net"

	server "github.com/zenkobert/sbtest-2/delivery/grpc/movie"
	repo "github.com/zenkobert/sbtest-2/repository"
	usecase "github.com/zenkobert/sbtest-2/usecase"
	"google.golang.org/grpc"
)

func main() {
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
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
