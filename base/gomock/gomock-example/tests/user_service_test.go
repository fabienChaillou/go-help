package tests

import (
	"errors"
	"gomock-example/models"
	"gomock-example/repository"
	"gomock-example/service"
	"testing"

	// "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGetUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	mockRepo.
		EXPECT().
		GetUserByID(1).
		Return(&models.User{ID: 1, Name: "Alice"}, nil)

	userService := service.UserService{Repo: mockRepo}

	name, err := userService.GetUserName(1)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", name)
}

func TestGetUserName_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	mockRepo.
		EXPECT().
		GetUserByID(2).
		Return(nil, errors.New("user not found"))

	userService := service.UserService{Repo: mockRepo}

	name, err := userService.GetUserName(2)

	assert.Error(t, err)
	assert.Equal(t, "", name)
}
