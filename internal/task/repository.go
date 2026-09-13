package task

import "database/sql"

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) create(t Task) (Task, error) {
	query := "INSERT INTO tasks (parent_id, title, description, created_at, updated_at, completed, deadline) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, t.ParentID, t.Title, t.Description, t.CreatedAt, t.UpdatedAt, t.Completed, t.Deadline)
	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:          id,
		ParentID:    t.ParentID,
		Title:       t.Title,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Completed:   t.Completed,
		Deadline:    t.Deadline,
	}, nil
}

func (r *TaskRepository) update(t Task) error {
	query := "UPDATE tasks SET parent_id = ?, title = ?, description = ?, created_at = ?, updated_at = ?, completed = ?, deadline = ?, is_archived = ? WHERE id = ?"
	_, err := r.db.Exec(query, t.ParentID, t.Title, t.Description, t.CreatedAt, t.UpdatedAt, t.Completed, t.Deadline, t.IsArchived, t.ID)
	return err
}

func (r *TaskRepository) getByID(id int64) (Task, error) {
	query := "SELECT * FROM tasks WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var t Task
	err := row.Scan(&t.ID, &t.ParentID, &t.Title, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.Completed, &t.Deadline, &t.IsArchived)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) getAll() ([]Task, error) {
	query := "SELECT * FROM tasks"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.ParentID, &t.Title, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.Completed, &t.Deadline, &t.IsArchived)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) delete(id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}