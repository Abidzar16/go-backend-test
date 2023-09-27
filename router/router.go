package router

import (
	nasabahRoutes "github.com/abidzar16/go-backend-test/internals/routes/nasabah"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/", logger.New())

	// Setup the Node Routes
	nasabahRoutes.SetupNoteRoutes(api)
}
