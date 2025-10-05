package internal

import "time"

type TaskStatus string

const (
    StatusPending   TaskStatus = "PENDING"
    StatusRunning   TaskStatus = "IN_PROGRESS"
    StatusSuccess   TaskStatus = "SUCCESS"
    StatusError     TaskStatus = "ERROR"
)

type Task struct {
    TaskID    string     `json:"task_id"`
    FileKey   string     `json:"file_key"`
    Status    TaskStatus `json:"status"`
    Result    string     `json:"result,omitempty"`
    Error     string     `json:"error,omitempty"`
    CreatedAt int64      `json:"created_at"`
    UpdatedAt int64      `json:"updated_at"`
}

func NewTask(taskID, fileKey string) *Task {
    now := time.Now().Unix()
    return &Task{
        TaskID:    taskID,
        FileKey:   fileKey,
        Status:    StatusPending,
        CreatedAt: now,
        UpdatedAt: now,
    }
}
