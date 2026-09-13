package task

type Task struct {
	ID			int64   `db:"id"`
	ParentID	int64   `db:"parent_id"`
	Title		string  `db:"title"`
	Description	string  `db:"description"`
	CreatedAt	string  `db:"created_at"`
	UpdatedAt	string  `db:"updated_at"`
	Completed	bool    `db:"completed"`
	Deadline	string  `db:"deadline"`
	IsArchived	bool    `db:"is_archived"`
}


// TODO rethink necceserity of this model, maybe just use Task and add a ParentID field to it
type SubTask struct {
	Task
	ParentID	int64   `db:"parent_id"`
}