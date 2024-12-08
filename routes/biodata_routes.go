package routes

import (
	controllers "data-referensi/app/controllers/biodata"
	"data-referensi/app/requests"

	"github.com/gofiber/fiber/v2"
)

func BiodataRoute(app fiber.Router) {
	biodata := app.Group("/biodata")

	/* Religions */
	religion := biodata.Group("religions")
	religionTrash := religion.Group("trashs")
	religionTrash.Get("/", requests.ValidatePagination, controllers.GetTrashReligions)
	religionTrash.Put("/:id", controllers.RestoreReligion)

	religion.Get("/", requests.ValidatePagination, controllers.GetReligions)
	religion.Get("/export", controllers.ExportReligions)
	religion.Get("/search", controllers.SearchReligions)
	religion.Get("/:id", controllers.GetReligion)
	religion.Post("/", requests.ValidateReligion, controllers.CreateReligion)
	religion.Post("/import", controllers.ImportReligions)
	religion.Put("/:id", requests.ValidateReligion, controllers.UpdateReligion)
	religion.Delete("/:id", controllers.DeleteReligion)

	/* Jobs */
	job := biodata.Group("jobs")
	jobTrash := job.Group("trashs")
	jobTrash.Get("/", requests.ValidatePagination, controllers.GetTrashJobs)
	jobTrash.Put("/:id", controllers.RestoreJob)

	job.Get("/", requests.ValidatePagination, controllers.GetJobs)
	job.Get("/export", controllers.ExportJobs)
	job.Get("/search", controllers.SearchJobs)
	job.Get("/:id", controllers.GetJob)
	job.Post("/", requests.ValidateJob, controllers.CreateJob)
	job.Post("/import", controllers.ImportJobs)
	job.Put("/:id", requests.ValidateJob, controllers.UpdateJob)
	job.Delete("/:id", controllers.DeleteJob)

	/* Ethnics */
	ethnic := biodata.Group("ethnics")
	ethnicTrash := ethnic.Group("trashs")
	ethnicTrash.Get("/", requests.ValidatePagination, controllers.GetTrashEthnics)
	ethnicTrash.Put("/:id", controllers.RestoreEthnic)

	ethnic.Get("/", requests.ValidatePagination, controllers.GetEthnics)
	ethnic.Get("/export", controllers.ExportEthnics)
	ethnic.Get("/search", controllers.SearchEthnics)
	ethnic.Get("/:id", controllers.GetEthnic)
	ethnic.Post("/", requests.ValidateEthnic, controllers.CreateEthnic)
	ethnic.Post("/import", controllers.ImportEthnics)
	ethnic.Put("/:id", requests.ValidateEthnic, controllers.UpdateEthnic)
	ethnic.Delete("/:id", controllers.DeleteEthnic)

	/* Almamater Sizes */
	almamaterSize := biodata.Group("almamater-sizes")
	almamaterSizeTrash := almamaterSize.Group("trashs")
	almamaterSizeTrash.Get("/", requests.ValidatePagination, controllers.GetTrashAlmamaterSizes)
	almamaterSizeTrash.Put("/:id", controllers.RestoreAlmamaterSize)

	almamaterSize.Get("/", requests.ValidatePagination, controllers.GetAlmamaterSizes)
	almamaterSize.Get("/export", controllers.ExportAlmamaterSizes)
	almamaterSize.Get("/search", controllers.SearchAlmamaterSizes)
	almamaterSize.Get("/:id", controllers.GetAlmamaterSize)
	almamaterSize.Post("/", requests.ValidateAlmamaterSize, controllers.CreateAlmamaterSize)
	almamaterSize.Post("/import", controllers.ImportAlmamaterSizes)
	almamaterSize.Put("/:id", requests.ValidateAlmamaterSize, controllers.UpdateAlmamaterSize)
	almamaterSize.Delete("/:id", controllers.DeleteAlmamaterSize)

	/* Marriage Statuses */
	marriageStatus := biodata.Group("marriage-statuses")
	marriageStatusTrash := marriageStatus.Group("trashs")
	marriageStatusTrash.Get("/", requests.ValidatePagination, controllers.GetTrashMarriageStatuses)
	marriageStatusTrash.Put("/:id", controllers.RestoreMarriageStatus)

	marriageStatus.Get("/", requests.ValidatePagination, controllers.GetMarriageStatuses)
	marriageStatus.Get("/export", controllers.ExportMarriageStatuses)
	marriageStatus.Get("/search", controllers.SearchMarriageStatuses)
	marriageStatus.Get("/:id", controllers.GetMarriageStatus)
	marriageStatus.Post("/", requests.ValidateMarriageStatus, controllers.CreateMarriageStatus)
	marriageStatus.Post("/import", controllers.ImportMarriageStatuses)
	marriageStatus.Put("/:id", requests.ValidateMarriageStatus, controllers.UpdateMarriageStatus)
	marriageStatus.Delete("/:id", controllers.DeleteMarriageStatus)

	/* Banks */
	bank := biodata.Group("banks")
	bankTrash := bank.Group("trashs")
	bankTrash.Get("/", requests.ValidatePagination, controllers.GetTrashBanks)
	bankTrash.Put("/:id", controllers.RestoreBank)

	bank.Get("/", requests.ValidatePagination, controllers.GetBanks)
	bank.Get("/export", controllers.ExportBanks)
	bank.Get("/search", controllers.SearchBanks)
	bank.Get("/:id", controllers.GetBank)
	bank.Post("/", requests.ValidateBank, controllers.CreateBank)
	bank.Post("/import", controllers.ImportBanks)
	bank.Put("/:id", requests.ValidateBank, controllers.UpdateBank)
	bank.Delete("/:id", controllers.DeleteBank)
}
