package data

import (
	"errors"
	"sync"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

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

func (m *TaskModel) Get(id int64) (*Task, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, task := range m.tasks {
        if task.ID == id {
            return &task, nil
        }
    }

    return nil, ErrRecordNotFound
}

func (m *TaskModel) Delete(id int64) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    index := -1
    for i, task := range m.tasks {
        if task.ID == id {
            index = i
            break
        }
    }

    if index == -1 {
        return ErrRecordNotFound
    }

    m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)

    return nil
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