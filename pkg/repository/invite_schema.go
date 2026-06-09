package repository

import (
	"esproject"
	"fmt"

	"gorm.io/gorm"
)

type InviteSchema struct {
	db *gorm.DB
}

func NewInviteSchema(db *gorm.DB) *InviteSchema {
	return &InviteSchema{db: db}
}

func (r InviteSchema) CreateInvite(project int64, user int64) (*esproject.ProjectInvite, error) {
	invite := &esproject.ProjectInvite{
		Project: project,
		User:    user,
	}
	err := r.db.Create(invite).Error
	return invite, err
}

func (r InviteSchema) GetInvite(id int64) (*esproject.ProjectInvite, error) {
	invite := esproject.ProjectInvite{}
	err := r.db.Find(&invite, "id = ?", id).Error
	if invite.Id == -1 {
		err = fmt.Errorf("unknown invite %d", id)
	}
	return &invite, err
}

func (r InviteSchema) FindInvites(query, param string) ([]esproject.ProjectInvite, error) {
	var invites []esproject.ProjectInvite
	err := r.db.Find(&invites, query, param).Error
	return invites, err
}

func (r InviteSchema) DeleteInvite(id int64) error {
	invite := esproject.ProjectInvite{}
	err := r.db.Find(&invite, "id = ?", id).Error
	if invite.Id == -1 {
		err = fmt.Errorf("unknown invite %d", id)
	}
	r.db.Delete(invite)
	return err
}
