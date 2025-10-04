package main

import (
	"context"
	"log"
	"net"

	"github.com/lyra/api-gateway/internal"
	"github.com/lyra/api-gateway/internal/clients"
   	pb "github.com/lyra/api-gateway/internal/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWhisperServiceServer
	whisperServiceAddr string
}

func (s *server) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	resp, err := clients.ProxyTranscribe(ctx, req, s.whisperServiceAddr)

	if err != nil {
		if _, ok := err.(*internal.UpstreamError); ok {
			log.Printf("Upstream error: %v", err)
		} else {
			log.Printf("Internal error: %v", err)
		}
		return nil, err
	}
	return resp, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "SERVING"}, nil
}

func main() {
	cfg, err := internal.LoadConfig()

	if err != nil {
		if _, ok := err.(*internal.ConfigError); ok {
			log.Fatalf("Configuration error: %v", err)
		} else {
			log.Fatalf("Startup error: %v", err)
		}
	}

	log.Printf("Starting API Gateway on port %s (address: %s), whisper-service: %s", cfg.GatewayPort, cfg.GatewayAddress, cfg.WhisperServiceAddr)

	listenAddr := ":" + cfg.GatewayPort
	lis, err := net.Listen("tcp", listenAddr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWhisperServiceServer(grpcServer, &server{whisperServiceAddr: cfg.WhisperServiceAddr})
	log.Printf("API Gateway gRPC server listening on %s", listenAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
