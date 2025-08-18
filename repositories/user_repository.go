package repositories

import (
	"final-project/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = ? LIMIT 1`
	if err := r.DB.Raw(query, email).Scan(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (name, email, password, role, created_at) VALUES (?, ?, ?, 'user', NOW())`
	return r.DB.Exec(query, user.Name, user.Email, user.Password).Error
}