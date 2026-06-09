package repository

import (
	"esproject"
	"fmt"

	"gorm.io/gorm"
)

type ProjectSchema struct {
	db *gorm.DB
}

func NewProjectSchema(db *gorm.DB) *ProjectSchema {
	return &ProjectSchema{db: db}
}

func (r ProjectSchema) CreateProject(owner int64, name string, tags []string) (*esproject.Project, error) {
	project := &esproject.Project{
		Name:    name,
		OwnerId: owner,
		Tags:    tags,
		Members: []int64{},
	}
	err := r.db.Create(project).Error
	fmt.Printf("Project created %d\n", project.Id)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	return project, err
}

func (r ProjectSchema) SaveProject(project *esproject.Project) (int64, error) {
	err := r.db.Save(project).Error
	return project.Id, err
}

func (r ProjectSchema) GetProject(id int64) (*esproject.Project, error) {
	project := esproject.Project{}
	err := r.db.Find(&project, "id = ?", id).Error
	if project.Id == -1 {
		err = fmt.Errorf("unknown account %d", id)
	}
	return &project, err
}

func (r ProjectSchema) GetProjects(start, limit int) []int64 {
	var result []esproject.Project
	r.db.Select("id").Limit(limit).Offset(start).Find(&result)
	var ans = make([]int64, len(result))
	for i, res := range result {
		ans[i] = res.Id
	}
	return ans
}

func (r ProjectSchema) GetProjectsCount() int64 {
	var count int64 = 0
	r.db.Model(&esproject.Project{}).Count(&count)
	return count
}

func (r ProjectSchema) DeleteProject(id int64) error {
	project := esproject.Project{}
	err := r.db.Find(&project, "id = ?", id).Error
	if project.Id == 0 {
		err = fmt.Errorf("unknown project %d", id)
	}
	r.db.Delete(project)
	return err
}
