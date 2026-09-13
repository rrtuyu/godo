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

type TaskInput struct {
	ParentID    int64  `json:"parent_id" default:"null"`
	Title       string `json:"title" default:"New Task"`
	Description string `json:"description" default:""`
	Deadline    string `json:"deadline" default:""`
	Completed   bool   `json:"completed" default:"false"`
}

func (ts *TaskService) CreateTask(input TaskInput) (Task, error) {
	newTask := Task{
		ParentID:    input.ParentID,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
	}
	task, err := ts.repo.create(newTask)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (ts *TaskService) UpdateTask(input Task) (Task, error) {
	task, err := ts.repo.update(input)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (ts *TaskService) GetTaskByID(id int64) (Task, error) {
	task, err := ts.repo.getByID(id)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (ts *TaskService) GetAllTasks() ([]Task, error) {
	tasks, err := ts.repo.getAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskService) DeleteTask(id int64) error {
	err := ts.repo.delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) ArchiveTask(id int64) error {
	task, err := ts.repo.getByID(id)
	if err != nil {
		return err
	}
	task.IsArchived = true
	_, err = ts.repo.update(task)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) UnarchiveTask(id int64) error {
	task, err := ts.repo.getByID(id)
	if err != nil {
		return err
	}
	task.IsArchived = false
	_, err = ts.repo.update(task)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) MarkTaskCompleted(id int64) error {
	task, err := ts.repo.getByID(id)
	if err != nil {
		return err
	}
	task.Completed = true
	_, err = ts.repo.update(task)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) MarkTaskIncomplete(id int64) error {
	task, err := ts.repo.getByID(id)
	if err != nil {
		return err
	}
	task.Completed = false
	_, err = ts.repo.update(task)
	if err != nil {
		return err
	}
	return nil
}