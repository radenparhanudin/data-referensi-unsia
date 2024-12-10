package routes

import (
	controllers "data-referensi/app/controllers/education"
	"data-referensi/app/requests"

	"github.com/gofiber/fiber/v2"
)

func EducationRoute(app fiber.Router) {
	educationGroup := app.Group("/education")

	/* Educational Levels */
	educationalLevel := educationGroup.Group("educational-levels")
	educationalLevelTrash := educationalLevel.Group("trashs")
	educationalLevelTrash.Get("/", requests.ValidatePagination, controllers.GetTrashEducationalLevels)
	educationalLevelTrash.Put("/:id", controllers.RestoreEducationalLevel)

	educationalLevel.Get("/", requests.ValidatePagination, controllers.GetEducationalLevels)
	educationalLevel.Get("/export", controllers.ExportEducationalLevels)
	educationalLevel.Get("/search", controllers.SearchEducationalLevels)
	educationalLevel.Get("/:id", controllers.GetEducationalLevel)
	educationalLevel.Post("/", requests.ValidateEducationalLevel, controllers.CreateEducationalLevel)
	educationalLevel.Post("/import", controllers.ImportEducationalLevels)
	educationalLevel.Put("/:id", requests.ValidateEducationalLevel, controllers.UpdateEducationalLevel)
	educationalLevel.Delete("/:id", controllers.DeleteEducationalLevel)

	/* Study Programs */
	studyProgram := educationGroup.Group("study-programs")
	studyProgramTrash := studyProgram.Group("trashs")
	studyProgramTrash.Get("/", requests.ValidatePagination, controllers.GetTrashStudyPrograms)
	studyProgramTrash.Put("/:id", controllers.RestoreStudyProgram)

	studyProgram.Get("/", requests.ValidatePagination, controllers.GetStudyPrograms)
	studyProgram.Get("/export", controllers.ExportStudyPrograms)
	studyProgram.Get("/search", controllers.SearchStudyPrograms)
	studyProgram.Get("/:id", controllers.GetStudyProgram)
	studyProgram.Post("/", requests.ValidateStudyProgram, controllers.CreateStudyProgram)
	studyProgram.Post("/import", controllers.ImportStudyPrograms)
	studyProgram.Put("/:id", requests.ValidateStudyProgram, controllers.UpdateStudyProgram)
	studyProgram.Delete("/:id", controllers.DeleteStudyProgram)

	/* Unsia Study Programs */
	unsiaStudyProgram := educationGroup.Group("unsia-study-programs")
	unsiaStudyProgramTrash := unsiaStudyProgram.Group("trashs")
	unsiaStudyProgramTrash.Get("/", requests.ValidatePagination, controllers.GetTrashUnsiaStudyPrograms)
	unsiaStudyProgramTrash.Put("/:id", controllers.RestoreUnsiaStudyProgram)

	unsiaStudyProgram.Get("/", requests.ValidatePagination, controllers.GetUnsiaStudyPrograms)
	unsiaStudyProgram.Get("/export", controllers.ExportUnsiaStudyPrograms)
	unsiaStudyProgram.Get("/search", controllers.SearchUnsiaStudyPrograms)
	unsiaStudyProgram.Get("/:id", controllers.GetUnsiaStudyProgram)
	unsiaStudyProgram.Post("/", requests.ValidateUnsiaStudyProgram, controllers.CreateUnsiaStudyProgram)
	unsiaStudyProgram.Post("/import", controllers.ImportUnsiaStudyPrograms)
	unsiaStudyProgram.Put("/:id", requests.ValidateUnsiaStudyProgram, controllers.UpdateUnsiaStudyProgram)
	unsiaStudyProgram.Delete("/:id", controllers.DeleteUnsiaStudyProgram)

	/* Educations */
	education := educationGroup.Group("educations")
	educationTrash := education.Group("trashs")
	educationTrash.Get("/", requests.ValidatePagination, controllers.GetTrashEducations)
	educationTrash.Put("/:id", controllers.RestoreEducation)

	education.Get("/", requests.ValidatePagination, controllers.GetEducations)
	education.Get("/export", controllers.ExportEducations)
	education.Get("/search", controllers.SearchEducations)
	education.Get("/by-educational-level/:educatioal_level_id", controllers.GetEducationByEducationalLevelId)
	education.Get("/:id", controllers.GetEducation)
	education.Post("/", requests.ValidateEducation, controllers.CreateEducation)
	education.Post("/import", controllers.ImportEducations)
	education.Put("/:id", requests.ValidateEducation, controllers.UpdateEducation)
	education.Delete("/:id", controllers.DeleteEducation)

}
