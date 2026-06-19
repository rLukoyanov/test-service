package config

import (
	"os"
)

type Config struct {
	Port     string
	GrpcPort string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}
	return &Config{Port: port, GrpcPort: grpcPort}
}
