package handlers

import (
    "context"
    "github.com/lyra/api-gateway/internal/clients"
    "github.com/lyra/api-gateway/internal/errors"
    pb "github.com/lyra/api-gateway/internal/pb"
)

func TranscribeHandler(ctx context.Context, req *pb.TranscribeRequest, whisperServiceAddr string) (*pb.TranscribeResponse, error) {
    if req == nil || req.FileKey == "" {
        return nil, errors.ValidationErrorf("file_key", "File key is required")
    }
    resp, err := clients.ProxyTranscribe(ctx, req, whisperServiceAddr)
    if err != nil {
        return nil, errors.HandlerErrorf("PROXY_ERROR", "Failed to proxy transcription: %v", err)
    }
    return resp, nil
}
