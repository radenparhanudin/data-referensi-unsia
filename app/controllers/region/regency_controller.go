package controllers

import "github.com/gofiber/fiber/v2"

func AllRegencies(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "AllRegencies"})
}

func ExportRegencies(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ExportRegencies"})
}

func RegencyById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RegencyById", "id": id})
}

func CreateRegency(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "CreateRegency"})
}

func ImportRegencies(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ImportRegencies"})
}

func UpdateRegency(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "UpdateRegency"})
}

func DeleteRegency(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeleteRegency", "id": id})
}

func TrashAllRegencies(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "TrashAllRegencies"})
}

func RestoreRegency(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RestoreRegency", "id": id})
}

func DeletePermanentRegency(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeletePermanentRegency", "id": id})
}
