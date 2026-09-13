package task

import "database/sql"

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTask(row rowScanner) (Task, error) {
	var t Task
	err := row.Scan(&t.ID, &t.ParentID, &t.Title, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.Completed, &t.Deadline, &t.IsArchived)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) create(t Task) (Task, error) {
	query := "INSERT INTO tasks (parent_id, title, description, completed, deadline) VALUES (?, ?, ?, ?, ?) RETURNING *"
	row := r.db.QueryRow(query, t.ParentID, t.Title, t.Description, t.Completed, t.Deadline)
	task, err := scanTask(row)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) update(t Task) (Task, error) {
	query := "UPDATE tasks SET parent_id = ?, title = ?, description = ?, completed = ?, deadline = ?, is_archived = ? WHERE id = ? RETURNING *"
	row := r.db.QueryRow(query, t.ParentID, t.Title, t.Description, t.Completed, t.Deadline, t.IsArchived, t.ID)
	task, err := scanTask(row)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) getByID(id int64) (Task, error) {
	query := "SELECT * FROM tasks WHERE id = ?"
	row := r.db.QueryRow(query, id)

	task, err := scanTask(row)
	if err != nil {
		return Task{}, err
	}
	return task, nil
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
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) delete(id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}