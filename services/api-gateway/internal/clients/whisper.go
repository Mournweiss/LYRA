package clients

import (
	context "context"
	log "log"

	"google.golang.org/grpc"
	pb "github.com/lyra/api-gateway/internal"
)

const whisperServiceAddr = "whisper-service:50052"

func ProxyTranscribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	conn, err := grpc.Dial(whisperServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to whisper-service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewWhisperServiceClient(conn)
	return client.Transcribe(ctx, req)
}
