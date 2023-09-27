package nasabahRoutes

import (
	nasabahHandler "github.com/abidzar16/go-backend-test/internals/handlers/nasabah"
	"github.com/gofiber/fiber/v2"
)

func SetupNoteRoutes(router fiber.Router) {
	nasabah := router.Group("/")

	nasabah.Post("/daftar", nasabahHandler.CreateNasabah)

	nasabah.Post("/tabung", nasabahHandler.TabungSaldo)

	nasabah.Post("/tarik", nasabahHandler.TarikSaldo)

	nasabah.Get("/saldo/:no_rekening", nasabahHandler.CekSaldo)

	nasabah.Get("/mutasi/:no_rekening", nasabahHandler.CekMutasi)

}
