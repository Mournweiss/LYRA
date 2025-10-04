package clients

import (
	context "context"
	log "log"

	"github.com/lyra/api-gateway/internal"
	"google.golang.org/grpc"
   	pb "github.com/lyra/api-gateway/internal/pb"
)

func ProxyTranscribe(ctx context.Context, req *pb.TranscribeRequest, serviceAddr string) (*pb.TranscribeResponse, error) {
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	
	if err != nil {
		log.Printf("Failed to connect to whisper-service: %v", err)
		return nil, internal.UpstreamErrorf("Failed to connect to whisper-service: %v", err)
	}

	defer conn.Close()

	client := pb.NewWhisperServiceClient(conn)
	return client.Transcribe(ctx, req)
}

