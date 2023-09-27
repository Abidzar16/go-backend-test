package nasabahHandler

import (
	"fmt"
	"math/rand"

	"github.com/abidzar16/go-backend-test/database"
	"github.com/abidzar16/go-backend-test/internals/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateFixedLengthStringOfNumber(length int, seed string) string {
	// Create a new rand.Source using the provided seed.
	source := rand.NewSource(int64(len(seed)))
	// Create a new rand.Rand using the rand.Source from step 2.
	random := rand.New(source)

	// Generate a random number between 0 and 9 for each character of the desired length string.
	var stringOfNumbers string
	for i := 0; i < length; i++ {
		randomNumber := random.Intn(10)
		stringOfNumbers += fmt.Sprintf("%d", randomNumber)
	}

	// Return the new string.
	return stringOfNumbers
}

func CreateNasabah(c *fiber.Ctx) error {
	type newNasabah struct {
		Nama  string `json:"nama"`
		Nik   string `json:"nik"`
		No_hp string `json:"no_hp"`
	}
	var newNasabahData newNasabah
	err := c.BodyParser(&newNasabahData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"remark": "Format input tidak tepat", "err": err})
	}

	db := database.DB
	var existedNasabah []model.Nasabah

	db.Find(&existedNasabah, "nik = ?", newNasabahData.Nik)
	if len(existedNasabah) >= 1 {
		return c.Status(400).JSON(fiber.Map{"remark": "NIK sudah terdaftar"})
	}

	db.Find(&existedNasabah, "no_hp = ?", newNasabahData.No_hp)
	if len(existedNasabah) >= 1 {
		return c.Status(400).JSON(fiber.Map{"remark": "No HP sudah terdaftar"})
	}

	var nasabah model.Nasabah
	nasabah.ID = uuid.New()
	nasabah.Nama = newNasabahData.Nama
	nasabah.Nik = newNasabahData.Nik
	nasabah.No_hp = newNasabahData.No_hp
	nasabah.Rekening = generateFixedLengthStringOfNumber(16, newNasabahData.Nik+newNasabahData.No_hp)
	nasabah.Saldo = 0

	err = db.Create(&nasabah).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"remark": "Tidak dapat membuat nasabah baru"})
	}

	// Return the created note
	return c.Status(200).JSON(fiber.Map{"no_rekening": nasabah.Rekening})
}

func TabungSaldo(c *fiber.Ctx) error {
	type newSaldo struct {
		No_rekening string `json:"no_rekening"`
		Nominal     int64  `json:"nominal"`
	}

	var newSaldoData newSaldo
	err := c.BodyParser(&newSaldoData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"remark": "Format input tidak tepat"})
	}

	db := database.DB
	var existedNasabah model.Nasabah

	db.Find(&existedNasabah, "rekening = ?", newSaldoData.No_rekening)
	if existedNasabah.ID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"remark": "No rekening tidak dikenal"})
	}

	existedNasabah.Saldo = existedNasabah.Saldo + newSaldoData.Nominal
	db.Save(&existedNasabah)

	var mutasiEntry model.Transaksi
	mutasiEntry.ID = uuid.New()
	mutasiEntry.Rekening = newSaldoData.No_rekening
	mutasiEntry.Kode_transaksi = "C"
	mutasiEntry.Nominal = newSaldoData.Nominal

	err = db.Create(&mutasiEntry).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"remark": "Tidak dapat membuat entri mutasi baru"})
	}

	return c.Status(200).JSON(fiber.Map{"saldo": existedNasabah.Saldo})
}

func TarikSaldo(c *fiber.Ctx) error {
	type newSaldo struct {
		No_rekening string `json:"no_rekening"`
		Nominal     int64  `json:"nominal"`
	}

	var newSaldoData newSaldo
	err := c.BodyParser(&newSaldoData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"remark": "Format input tidak tepat"})
	}

	db := database.DB
	var existedNasabah model.Nasabah

	db.Find(&existedNasabah, "rekening = ?", newSaldoData.No_rekening)
	if existedNasabah.ID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"remark": "No rekening tidak dikenal"})
	}

	if existedNasabah.Saldo-newSaldoData.Nominal < 0 {
		return c.Status(400).JSON(fiber.Map{"remark": "Saldo anda kurang cukup"})
	}

	existedNasabah.Saldo = existedNasabah.Saldo - newSaldoData.Nominal
	db.Save(&existedNasabah)

	var mutasiEntry model.Transaksi
	mutasiEntry.ID = uuid.New()
	mutasiEntry.Rekening = newSaldoData.No_rekening
	mutasiEntry.Kode_transaksi = "D"
	mutasiEntry.Nominal = newSaldoData.Nominal

	err = db.Create(&mutasiEntry).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"remark": "Tidak dapat membuat entri mutasi baru"})
	}

	// Return the created note
	return c.Status(200).JSON(fiber.Map{"saldo": existedNasabah.Saldo})
}

// /saldo/{no_rekening}
func CekSaldo(c *fiber.Ctx) error {
	no_rekening := c.Params("no_rekening")

	db := database.DB
	var existedNasabah model.Nasabah

	db.Find(&existedNasabah, "Rekening = ?", no_rekening)
	if existedNasabah.ID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"remark": "No rekening tidak dikenal"})
	}

	// Return the created note
	return c.Status(200).JSON(fiber.Map{"saldo": existedNasabah.Saldo})
}

// /mutasi/{no_rekening}
func CekMutasi(c *fiber.Ctx) error {
	no_rekening := c.Params("no_rekening")

	db := database.DB
	var existedNasabah model.Nasabah
	var listTransaksi []model.Transaksi

	db.Find(&existedNasabah, "Rekening = ?", no_rekening)
	if existedNasabah.ID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"remark": "No rekening tidak dikenal"})
	}

	db.Find(&listTransaksi, "Rekening = ?", no_rekening)
	if len(listTransaksi) == 0 {
		return c.Status(400).JSON(fiber.Map{"remark": "Belum ada transaksi"})
	}

	return c.Status(200).JSON(fiber.Map{"mutasi": listTransaksi})
}
