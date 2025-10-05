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
    ID        string     `json:"id"`
    Status    TaskStatus `json:"status"`
    Result    string     `json:"result,omitempty"`
    Error     string     `json:"error,omitempty"`
    CreatedAt int64      `json:"created_at"`
    UpdatedAt int64      `json:"updated_at"`
}

func NewTask(id string) *Task {
    now := time.Now().Unix()
    return &Task{
        ID:        id,
        Status:    StatusPending,
        CreatedAt: now,
        UpdatedAt: now,
    }
}
