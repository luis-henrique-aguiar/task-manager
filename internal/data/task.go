package data

import (
	"sync"
	"time"
)

type Task struct {
    ID         int64     `json:"id"`
    CreatedAt  time.Time `json:"created_at"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    Done       bool      `json:"done"`
    Version    int32     `json:"version"`
}

type TaskModel struct {
    mu     sync.Mutex
    tasks  []Task
    nextID int64
}

func NewTaskModel() *TaskModel {
    return &TaskModel{
        tasks: []Task{},
        nextID: 1,
    }
}

func (m *TaskModel) Insert(task *Task) {
    m.mu.Lock()
    defer m.mu.Unlock()

    task.ID = m.nextID
    task.CreatedAt = time.Now()
    task.Version = 1

    m.tasks = append(m.tasks, *task)
    m.nextID++
}

func (m *TaskModel) GetAll() []Task {
    m.mu.Lock()
    defer m.mu.Unlock()

    return m.tasks
}