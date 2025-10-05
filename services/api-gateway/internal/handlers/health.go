package handlers

import (
    "context"
    pb "github.com/lyra/api-gateway/internal/pb"
)

func HealthCheckHandler(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    return &pb.HealthCheckResponse{Status: "SERVING"}, nil
}
