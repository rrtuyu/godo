package task

import (
	"database/sql"
)

type TaskService struct {
	repo *TaskRepository
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{
		repo: NewTaskRepository(db),
	}
}

func (ts *TaskService) CreateTask() (t Task) (Task, error) {}