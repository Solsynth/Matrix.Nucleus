package api

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/models"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/server/exts"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func listRelease(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	id, _ := c.ParamsInt("productId", 0)

	count, err := services.CountRelease(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListRelease(id, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func calcReleaseToInstall(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("productId", 0)

	var data struct {
		CurrentVersion string `json:"current"`
		TargetVersion  string `json:"target"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	releases, err := services.CalcReleaseToInstall(id, data.CurrentVersion, data.TargetVersion)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(releases)
}

func getRelease(c *fiber.Ctx) error {
	productId, _ := c.ParamsInt("productId", 0)
	id, _ := c.ParamsInt("releaseId", 0)

	if item, err := services.GetReleaseWithProduct(uint(id), uint(productId)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(item)
	}
}

func createRelease(c *fiber.Ctx) error {
	if err := sec.EnsureGrantedPerm(c, "CreateMaReleases", true); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	productId, _ := c.ParamsInt("productId", 0)

	var data struct {
		Version     string                             `json:"version" validate:"required"`
		Type        int                                `json:"type"`
		Channel     string                             `json:"channel" validate:"required"`
		Title       string                             `json:"title" validate:"required,max=1024"`
		Description string                             `json:"description" validate:"required,max=4096"`
		Content     string                             `json:"content" validate:"required"`
		Assets      map[string]models.ReleaseAsset     `json:"assets" validate:"required"`
		Installers  map[string]models.ReleaseInstaller `json:"installers" validate:"required"`
		Runners     map[string]models.ReleaseRunner    `json:"runners" validate:"required"`
		Attachments []string                           `json:"attachments"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	product, err := services.GetProductWithUser(uint(productId), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	release := models.ProductRelease{
		Version:    data.Version,
		Type:       models.ProductReleaseType(data.Type),
		Channel:    data.Channel,
		Assets:     datatypes.NewJSONType(data.Assets),
		Installers: datatypes.NewJSONType(data.Installers),
		Runners:    datatypes.NewJSONType(data.Runners),
		ProductID:  product.ID,
		Meta: models.ProductReleaseMeta{
			Title:       data.Title,
			Description: data.Description,
			Content:     data.Content,
			Attachments: data.Attachments,
		},
	}

	if release, err := services.NewRelease(release); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(release)
	}
}

func updateRelease(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	productId, _ := c.ParamsInt("productId", 0)
	id, _ := c.ParamsInt("releaseId", 0)

	var data struct {
		Version     string                             `json:"version" validate:"required"`
		Type        int                                `json:"type"`
		Channel     string                             `json:"channel" validate:"required"`
		Title       string                             `json:"title" validate:"required,max=1024"`
		Description string                             `json:"description" validate:"required,max=4096"`
		Content     string                             `json:"content" validate:"required"`
		Assets      map[string]models.ReleaseAsset     `json:"assets" validate:"required"`
		Installers  map[string]models.ReleaseInstaller `json:"installers" validate:"required"`
		Runners     map[string]models.ReleaseRunner    `json:"runners" validate:"required"`
		Attachments []string                           `json:"attachments"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	product, err := services.GetProductWithUser(uint(productId), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	release, err := services.GetReleaseWithProduct(uint(id), product.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	release.Version = data.Version
	release.Type = models.ProductReleaseType(data.Type)
	release.Channel = data.Channel
	release.Assets = datatypes.NewJSONType(data.Assets)
	release.Installers = datatypes.NewJSONType(data.Installers)
	release.Runners = datatypes.NewJSONType(data.Runners)
	release.Meta.Title = data.Title
	release.Meta.Description = data.Description
	release.Meta.Content = data.Content
	release.Meta.Attachments = data.Attachments

	if release, err := services.UpdateRelease(release); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(release)
	}
}

func deleteRelease(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	productId, _ := c.ParamsInt("productId", 0)
	id, _ := c.ParamsInt("releaseId", 0)

	product, err := services.GetProductWithUser(uint(productId), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	release, err := services.GetReleaseWithProduct(uint(id), product.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if _, err := services.DeleteRelease(release); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
