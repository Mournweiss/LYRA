package server

import (
	"context"
	"log"
	"net"

	"github.com/lyra/api-gateway/internal/clients"
	handlers "github.com/lyra/api-gateway/internal/handlers"
	pb "github.com/lyra/api-gateway/internal/pb"
	"github.com/lyra/api-gateway/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedWhisperServiceServer
	whisperServiceAddr string
	redisClient       *clients.RedisClient
}

func (s *Server) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	return handlers.TranscribeHandler(ctx, req, s.whisperServiceAddr)
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return handlers.HealthCheckHandler(ctx, req)
}

func (s *Server) CreateTranscriptionTask(ctx context.Context, req *pb.CreateTranscriptionTaskRequest) (*pb.CreateTaskResponse, error) {
	return handlers.CreateTranscriptionTaskHandler(ctx, req, s.redisClient)
}

func (s *Server) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (*pb.GetTaskStatusResponse, error) {
	return handlers.GetTaskStatusHandler(ctx, req, s.redisClient)
}

func StartServer(cfg *config.Config) {
	log.Printf("Starting API Gateway on port %s (address: %s), whisper-service: %s", cfg.GatewayPort, cfg.GatewayAddress, cfg.WhisperServiceAddr)

	redisClient := clients.NewRedisClient(cfg.RedisHost, cfg.RedisPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handlers.StartTaskWorker(ctx, redisClient, cfg.WhisperServiceAddr, cfg.WorkerConcurrency)

	listenAddr := ":" + cfg.GatewayPort
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWhisperServiceServer(grpcServer, &Server{
		whisperServiceAddr: cfg.WhisperServiceAddr,
		redisClient:       redisClient,
	})
	log.Printf("API Gateway gRPC server listening on %s", listenAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
