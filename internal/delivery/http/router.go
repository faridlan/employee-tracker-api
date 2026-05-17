package http

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

type AppHandlers struct {
	Employee      *EmployeeHandler
	Category      *CategoryHandler
	Product       *ProductHandler
	Target        *TargetHandler
	Achievement   *AchievementHandler
	MeetingMinute *MeetingMinuteHandler
}

func SetupRoutes(app *fiber.App, h AppHandlers) {

	prometheus := fiberprometheus.New("employee_tracker_api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Hello dari Staging! CI/CD With Selft-Hosted Runner. 🚀",
			"version": "1.0.1-beta",
		})
	})

	// ==========================================
	// 1. Karyawan
	// ==========================================
	api.Post("/employees", h.Employee.RegisterEmployee)
	api.Get("/employees", h.Employee.GetAllEmployee)
	api.Get("/employees/:id", h.Employee.GetEmployeeDetail)
	api.Put("/employees/:id", h.Employee.UpdateEmployee)
	api.Get("/employees/:employee_id/performance", h.Target.GetEmployeePerformance)

	// ==========================================
	// 2. Master Data: Kategori
	// ==========================================
	api.Post("/categories", h.Category.CreateCategory)
	api.Get("/categories", h.Category.GetAllCategories)
	api.Get("/categories/:id", h.Category.GetCategoryDetail)
	api.Put("/categories/:id", h.Category.UpdateCategory)

	// ==========================================
	// 3. Master Data: Produk
	// ==========================================
	api.Post("/products", h.Product.CreateProduct)
	api.Get("/products", h.Product.GetAllProducts)
	api.Get("/products/:id", h.Product.GetProductByID)
	api.Put("/products/:id", h.Product.UpdateProduct)
	api.Get("/categories/:category_id/products", h.Product.GetProductsByCategory)

	// ==========================================
	// 4. Siklus Target
	// ==========================================
	api.Post("/targets", h.Target.AssignTarget)
	api.Get("/targets", h.Target.GetAllTargets)
	api.Patch("/targets/:id/nominal", h.Target.UpdateTargetNominal)
	api.Delete("/targets/:id", h.Target.DeleteTarget)

	// ==========================================
	// 5. Siklus Pencapaian (Achievement)
	// ==========================================
	api.Post("/achievements", h.Achievement.RecordAchievement)
	api.Get("/targets/:target_id/achievements", h.Achievement.GetAchievementsByTarget)

	// ==========================================
	// 6. Siklus Rapat & Notulen (Meeting Minutes)
	// ==========================================
	api.Post("/meetings", h.MeetingMinute.CreateMeeting)
	api.Get("/meetings", h.MeetingMinute.GetAllMeetings)
	api.Get("/meetings/:id", h.MeetingMinute.GetMeetingDetail)
	api.Put("/meetings/:id", h.MeetingMinute.UpdateMeeting)
	api.Patch("/meetings/results/:resultId/status", h.MeetingMinute.UpdateTaskStatus)
}
