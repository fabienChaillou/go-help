package service

import (
	"gomock-example/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) GetUserName(id int) (string, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
