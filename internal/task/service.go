package task

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TaskService struct {
	repo *TaskRepository
}

func NewTaskService(db *sqlx.DB) *TaskService {
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

type TaskOutput struct {
	Task
	SubTasks	[]TaskOutput `json:"subtasks"`
}

func toTaskOutput(parentTask Task, subTasks []TaskOutput) TaskOutput {
	return TaskOutput{
		Task: parentTask,
		SubTasks: subTasks,
	}
}

func getTaskMaps(tasks []Task) (map[int64]Task, map[int64][]Task) {
	// returns 2 maps
	// 1st is Task by it's ID
	// 2nd is Task list by it's ParentID
	taskMap := make(map[int64]Task)
	childreByParentMap := make(map[int64][]Task)
	for _, t := range tasks {
		taskMap[t.ID] = t
		childreByParentMap[t.ParentID] = append(childreByParentMap[t.ParentID], t)
	}
	return taskMap, childreByParentMap
}

func (ts *TaskService) buildTaskTree(rootId int64) (TaskOutput, error) {
	allTasks, err := ts.repo.getAll()
	if err != nil {
		return TaskOutput{}, nil
	}

	taskById, childrenByParent := getTaskMaps(allTasks)

	root, ok := taskById[rootId]
	if !ok {
		return TaskOutput{}, sql.ErrNoRows
	}
	return assembleTree(root, childrenByParent), nil
}

func (ts *TaskService) buildAllTasks() ([]TaskOutput, error) {
	allTasks, err := ts.repo.getAll()
	if err != nil {
		return nil, err
	}

	_, childrenByParent := getTaskMaps(allTasks)
	outputTasks := assembleAlltasksTree(allTasks, childrenByParent)

	return outputTasks, nil
}

func assembleTree(task Task, childrenByParent map[int64][]Task) TaskOutput {
	children := childrenByParent[task.ID]
	subOutPuts := make([]TaskOutput, 0, len(children))

	for _, child := range children{
		subOutPuts = append(subOutPuts, assembleTree(child, childrenByParent))
	}

	return toTaskOutput(task, subOutPuts)
}

func assembleAlltasksTree(tasks []Task, childrenByParent map[int64][]Task) []TaskOutput {
	var outputTasks []TaskOutput
	for _, task := range tasks {
		ot := assembleTree(task, childrenByParent)
		outputTasks = append(outputTasks, ot)
	}
	return outputTasks
}

func (ts *TaskService) CreateTask(input TaskInput) (TaskOutput, error) {
	newTask := Task{
		ParentID:    input.ParentID,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
	}
	taskObj, err := ts.repo.create(newTask)
	if err != nil {
		return TaskOutput{}, err
	}
	task := toTaskOutput(taskObj, nil)
	return task, nil
}

func (ts *TaskService) UpdateTask(input Task) (TaskOutput, error) {
	task, err := ts.repo.update(input)
	if err != nil {
		return TaskOutput{}, err
	}
	return ts.buildTaskTree(task.ID)
}

func (ts *TaskService) GetTaskByID(id int64) (TaskOutput, error) {
	return ts.buildTaskTree(id)
}

func (ts *TaskService) GetAllTasks() ([]TaskOutput, error) {
	tasks, err := ts.buildAllTasks()
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