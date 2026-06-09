package repository

import (
	"esproject"

	"gorm.io/gorm"
)

type Project interface {
	CreateProject(owner int64, name string, tags []string) (*esproject.Project, error)
	SaveProject(user *esproject.Project) (int64, error)
	GetProject(id int64) (*esproject.Project, error)
	GetProjects(start, limit int) []int64
	DeleteProject(id int64) error
	GetProjectsCount() int64
}

type ProjectInvite interface {
	CreateInvite(project int64, user int64) (*esproject.ProjectInvite, error)
	GetInvite(id int64) (*esproject.ProjectInvite, error)
	FindInvites(query, param string) ([]esproject.ProjectInvite, error)
	DeleteInvite(id int64) error
}

type Repository struct {
	Project
	ProjectInvite
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Project:       NewProjectSchema(db),
		ProjectInvite: NewInviteSchema(db),
	}
}
