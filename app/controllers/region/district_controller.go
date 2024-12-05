package controllers

import "github.com/gofiber/fiber/v2"

func AllDistricts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "AllDistricts"})
}

func ExportDistricts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ExportDistricts"})
}

func DistrictById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DistrictById", "id": id})
}

func CreateDistrict(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "CreateDistrict"})
}

func ImportDistricts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ImportDistricts"})
}

func UpdateDistrict(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "UpdateDistrict"})
}

func DeleteDistrict(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeleteDistrict", "id": id})
}

func TrashAllDistricts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "TrashAllDistricts"})
}

func RestoreDistrict(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RestoreDistrict", "id": id})
}

func DeletePermanentDistrict(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeletePermanentDistrict", "id": id})
}
