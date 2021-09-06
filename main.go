package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	server "github.com/zenkobert/sbtest-2/delivery/grpc"
	repo "github.com/zenkobert/sbtest-2/repository"
	usecase "github.com/zenkobert/sbtest-2/usecase"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	grpcPort, restPort, apiKey string
)

func init() {
	err := godotenv.Load("shouldnotbeuploaded.env")
	if err != nil {
		log.Fatalln(err)
	}

	grpcPort = getEnvVariable("GRPC_PORT")
	restPort = getEnvVariable("REST_PORT")
	apiKey = getEnvVariable("API_KEY")
}

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
	movieRepo := repo.NewMovieRepo(apiKey)
	movieUsecase := usecase.NewMovieUsecase(movieRepo)
	movieServer := server.NewMovieServer(movieUsecase)

	grpcServer := grpc.NewServer()
	server.RegisterSearchMovieServer(grpcServer, movieServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Println(err)
	}

	log.Printf("GRPC Server Started. Listening to port %s", grpcPort)
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
		fmt.Sprintf("127.0.0.1:%s", grpcPort),
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Println(err)
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("REST HTTP Server Started. Listening to port %s", restPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", restPort), mux)
}

func getEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Can't find env variable : %s\n", key)
	}

	return value
}
