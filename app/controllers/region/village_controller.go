package controllers

import "github.com/gofiber/fiber/v2"

func AllVillages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "AllVillages"})
}

func ExportVillages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ExportVillages"})
}

func VillageById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "VillageById", "id": id})
}

func CreateVillage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "CreateVillage"})
}

func ImportVillages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ImportVillages"})
}

func UpdateVillage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "UpdateVillage"})
}

func DeleteVillage(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeleteVillage", "id": id})
}

func TrashAllVillages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "TrashAllVillages"})
}

func RestoreVillage(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RestoreVillage", "id": id})
}

func DeletePermanentVillage(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeletePermanentVillage", "id": id})
}
