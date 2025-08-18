package services

import (
	"errors"
	"final-project/models"
	"final-project/repositories"
	"final-project/utils"
	"fmt"
)


type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (as *UserService) Register(user *models.User) (string, uint, error) {

	// check if user already exists
	existingUser, err := as.repo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return "", 0, errors.New("email already registered")
	}

	// hash the password before di save
	if err := user.HashPassword(user.Password); err != nil {
		return "", 0, errors.New("failed to hashing password")
	}

	// Save user
	if err := as.repo.Create(user); err != nil {
		return "", 0, errors.New("failed to create user")
	}

	// generate token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}

	return token, user.ID, nil
}

func (as *UserService) Login(loginRequest *models.LoginRequest) (string, error) {

	// checker user is exist
	user, err := as.repo.FindByEmail(loginRequest.Email)
	  if err != nil || user.ID == 0 {
        return "", errors.New("invalid email or password")
    }
	// checker password
	if err := user.CheckPassword(loginRequest.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	fmt.Println("Generating token for user service:", user.ID, "with role:", user.Role)
	// generate token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}