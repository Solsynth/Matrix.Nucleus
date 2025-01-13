package services

import (
	"fmt"

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

func CalcReleaseToInstall(product int, current, target string) ([]models.ProductRelease, error) {
	var targetRelease models.ProductRelease
	if err := database.C.
		Where("product_id = ? AND version = ?", product, target).
		First(&targetRelease).Error; err != nil {
		return nil, fmt.Errorf("target release was not found: %v", err)
	}

	if targetRelease.Type == models.ReleaseTypeFull {
		return []models.ProductRelease{targetRelease}, nil
	}

	var lastFullRelease models.ProductRelease
	if err := database.C.
		Where("product_id = ? AND type = ?", product, models.ReleaseTypeFull).
		Order("created_at DESC").
		First(&lastFullRelease).Error; err != nil {
		return nil, fmt.Errorf("failed to find last full release: %v", err)
	}

	var plannedRelease []models.ProductRelease
	if err := database.C.
		Where("product_id = ? AND version > ? AND version <= ?", product, lastFullRelease.Version, target).
		Order("version ASC").
		Find(&plannedRelease).Error; err != nil {
		return nil, fmt.Errorf("failed to find planned releases: %v", err)
	}

	return plannedRelease, nil
}

func ListRelease(product int, take, offset int) ([]models.ProductRelease, error) {
	var items []models.ProductRelease
	if err := database.C.
		Where("product_id = ?", product).Order("created_at DESC").
		Preload("Meta").
		Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

func GetRelease(id uint) (models.ProductRelease, error) {
	var item models.ProductRelease
	if err := database.C.Where("id = ?", id).
		Preload("Meta").First(&item).Error; err != nil {
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
