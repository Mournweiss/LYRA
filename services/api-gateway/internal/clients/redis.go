package clients

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/lyra/api-gateway/internal/models"
)

type RedisClient struct {
    client *redis.Client
}

func NewRedisClient(host, port string) *RedisClient {
    addr := fmt.Sprintf("%s:%s", host, port)
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
        DB:   0,
    })
    return &RedisClient{client: rdb}
}

func (r *RedisClient) SaveTask(ctx context.Context, task *models.Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return fmt.Errorf("failed to marshal task: %w", err)
    }
    key := r.taskKey(task.TaskID)
    return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisClient) GetTask(ctx context.Context, id string) (*models.Task, error) {
    key := r.taskKey(id)
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Not found
    } else if err != nil {
        return nil, fmt.Errorf("failed to get task from redis: %w", err)
    }
    var task models.Task
    if err := json.Unmarshal([]byte(val), &task); err != nil {
        return nil, fmt.Errorf("failed to unmarshal task: %w", err)
    }
    return &task, nil
}

func (r *RedisClient) UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus, result, errMsg string) error {
    task, err := r.GetTask(ctx, id)
    if err != nil {
        return err
    }
    if task == nil {
        return fmt.Errorf("task not found: %s", id)
    }
    task.Status = status
    task.Result = result
    task.Error = errMsg
    task.UpdatedAt = time.Now().Unix()
    return r.SaveTask(ctx, task)
}

func (r *RedisClient) ListPendingTaskIDs(ctx context.Context) ([]string, error) {
    keys, err := r.client.Keys(ctx, "task:*").Result()
    if err != nil {
        return nil, err
    }
    var pending []string
    for _, key := range keys {
        taskID := key[len("task:"):]
        task, err := r.GetTask(ctx, taskID)
        if err != nil || task == nil {
            continue
        }
        if task.Status == models.StatusPending {
            pending = append(pending, taskID)
        }
    }
    return pending, nil
}

func (r *RedisClient) taskKey(id string) string {
    return "task:" + id
}
