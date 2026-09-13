package group

import (
	"github.com/jmoiron/sqlx"
)


type GroupService struct {
	repo *GroupRepository
}

func NewGroupService(db *sqlx.DB) *GroupService {
	return &GroupService{
		repo: NewGroupRepository(db),
	}
}

type GroupInput struct {
	Name		string `json:"name" default:"New group"`
	Description string	`json:"description" default:""`
}

func (gs *GroupService) CreateGroup(input GroupInput) (Group, error) {
	newGroup := Group{
		Name:			input.Name,
		Description:	input.Name,
	}
	group, err := gs.repo.create(newGroup)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (gs *GroupService) UpdateGroup(input Group) (Group, error) {
	group, err := gs.repo.update(input)
	if err != nil {
		return Group{}, nil
	}

	return group, nil
}

func (gs *GroupService) DeleteGroup(id int64) error {
	if err := gs.repo.delete(id); err != nil {
		return err
	}
	return nil
}

func (gs *GroupService) AddTaskToGroup(taskId int64, groupId int64) error {
	if err := gs.repo.addTaskToGroup(taskId, groupId); err != nil {
		return err
	}
	return nil
}