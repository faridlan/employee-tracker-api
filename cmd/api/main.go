package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// Sesuaikan module path ini dengan nama module project kamu di go.mod
	"github.com/faridlan/employee-tracker-api/docs"
	"github.com/faridlan/employee-tracker-api/internal/config"
	myHttp "github.com/faridlan/employee-tracker-api/internal/delivery/http"
	"github.com/faridlan/employee-tracker-api/internal/repository/postgres"
	"github.com/faridlan/employee-tracker-api/internal/usecase"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/faridlan/employee-tracker-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title Employee Tracker API
// @version 1.0
// @description Ini adalah dokumentasi API untuk sistem pencatatan target dan performa Karyawan (Employee Tracker).
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}
func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan environment variable dari sistem")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Asumsi kamu memiliki package config untuk setup GORM PostgreSQL sama seperti project sebelumnya
	db := config.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)

	// ==========================================
	// 1. INISIASI REPOSITORY (Layer Data)
	// ==========================================
	employeeRepo := postgres.NewEmployeeRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	productRepo := postgres.NewProductRepository(db)
	targetRepo := postgres.NewTargetRepository(db)
	achievementRepo := postgres.NewAchievementRepository(db)
	notulenRepo := postgres.NewMeetingMinuteRepository(db)

	// ==========================================
	// 2. INISIASI USECASE (Layer Business Logic)
	// ==========================================
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	productUsecase := usecase.NewProductUsecase(productRepo, categoryRepo)
	targetUsecase := usecase.NewTargetUsecase(targetRepo, employeeRepo, productRepo)
	achievementUsecase := usecase.NewAchievementUsecase(achievementRepo, targetRepo)
	notuleUsecase := usecase.NewMeetingMinuteUsecase(notulenRepo)

	// ==========================================
	// 3. INISIASI HANDLER (Layer Delivery)
	// ==========================================
	employeeHandler := myHttp.NewEmployeeHandler(employeeUsecase)
	categoryHandler := myHttp.NewCategoryHandler(categoryUsecase)
	productHandler := myHttp.NewProductHandler(productUsecase)
	targetHandler := myHttp.NewTargetHandler(targetUsecase)
	achievementHandler := myHttp.NewAchievementHandler(achievementUsecase)
	notulenHandler := myHttp.NewMeetingMinuteHandler(notuleUsecase)

	// ==========================================
	// 4. BUNGKUS KE DALAM STRUCT REGISTRY
	// ==========================================
	handlers := myHttp.AppHandlers{
		Employee:      employeeHandler,
		Category:      categoryHandler,
		Product:       productHandler,
		Target:        targetHandler,
		Achievement:   achievementHandler,
		MeetingMinute: notulenHandler,
	}

	// Menjalankan migrasi golang-migrate jika diatur dalam config
	dbURL := os.Getenv("DB_URL")
	if dbURL != "" {
		config.RunDBMigration(dbURL)
	}

	// Setup Swagger Host
	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
	}

	app := fiber.New()

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: false,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// Endpoint Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// ==========================================
	// 5. DAFTARKAN SEMUA ROUTE KE FIBER
	// ==========================================
	myHttp.SetupRoutes(app, handlers)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Menjalankan server dalam Goroutine
	go func() {
		slog.Info("Starting Employee Tracker API Server", slog.String("port", port))
		if err := app.Listen(":" + port); err != nil {
			slog.Error("Server failed to start", slog.String("detail", err.Error()))
		}
	}()

	// Graceful Shutdown Setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Block main thread until signal is received

	slog.Info("Menerima sinyal mati, mematikan server dengan sopan...")

	if err := app.Shutdown(); err != nil {
		slog.Error("Server dipaksa mati karena error", slog.String("detail", err.Error()))
	}

	slog.Info("Employee Tracker API Server berhasil dimatikan dengan aman. Sampai jumpa!")
}
