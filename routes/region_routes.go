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
	countryTrash.Get("/", requests.ValidatePagination, controllers.GetTrashCountries)
	countryTrash.Put("/:id", controllers.RestoreCountry)

	country.Get("/", requests.ValidatePagination, controllers.GetCountries)
	country.Get("/export", controllers.ExportCountries)
	country.Get("/search", controllers.SearchCountries)
	country.Get("/:id", controllers.GetCountry)
	country.Post("/", requests.ValidateCountry, controllers.CreateCountry)
	country.Post("/import", controllers.ImportCountries)
	country.Put("/:id", requests.ValidateCountry, controllers.UpdateCountry)
	country.Delete("/:id", controllers.DeleteCountry)

	/* Provinces */
	province := region.Group("provinces")
	provinceTrash := province.Group("trashs")
	provinceTrash.Get("/", requests.ValidatePagination, controllers.GetTrashProvinces)
	provinceTrash.Put("/:id", controllers.RestoreProvince)

	province.Get("/", requests.ValidatePagination, controllers.GetProvinces)
	province.Get("/export", controllers.ExportProvinces)
	province.Get("/search", controllers.SearchProvinces)
	province.Get("/by-country/:country_id", controllers.GetProvinceByCountryId)
	province.Get("/:id", controllers.GetProvince)
	province.Post("/", requests.ValidateProvince, controllers.CreateProvince)
	province.Post("/import", controllers.ImportProvinces)
	province.Put("/:id", requests.ValidateProvince, controllers.UpdateProvince)
	province.Delete("/:id", controllers.DeleteProvince)

	/* Cities */
	city := region.Group("cities")
	cityTrash := city.Group("trashs")
	cityTrash.Get("/", requests.ValidatePagination, controllers.GetTrashCities)
	cityTrash.Put("/:id", controllers.RestoreCity)

	city.Get("/", requests.ValidatePagination, controllers.GetCities)
	city.Get("/export", controllers.ExportCities)
	city.Get("/search", controllers.SearchCities)
	city.Get("/by-province/:province_id", controllers.GetCityByProvinceId)
	city.Get("/:id", controllers.GetCity)
	city.Post("/", requests.ValidateCity, controllers.CreateCity)
	city.Post("/import", controllers.ImportCities)
	city.Put("/:id", requests.ValidateCity, controllers.UpdateCity)
	city.Delete("/:id", controllers.DeleteCity)

	/* Districts */
	district := region.Group("districts")
	districtTrash := district.Group("trashs")
	districtTrash.Get("/", requests.ValidatePagination, controllers.GetTrashDistricts)
	districtTrash.Put("/:id", controllers.RestoreDistrict)

	district.Get("/", requests.ValidatePagination, controllers.GetDistricts)
	district.Get("/export", controllers.ExportDistricts)
	district.Get("/search", controllers.SearchDistricts)
	district.Get("/by-city/:city_id", controllers.GetDistrictByCityId)
	district.Get("/:id", controllers.GetDistrict)
	district.Post("/", requests.ValidateDistrict, controllers.CreateDistrict)
	district.Post("/import", controllers.ImportDistricts)
	district.Put("/:id", requests.ValidateDistrict, controllers.UpdateDistrict)
	district.Delete("/:id", controllers.DeleteDistrict)

	/* Villages */
	village := region.Group("villages")
	villageTrash := village.Group("trashs")
	villageTrash.Get("/", requests.ValidatePagination, controllers.GetTrashVillages)
	villageTrash.Put("/:id", controllers.RestoreVillage)

	village.Get("/", requests.ValidatePagination, controllers.GetVillages)
	village.Get("/export", controllers.ExportVillages)
	village.Get("/search", controllers.SearchVillages)
	village.Get("/:id", controllers.GetVillage)
	village.Get("/by-district/:district_id", controllers.GetVillageByDistrictId)
	village.Post("/", requests.ValidateVillage, controllers.CreateVillage)
	village.Post("/import", controllers.ImportVillages)
	village.Put("/:id", requests.ValidateVillage, controllers.UpdateVillage)
	village.Delete("/:id", controllers.DeleteVillage)
}
