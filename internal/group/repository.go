package group

import (
	"godo/internal/task"

	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) create(g Group) (Group, error) {
	var group Group
	query := "INSERT INTO groups (name, description) VALUES (?, ?) RETURNING"
	err := r.db.Get(&group, query, g.Name, g.Description)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (r *GroupRepository) update(g Group) (Group, error) {
	var group Group
	query := "UPDATE groups SET name = ? description = ? WHERE id = ? RETURNING *"
	err := r.db.Get(&group, query, g.Name, g.Description, g.ID)
	if err != nil {
		return Group{}, nil
	}

	return group, nil
}

func (r *GroupRepository) getByID(id int64) (Group, error) {
	var group Group
	query := "SELECT * FROM groups WHERE id = ?"
	err := r.db.Get(&group, query, id)
	if err != nil {
		return Group{}, nil
	}
	return group, nil
}

func (r *GroupRepository) getGroupTasks(groupID int64) ([]task.Task, error) {
	var tasks []task.Task
	query := "SELECT t.* FROM tasks t JOIN group_task gt ON t.id = gt.task_id WHERE gt.group_id = ?"
	err := r.db.Select(&tasks, query, groupID)
	if err != nil {
		return nil, err
	}
	return tasks, nil 
}

func (r *GroupRepository) getAll() ([]Group, error) {
	var groups []Group

	query := "SELECT * FROM groups"
	err := r.db.Select(&groups, query)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) delete(id int64) error {
	query := "DELETE FROM groups WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *GroupRepository) addTaskToGroup(taskID int64, groupID int64) error {
	query := "INSERT INTO group_task (group_id, task_id) VALUES (?, ?)"
	_, err := r.db.Exec(query, groupID, taskID)
	return err
}

func (r *GroupRepository) removeTaskFromGroup(taskID int64, groupID int64) error {
	query := "DELETE FROM group_task WHERE group_id = ? AND task_id = ?"
	_, err := r.db.Exec(query, groupID, taskID)
	return err
}