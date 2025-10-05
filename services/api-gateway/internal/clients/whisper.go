package clients

import (
	context "context"
	log "log"
	pb "github.com/lyra/api-gateway/internal/pb"
	"google.golang.org/grpc"
)

func ProxyTranscribe(ctx context.Context, req *pb.TranscribeRequest, serviceAddr string) (*pb.TranscribeResponse, error) {
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to whisper-service: %v", err)
		return nil, err
	}
	defer conn.Close()
	client := pb.NewWhisperServiceClient(conn)
	return client.Transcribe(ctx, req)
}

