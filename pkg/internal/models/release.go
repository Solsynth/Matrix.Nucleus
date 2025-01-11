package models

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cruda"
	"gorm.io/datatypes"
)

type ProductReleaseType int

const (
	ReleaseTypeFull = ProductReleaseType(iota)
	ReleaseTypePatch
)

type ProductRelease struct {
	cruda.BaseModel

	Version    string                                          `json:"version"`
	Type       ProductReleaseType                              `json:"type"`
	Channel    string                                          `json:"channel"`
	Assets     datatypes.JSONType[map[string]ReleaseAsset]     `json:"assets"`
	Runners    datatypes.JSONType[map[string]ReleaseRunner]    `json:"runners"`
	Installers datatypes.JSONType[map[string]ReleaseInstaller] `json:"installers"`

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

type ReleaseAsset struct {
	URI         string `json:"uri" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

type ReleaseInstaller struct {
	Workdir string                  `json:"workdir"`
	Script  string                  `json:"script"`
	Patches []ReleaseInstallerPatch `json:"patches"`
}

type ReleaseInstallerPatch struct {
	Action string `json:"action" validate:"required"`
	Glob   string `json:"glob" validate:"required"`
}

type ReleaseRunner struct {
	Workdir string `json:"workdir"`
	Script  string `json:"script"`
	Label   string `json:"label"`
}
