package api

import (
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App, baseURL string) {
	api := app.Group(baseURL).Name("API")
	{
		products := api.Group("/products")
		{
			products.Get("/", listProduct)
			products.Get("/created", listCreatedProduct)
			products.Get("/:productId", getProduct)
			products.Post("/", createProduct)
			products.Put("/:productId", updateProduct)
			products.Delete("/:productId", deleteProduct)

			releases := products.Group("/:productId/releases")
			{
				releases.Get("/", listRelease)
				releases.Post("/calc", calcReleaseToInstall)
				releases.Get("/:releaseId", getRelease)
				releases.Post("/", createRelease)
				releases.Put("/:releaseId", updateRelease)
				releases.Delete("/:releaseId", deleteRelease)
			}
		}
	}
}
