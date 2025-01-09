package services

import (
	"git.solsynth.dev/matrix/nucleus/pkg/internal/database"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CountRelease(product int) (int64, error) {
	var count int64
	if err := database.C.Model(&models.ProductRelease{}).Where("product_id = ?", product).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func ListRelease(product int, take, offset int) ([]models.ProductRelease, error) {
	var items []models.ProductRelease
	if err := database.C.
		Where("product_id = ?", product).Preload("Meta").
		Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

func GetRelease(id uint) (models.ProductRelease, error) {
	var item models.ProductRelease
	if err := database.C.Where("id = ?", id).Preload("Meta").First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func GetReleaseWithProduct(id uint, product uint) (models.ProductRelease, error) {
	var item models.ProductRelease
	if err := database.C.Where("id = ? AND product_id = ?", id, product).Preload("Meta").First(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func NewRelease(item models.ProductRelease) (models.ProductRelease, error) {
	if err := database.C.Create(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func UpdateRelease(item models.ProductRelease) (models.ProductRelease, error) {
	if err := database.C.Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Save(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func DeleteRelease(item models.ProductRelease) (models.ProductRelease, error) {
	if err := database.C.Select(clause.Associations).Delete(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}
