package handlers

import (
    "context"
    "github.com/lyra/api-gateway/internal/clients"
    "github.com/lyra/api-gateway/internal"
    pb "github.com/lyra/api-gateway/internal/pb"
)

func TranscribeHandler(ctx context.Context, req *pb.TranscribeRequest, whisperServiceAddr string) (*pb.TranscribeResponse, error) {
    if req == nil || len(req.FileContent) == 0 {
        return nil, internal.ValidationErrorf("file_content", "Audio or video file content is required")
    }
    resp, err := clients.ProxyTranscribe(ctx, req, whisperServiceAddr)
    if err != nil {
        return nil, internal.HandlerErrorf("PROXY_ERROR", "Failed to proxy transcription: %v", err)
    }
    return resp, nil
}
