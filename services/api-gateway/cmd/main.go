package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/lyra/api-gateway/proto"
)

type server struct {
	pb.UnimplementedWhisperServiceServer
}

func (s *server) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	return &pb.TranscribeResponse{
		Text:  "",
		Error: "Not implemented",
	}, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: "SERVING"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterWhisperServiceServer(grpcServer, &server{})
	log.Println("API Gateway gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
