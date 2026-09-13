package task

type Task struct {
	ID          int64  `db:"id" json:"id"`
	ParentID    int64  `db:"parent_id" json:"parent_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
	Completed   bool   `db:"completed" json:"completed"`
	Deadline    string `db:"deadline" json:"deadline"`
	IsArchived  bool   `db:"is_archived" json:"is_archived"`
}