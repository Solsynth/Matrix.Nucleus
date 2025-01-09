package models

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cruda"
	"gorm.io/datatypes"
)

type ProductReleaseType int

const (
	ReleaseTypeMinor = ProductReleaseType(iota)
	ReleaseTypeRegular
	ReleaseTypeMajor
)

type ProductRelease struct {
	cruda.BaseModel

	Version string                             `json:"version"`
	Type    ProductReleaseType                 `json:"type"`
	Channel string                             `json:"channel"`
	Assets  datatypes.JSONType[map[string]any] `json:"assets"`

	ProductID uint               `json:"product_id"`
	Meta      ProductReleaseMeta `json:"meta" gorm:"foreignKey:ReleaseID"`
}

type ProductReleaseMeta struct {
	cruda.BaseModel

	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Content     string                      `json:"content"`
	Attachments datatypes.JSONSlice[string] `json:"attachments"`
	ReleaseID   uint                        `json:"release_id"`
}
