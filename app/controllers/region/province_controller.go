package controllers

import "github.com/gofiber/fiber/v2"

func AllProvinces(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "AllProvinces"})
}

func ExportProvinces(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ExportProvinces"})
}

func ProvinceById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "ProvinceById", "id": id})
}

func CreateProvince(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "CreateProvince"})
}

func ImportProvinces(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ImportProvinces"})
}

func UpdateProvince(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "UpdateProvince"})
}

func DeleteProvince(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeleteProvince", "id": id})
}

func TrashAllProvinces(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "TrashAllProvinces"})
}

func RestoreProvince(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RestoreProvince", "id": id})
}

func DeletePermanentProvince(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeletePermanentProvince", "id": id})
}
