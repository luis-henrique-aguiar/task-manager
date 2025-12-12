package data

import (
	"database/sql"
	"errors"
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
    DB *sql.DB
}

func NewTaskModel(db *sql.DB) *TaskModel {
    return &TaskModel{DB: db}
}

func (m *TaskModel) Get(id int64) (*Task, error) {
    query := `
        SELECT id, created_at, title, content, done, version
        FROM tasks
        WHERE id = $1
    `

    var task Task

    err := m.DB.QueryRow(query, id).Scan(
        &task.ID,
        &task.CreatedAt,
        &task.Title,
        &task.Content,
        &task.Done,
        &task.Version,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrRecordNotFound
        }
        return nil, err
    }

    return &task, nil
}

func (m *TaskModel) Delete(id int64) error {
    query := `DELETE FROM tasks WHERE id = $1`

    result, err := m.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrRecordNotFound
    }

    return nil
}

func (m *TaskModel) Insert(task *Task) error {
    query := `
        INSERT INTO tasks (title, content, done, version, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

    task.CreatedAt = time.Now()
    task.Version = 1
    task.Done = false

    err := m.DB.QueryRow(
        query,
        task.Title,
        task.Content,
        task.Done,
        task.Version,
        task.CreatedAt,
    ).Scan(&task.ID)
    
    if err != nil {
        return err
    }

    return nil
}

func (m *TaskModel) GetAll() ([]Task, error) {
    query := `SELECT id, created_at, title, content, done, version FROM tasks`

    rows, err := m.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []Task

    for rows.Next() {
        var task Task
        err := rows.Scan(
            &task.ID,
            &task.CreatedAt,
            &task.Title,
            &task.Content,
            &task.Done,
            &task.Version,
        )
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (m *TaskModel) Update(task *Task) error {
    query := `
        UPDATE tasks
        SET title = $1, content = $2, done = $3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version
    `

    err := m.DB.QueryRow(
        query,
        task.Title,
        task.Content,
        task.Done,
        task.ID,
        task.Version,
    ).Scan(&task.Version)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrRecordNotFound
        }
        return err
    }

    return nil
}