package services

import (
	"git.solsynth.dev/matrix/nucleus/pkg/internal/database"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CountProduct() (int64, error) {
	var count int64
	if err := database.C.Model(&models.Product{}).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func CountCreatedProduct(user uint) (int64, error) {
	var count int64
	if err := database.C.Model(&models.Product{}).Where("account_id = ?", user).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func ListProduct(take, offset int) ([]models.Product, error) {
	var items []models.Product
	if err := database.C.Limit(take).Offset(offset).Preload("Meta").Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

func ListCreatedProduct(user uint, take, offset int) ([]models.Product, error) {
	var items []models.Product
	if err := database.C.
		Where("account_id = ?", user).
		Preload("Meta").
		Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

func GetProduct(id uint) (models.Product, error) {
	var item models.Product
	if err := database.C.Where("id = ?", id).Preload("Meta").First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func GetProductWithUser(id uint, user uint) (models.Product, error) {
	var item models.Product
	if err := database.C.Where("id = ? AND account_id = ?", id, user).Preload("Meta").First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func GetProductByAlias(alias string) (models.Product, error) {
	var item models.Product
	if err := database.C.Where("alias = ?", alias).Preload("Meta").First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func NewProduct(item models.Product) (models.Product, error) {
	if err := database.C.Create(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func UpdateProduct(item models.Product) (models.Product, error) {
	if err := database.C.Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Save(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func DeleteProduct(item models.Product) (models.Product, error) {
	if err := database.C.Select(clause.Associations).Delete(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}
