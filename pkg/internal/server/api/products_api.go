package api

import (
	"strconv"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/models"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/server/exts"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listProduct(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	count, err := services.CountProduct()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListProduct(take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listCreatedProduct(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	count, err := services.CountCreatedProduct(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListCreatedProduct(user.ID, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func getProduct(c *fiber.Ctx) error {
	alias := c.Params("productId")

	var item models.Product
	if numericId, err := strconv.Atoi(alias); err == nil {
		item, err = services.GetProduct(uint(numericId))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		item, err = services.GetProduct(uint(numericId))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(item)
}

func createProduct(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Name         string   `json:"name" validate:"required,max=256"`
		Description  string   `json:"description" validate:"max=4096"`
		Introduction string   `json:"introduction"`
		Alias        string   `json:"alias" validate:"required"`
		Tags         []string `json:"tags"`
		Attachments  []string `json:"attachments"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	product := models.Product{
		Name:        data.Name,
		Alias:       data.Alias,
		Description: data.Description,
		Tags:        data.Tags,
		Meta: models.ProductMeta{
			Introduction: data.Introduction,
			Attachments:  data.Attachments,
		},
		AccountID: user.ID,
	}

	if product, err := services.NewProduct(product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(product)
	}
}

func updateProduct(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("productId", 0)

	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	var data struct {
		Icon         string   `json:"icon"`
		Name         string   `json:"name" validate:"required,max=256"`
		Description  string   `json:"description" validate:"max=4096"`
		Introduction string   `json:"introduction"`
		Alias        string   `json:"alias" validate:"required"`
		Tags         []string `json:"tags"`
		Previews     []string `json:"previews"`
		Attachments  []string `json:"attachments"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	product, err := services.GetProductWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	product.Icon = data.Icon
	product.Name = data.Name
	product.Description = data.Description
	product.Tags = data.Tags
	product.Previews = data.Previews
	product.Meta.Introduction = data.Introduction
	product.Meta.Attachments = data.Attachments

	if product, err := services.UpdateProduct(product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(product)
	}
}

func deleteProduct(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("productId", 0)

	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("nex_user").(*sec.UserInfo)

	product, err := services.GetProductWithUser(uint(id), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if _, err := services.DeleteProduct(product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
