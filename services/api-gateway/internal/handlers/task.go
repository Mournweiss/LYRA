package handlers

import (
    "context"
    "log"
    "time"
    "github.com/google/uuid"
    "github.com/lyra/api-gateway/internal"
    "github.com/lyra/api-gateway/internal/clients"
    pb "github.com/lyra/api-gateway/internal/pb"
)

func CreateTranscriptionTaskHandler(ctx context.Context, req *pb.TranscribeRequest, redisClient *clients.RedisClient) (*pb.CreateTaskResponse, error) {
    if req == nil || len(req.FileContent) == 0 {
        return &pb.CreateTaskResponse{TaskId: "", Error: "Audio or video file content is required"}, internal.ValidationErrorf("file_content", "Audio or video file content is required")
    }
    taskID := uuid.NewString()
    task := internal.NewTask(taskID)
    err := redisClient.SaveTask(ctx, task)
    if err != nil {
        return &pb.CreateTaskResponse{TaskId: "", Error: "Failed to create task"}, internal.HandlerErrorf("REDIS_ERROR", "Failed to save task in Redis: %v", err)
    }
    return &pb.CreateTaskResponse{TaskId: taskID, Error: ""}, nil
}

func GetTaskStatusHandler(ctx context.Context, req *pb.GetTaskStatusRequest, redisClient *clients.RedisClient) (*pb.GetTaskStatusResponse, error) {
    if req == nil || req.TaskId == "" {
        return &pb.GetTaskStatusResponse{Status: "ERROR", Error: "Task ID is required"}, internal.ValidationErrorf("task_id", "Task ID is required")
    }
    task, err := redisClient.GetTask(ctx, req.TaskId)
    if err != nil {
        return &pb.GetTaskStatusResponse{Status: "ERROR", Error: "Failed to get task status"}, internal.HandlerErrorf("REDIS_ERROR", "Failed to get task from Redis: %v", err)
    }
    if task == nil {
        return &pb.GetTaskStatusResponse{Status: "NOT_FOUND", Error: "Task not found"}, nil
    }
    return &pb.GetTaskStatusResponse{
        Status:    string(task.Status),
        Result:    task.Result,
        Error:     task.Error,
        CreatedAt: task.CreatedAt,
        UpdatedAt: task.UpdatedAt,
    }, nil
}

func StartTaskWorker(ctx context.Context, redisClient *clients.RedisClient, whisperServiceAddr string, concurrency int) {
    taskChan := make(chan string, 100)

    go func() {
        for {
            select {
            case <-ctx.Done():
                log.Println("Task scanner stopped")
                close(taskChan)
                return
            default:
                keys, err := redisClient.Client().Keys(ctx, "task:*").Result()
                if err != nil {
                    log.Printf("Failed to scan Redis: %v", err)
                    continue
                }
                for _, key := range keys {
                    task, err := redisClient.GetTask(ctx, key[len("task:"):])
                    if err != nil || task == nil {
                        continue
                    }
                    if task.Status == internal.StatusPending {
                        select {
                        case taskChan <- task.ID:
                        case <-ctx.Done():
                            close(taskChan)
                            return
                        }
                    }
                }
            }
        }
    }()

    for i := 0; i < concurrency; i++ {
        go func(workerID int) {
            for taskID := range taskChan {
                select {
                case <-ctx.Done():
                    log.Printf("[Worker-%d] Stopping", workerID)
                    return
                default:
                    processTask(ctx, redisClient, whisperServiceAddr, taskID, workerID)
                }
            }
        }(i + 1)
    }
}

func processTask(ctx context.Context, redisClient *clients.RedisClient, whisperServiceAddr, taskID string, workerID int) {
    task, err := redisClient.GetTask(ctx, taskID)
    if err != nil || task == nil {
        log.Printf("[Worker-%d] Failed to fetch task %s: %v", workerID, taskID, err)
        return
    }
    log.Printf("[Worker-%d] Processing task %s", workerID, task.ID)
    _ = redisClient.UpdateTaskStatus(ctx, task.ID, internal.StatusRunning, "", "")
    req := &pb.TranscribeRequest{
        FileContent: []byte(task.Result),
        FileName:    "async-task.wav",
    }
    resp, err := clients.ProxyTranscribe(ctx, req, whisperServiceAddr)
    if err != nil {
        log.Printf("[Worker-%d] Transcription failed for task %s: %v", workerID, task.ID, err)
        _ = redisClient.UpdateTaskStatus(ctx, task.ID, internal.StatusError, "", err.Error())
        return
    }
    _ = redisClient.UpdateTaskStatus(ctx, task.ID, internal.StatusSuccess, resp.Text, resp.Error)
    log.Printf("[Worker-%d] Task %s completed", workerID, task.ID)
}
