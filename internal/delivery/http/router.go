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

	// Karyawan
	api.Post("/employees", h.Employee.RegisterEmployee)
	api.Get("/employees/:id", h.Employee.GetEmployeeDetail)
	api.Get("/employees/:employee_id/performance", h.Target.GetEmployeePerformance) // <--- Kalkulasi ada di sini

	// Master Data: Kategori & Produk
	api.Post("/categories", h.Category.CreateCategory)
	api.Get("/categories", h.Category.GetAllCategories)
	api.Post("/products", h.Product.CreateProduct)
	api.Get("/categories/:category_id/products", h.Product.GetProductsByCategory)

	// Siklus Target & Pencapaian
	api.Post("/targets", h.Target.AssignTarget)
	api.Post("/achievements", h.Achievement.RecordAchievement)

}
