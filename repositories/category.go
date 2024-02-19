package repositories

import (
	"github.com/google/uuid"
	"github.com/rammyblog/spendwise/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	return repo.db.Create(category).Error
}

func (repo *CategoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := repo.db.First(&category, id).Error
	return &category, err
}

func (repo *CategoryRepository) FindAll() ([]models.Category, error) {

	var categories []models.Category
	err := repo.db.Find(&categories).Error
	return categories, err
}

func (repo *CategoryRepository) Update(id uuid.UUID, category *models.Category) error {
	var existingCategory models.Category
	err := repo.db.First(&existingCategory, id).Error
	if err != nil {
		return err
	}
	return repo.db.Model(&existingCategory).Updates(category).Error
}
