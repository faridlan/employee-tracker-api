package http

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

type AppHandlers struct {
	Employee    *EmployeeHandler
	Category    *CategoryHandler
	Product     *ProductHandler
	Target      *TargetHandler
	Achievement *AchievementHandler
}

func SetupRoutes(app *fiber.App, h AppHandlers) {

	prometheus := fiberprometheus.New("employee_tracker_api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	api := app.Group("/api")

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
	api.Patch("/targets/:id/nominal", h.Target.UpdateTargetNominal)
	api.Delete("/targets/:id", h.Target.DeleteTarget)

	// ==========================================
	// 5. Siklus Pencapaian (Achievement)
	// ==========================================
	api.Post("/achievements", h.Achievement.RecordAchievement)
	api.Get("/targets/:target_id/achievements", h.Achievement.GetAchievementsByTarget)
}
