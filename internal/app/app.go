package app

import (
	"context"
	"cqs-kanban/config"
	"cqs-kanban/database"
	"cqs-kanban/internal/handler"
	"cqs-kanban/internal/logger"
	"cqs-kanban/internal/middleware"
	"cqs-kanban/internal/repository"
	"cqs-kanban/internal/service"

	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	flogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// App represents the application
type App struct {
	config *config.Config
	fiber  *fiber.App
	db     database.Database
	// Handlers
	handlers []handler.BaseHandler

	authService      service.AuthService
	operationService service.OperationService
	menuService      service.MenuService
	forecastService  service.ForecastService
}

// New creates a new application instance
func New(cfg *config.Config, db database.Database) *App {
	app := &App{
		config: cfg,
		db:     db,
	}

	// Initialize Fiber
	app.fiber = fiber.New(fiber.Config{
		AppName:      cfg.Server.Name,
		ErrorHandler: errorHandler,
	})

	// Setup middleware
	app.fiber.Use(recover.New())
	app.fiber.Use(flogger.New())
	app.fiber.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))

	// Report
	ctx := context.Background()
	baseErpRepo := repository.NewBaseERPRepository(ctx, app.db.ERPDB())
	reportRepo := repository.NewReportRepo(app.db.DB())
	operationRepo := repository.NewOperationRepo(app.db.DB())
	logger := logger.NewConsoleLogger()
	reportService := service.NewReportService(reportRepo, baseErpRepo, operationRepo, logger)
	reportHandler := handler.NewReportHandler(reportService)
	// Menu
	menuRepo := repository.NewMenuRepo(app.db.DB())
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	// Department
	departmentRepo := repository.NewDepartmentRepo(app.db.DB())
	departmentService := service.NewDepartmentService(departmentRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	// Users
	userRepo := repository.NewUserRepo(app.db.DB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	// SalesCopi04

	saleCopi04Repo := repository.NewSaleCopi04(app.db.ERPDB())
	saleCopi04Service := service.NewSaleCopi04Service(saleCopi04Repo)
	saleCopi04Handler := handler.NewSaleCopi04Handler(saleCopi04Service)
	// Copma
	copmaRepo := repository.NewCopmaRepo(app.db.ERPDB())
	copmaService := service.NewCopmaService(copmaRepo)
	copmaHandler := handler.NewCopmaHandler(copmaService)
	// Forecast
		forecasrRepo := repository.NewForecastRepo(app.db.ERPDB())
	forecastService := service.NewForecastService(forecasrRepo, operationRepo, logger)
	forecastHandler := handler.NewForecastHandler(forecastService)
	authService := service.NewAuthService(userRepo, app.config)
	authHandler := handler.NewAuthHandler(authService)
	adminService := service.NewAdminService(operationRepo, userRepo, departmentRepo, reportRepo, logger)
	adminHandler := handler.NewAdminHandler(adminService)
	app.authService = authService
	app.handlers = []handler.BaseHandler{
		reportHandler,
		menuHandler,
		departmentHandler,
		userHandler,
		authHandler,
		adminHandler,
		saleCopi04Handler,
		copmaHandler,
		forecastHandler,
	}

	return app
}

// SetupRoutes configures the application routes
func (a *App) SetupRoutes() {
	// Health check endpoint
	a.fiber.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"name":   a.config.Server.Name,
			"env":    a.config.Server.Env,
		})
	})

	// API routes
	api := a.fiber.Group("/api")
	// White list routes
	whitelist := []string{
		"/api/auth/login",
		"api/forecasts",
		"/api/admin/dashboard",
		"/api/admin/dashboard/access-trend",
	}

	// Protected routes
	protected := api.Group("/", middleware.JWTMiddleware(a.authService, whitelist))

	// Setup all handler routes
	for _, handler := range a.handlers {
		handler.SetupRoutes(protected)
	}
	// 404 handler
	a.fiber.Use(func(c fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Not Found",
			"error":   "The requested resource does not exist",
		})
	})
}

func (a *App) Start() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", a.config.Server.Port)
		if err := a.fiber.Listen(addr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", a.config.Server.Port)
	<-sigChan
	log.Println("Shutting down server...")

	if err := a.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	if err := a.fiber.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}
