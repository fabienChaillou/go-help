package repository

import "gomock-example/models"

type UserRepository interface {
	GetUserByID(id int) (*models.User, error)
}
