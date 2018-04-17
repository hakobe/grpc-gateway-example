package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb_articles "github.com/hakobe/grpc-gateway-example/articles"
)

func serve(hostPort string, grpcHostPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb_articles.RegisterArticlesServiceHandlerFromEndpoint(ctx, mux, grpcHostPort, opts)
	if err != nil {
		return errors.Wrap(err, "Cannot register article service handler")
	}

	log.Println("Starting gateway on " + hostPort)
	return http.ListenAndServe(hostPort, mux)
}

func main() {
	hostPort := os.Getenv("HOST_PORT")
	if hostPort == "" {
		hostPort = "0.0.0.0:5050"
	}
	grpcHostPort := os.Getenv("GRPC_HOST_PORT")
	if grpcHostPort == "" {
		grpcHostPort = "0.0.0.0:5000"
	}

	err := serve(hostPort, grpcHostPort)
	if err != nil {
		log.Fatal(err.Error())
	}
}
