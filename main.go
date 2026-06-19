package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"test-service/config"
	"test-service/grpcserver"
	"test-service/handlers"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/one",   handlers.One)
	mux.HandleFunc("/two",   handlers.Two)
	mux.HandleFunc("/three", handlers.Three)
	mux.HandleFunc("/four",  handlers.Four)

	httpAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("starting HTTP server on %s", httpAddr)
	go func() {
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	grpcAddr := fmt.Sprintf(":%s", cfg.GrpcPort)
	log.Printf("starting gRPC server on %s", grpcAddr)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}
	s := grpc.NewServer()
	grpcserver.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
