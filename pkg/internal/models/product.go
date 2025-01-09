package models

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cruda"
	"gorm.io/datatypes"
)

type Product struct {
	cruda.BaseModel

	Icon        string                      `json:"icon"` // random id of atttachment
	Name        string                      `json:"name"`
	Alias       string                      `json:"alias" gorm:"uniqueIndex"`
	Description string                      `json:"description"`
	Previews    datatypes.JSONSlice[string] `json:"previews"` // random id of attachments
	Tags        datatypes.JSONSlice[string] `json:"tags"`

	Meta      ProductMeta      `json:"meta" gorm:"foreignKey:ProductID"`
	Releases  []ProductRelease `json:"releases" gorm:"foreignKey:ProductID"`
	AccountID uint             `json:"account_id"`
}

type ProductMeta struct {
	cruda.BaseModel

	Introduction string                      `json:"introduction"`
	Attachments  datatypes.JSONSlice[string] `json:"attachments"` // random id of attachments

	ProductID uint `json:"product_id"`
}
