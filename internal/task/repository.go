package task

import "github.com/jmoiron/sqlx"

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) create(t Task) (Task, error) {
	var task Task
	query := "INSERT INTO tasks (parent_id, title, description, completed, deadline) VALUES (?, ?, ?, ?, ?) RETURNING *"
	err := r.db.Get(&task, query, t.ParentID, t.Title, t.Description, t.Completed, t.Deadline)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) update(t Task) (Task, error) {
	var task Task
	query := "UPDATE tasks SET parent_id = ?, title = ?, description = ?, completed = ?, deadline = ?, is_archived = ? WHERE id = ? RETURNING *"
	err := r.db.Get(&task, query, t.ParentID, t.Title, t.Description, t.Completed, t.Deadline, t.IsArchived, t.ID)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) getByID(id int64) (Task, error) {
	var task Task
	query := "SELECT * FROM tasks WHERE id = ?"
	err := r.db.Get(&task, query, id)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) getAll() ([]Task, error) {
	var tasks []Task
	query := "SELECT * FROM tasks"
	err := r.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) delete(id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TaskRepository) getSubTasks(parentId int64) ([]Task, error) {
	var tasks []Task
	query := "SELECT * FROM tasks WHERE parent_id = ?"
	if err := r.db.Select(&tasks, query, parentId); err != nil {
		return nil, err
	}
	return tasks, nil
}