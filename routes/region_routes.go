package routes

import (
	controllers "data-referensi/app/controllers/region"
	"data-referensi/app/requests"

	"github.com/gofiber/fiber/v2"
)

func RegionRoute(app fiber.Router) {
	region := app.Group("/region")

	/* Countries */
	country := region.Group("countries")
	countryTrash := country.Group("trashs")
	countryTrash.Get("/", controllers.TrashAllCountries)
	countryTrash.Put("/:id", controllers.RestoreCountry)
	countryTrash.Delete("/:id", controllers.DeletePermanentCountry)

	country.Get("/", requests.ValidatePagination, controllers.AllCountries)
	country.Get("/export", controllers.ExportCountries)
	country.Get("/:id", controllers.CountryById)
	country.Post("/", controllers.CreateCountry)
	country.Post("/import", controllers.ImportCountries)
	country.Put("/:id", controllers.UpdateCountry)
	country.Delete("/:id", controllers.DeleteCountry)

	/* Provinces */
	province := region.Group("provinces")
	provinceTrash := province.Group("trashs")
	provinceTrash.Get("/", controllers.TrashAllProvinces)
	provinceTrash.Put("/:id", controllers.RestoreProvince)
	provinceTrash.Delete("/:id", controllers.DeletePermanentProvince)

	province.Get("/", controllers.AllProvinces)
	province.Get("/export", controllers.ExportProvinces)
	province.Get("/:id", controllers.ProvinceById)
	province.Post("/", controllers.CreateProvince)
	province.Post("/import", controllers.ImportProvinces)
	province.Put("/:id", controllers.UpdateProvince)
	province.Delete("/:id", controllers.DeleteProvince)

	/* Regencies */
	regency := region.Group("regencies")
	regencyTrash := regency.Group("trashs")
	regencyTrash.Get("/", controllers.TrashAllRegencies)
	regencyTrash.Put("/:id", controllers.RestoreRegency)
	regencyTrash.Delete("/:id", controllers.DeletePermanentRegency)

	regency.Get("/", controllers.AllRegencies)
	regency.Get("/export", controllers.ExportRegencies)
	regency.Get("/:id", controllers.RegencyById)
	regency.Post("/", controllers.CreateRegency)
	regency.Post("/import", controllers.ImportRegencies)
	regency.Put("/:id", controllers.UpdateRegency)
	regency.Delete("/:id", controllers.DeleteRegency)

	/* Districts */
	district := region.Group("districts")
	districtTrash := district.Group("trashs")
	districtTrash.Get("/", controllers.TrashAllDistricts)
	districtTrash.Put("/:id", controllers.RestoreDistrict)
	districtTrash.Delete("/:id", controllers.DeletePermanentDistrict)

	district.Get("/", controllers.AllDistricts)
	district.Get("/export", controllers.ExportDistricts)
	district.Get("/:id", controllers.DistrictById)
	district.Post("/", controllers.CreateDistrict)
	district.Post("/import", controllers.ImportDistricts)
	district.Put("/:id", controllers.UpdateDistrict)
	district.Delete("/:id", controllers.DeleteDistrict)

	/* Villages */
	village := region.Group("villages")
	villageTrash := village.Group("trashs")
	villageTrash.Get("/", controllers.TrashAllVillages)
	villageTrash.Put("/:id", controllers.RestoreVillage)
	villageTrash.Delete("/:id", controllers.DeletePermanentVillage)

	village.Get("/", controllers.AllVillages)
	village.Get("/export", controllers.ExportVillages)
	village.Get("/:id", controllers.VillageById)
	village.Post("/", controllers.CreateVillage)
	village.Post("/import", controllers.ImportVillages)
	village.Put("/:id", controllers.UpdateVillage)
	village.Delete("/:id", controllers.DeleteVillage)
}
