package service

import (
	"esproject"
)

type Authorization interface {
	Authorize(header string) (*esproject.AuthorizedUser, error)
	AuthorizeById(id int64) (*esproject.AuthorizedUser, error)
	AuthorizeAndUpdatePicture(header string, picUrl string) (*esproject.AuthorizedUser, error)
}

type Service struct {
	Authorization
}

func NewService() *Service {
	return &Service{
		Authorization: NewAuthService(),
	}
}
