package group

type Group struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}


// M2M relationship between groups and tasks
type TaskGroup struct {
	GroupID int64 `db:"group_id"`
	TaskID  int64 `db:"task_id"`
}