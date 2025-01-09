package database

import (
	"git.solsynth.dev/matrix/nucleus/pkg/internal/models"
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{
	&models.Product{},
	&models.ProductMeta{},
	&models.ProductRelease{},
	&models.ProductReleaseMeta{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		AutoMaintainRange...,
	); err != nil {
		return err
	}

	return nil
}
