package controller

import (
	"github.com/rammyblog/spendwise/config"
	"github.com/rammyblog/spendwise/models"
	"github.com/rammyblog/spendwise/repositories"
)

func GetCategories() ([]models.Category, error) {
	categoryRepo := repositories.NewCategoryRepository(config.GlobalConfig.DB)
	categories, err := categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
