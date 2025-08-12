package services

import (
	"errors"
	"final-project/models"
	"final-project/utils"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (as *UserService) Register(user *models.User) (string, uint, error) {
	// hash the password before di save
	if err := user.HashPassword(user.Password); err != nil {
		return "", 0, errors.New("failed to hashing password")
	}

	// create in db
	query := `INSERT INTO users (name, email, password, role, created_at) VALUES (?, ?, ?, 'user', NOW())`
	result := as.DB.Exec(query, user.Name, user.Email, user.Password)
	if result.Error != nil {
		return "", 0, errors.New("failed to create user")
	}
	fmt.Println("Rows affected:", result.RowsAffected)
	userId := uint(result.RowsAffected)
	user.ID = userId

	// generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}

	return token, user.ID, nil
}

func (as *UserService) Login(loginRequest *models.LoginRequest) (string, uint, error) {
	var user models.User

	// checker user is exist
	query := `SELECT * FROM users WHERE email = ? LIMIT 1`
	if err := as.DB.Raw(query, loginRequest.Email).Scan(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, errors.New("invalid email or password")
		}
		return "", 0, err
	}

	// checker password
	if err := user.CheckPassword(loginRequest.Password); err != nil {
		return "", 0, errors.New("invalid email or password")
	}

	// generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}
	return token, user.ID, nil
}