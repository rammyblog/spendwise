package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rammyblog/spendwise/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(user *models.User) error {
	fmt.Println("User: ", user)
	return repo.db.Create(user).Error
}

func (repo *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := repo.db.First(&user, id).Error
	return &user, err
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repo *UserRepository) Update(email string, user *models.User) error {
	var existingUser models.User
	err := repo.db.Where("email = ?", email).First(&existingUser).Error
	if err != nil {
		return err
	}
	return repo.db.Model(&existingUser).Updates(user).Error

}

func (repo *UserRepository) Delete(user *models.User) error {
	return repo.db.Delete(user).Error
}

func (repo *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := repo.db.Find(&users).Error
	return users, err
}
