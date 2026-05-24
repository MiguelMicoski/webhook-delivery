package service

import "awesomeProject/internal/model"

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) FindAll() []model.User {
	return []model.User{
		{
			ID:    1,
			Name:  "Ana Silva",
			Email: "ana@example.com",
		},
		{
			ID:    2,
			Name:  "Bruno Costa",
			Email: "bruno@example.com",
		},
	}
}
