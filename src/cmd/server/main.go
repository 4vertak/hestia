package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/4vertak/hestia/docs" // сгенерированная документация
	"github.com/4vertak/hestia/internal/config"
	"github.com/4vertak/hestia/internal/database"
	"github.com/4vertak/hestia/internal/handler"
	"github.com/4vertak/hestia/internal/handler/middleware"
	"github.com/4vertak/hestia/internal/repository/postgres"
	"github.com/4vertak/hestia/internal/service"
)

// @title HESTIA — Home Equipment Store Transactional Interface API
// @version 1.0
// @description REST API for an appliance store.
// @BasePath /api/v1
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("connected to database")

	// Запуск миграций (путь относительно рабочей директории)
	if err := database.RunMigrations(db, "migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Инициализация репозиториев
	addressRepo := postgres.NewAddressRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	supplierRepo := postgres.NewSupplierRepository(db)
	productRepo := postgres.NewProductRepository(db)
	imageRepo := postgres.NewImageRepository(db)

	// Инициализация сервисов
	clientSvc := service.NewClientService(clientRepo, addressRepo)
	supplierSvc := service.NewSupplierService(supplierRepo, addressRepo)
	productSvc := service.NewProductService(productRepo, imageRepo)
	imageSvc := service.NewImageService(imageRepo, productRepo)

	// Инициализация ручек
	clientH := handler.NewClientHandler(clientSvc)
	supplierH := handler.NewSupplierHandler(supplierSvc)
	productH := handler.NewProductHandler(productSvc)
	imageH := handler.NewImageHandler(imageSvc)

	// Настройка роутера
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа API v1
	v1 := router.Group("/api/v1")
	{
		clients := v1.Group("/clients")
		{
			clients.POST("", clientH.Create)
			clients.DELETE("/:id", clientH.Delete)
			clients.GET("", clientH.SearchByName)
			clients.GET("/all", clientH.GetAll)
			clients.PATCH("/:id/address", clientH.UpdateAddress)
		}
		products := v1.Group("/products")
		{
			products.POST("", productH.Create)
			products.PATCH("/:id/decrease", productH.DecreaseStock)
			products.GET("/:id", productH.GetByID)
			products.GET("", productH.GetAllAvailable)
			products.DELETE("/:id", productH.Delete)
		}
		suppliers := v1.Group("/suppliers")
		{
			suppliers.POST("", supplierH.Create)
			suppliers.PATCH("/:id/address", supplierH.UpdateAddress)
			suppliers.DELETE("/:id", supplierH.Delete)
			suppliers.GET("", supplierH.GetAll)
			suppliers.GET("/:id", supplierH.GetByID)
		}
		images := v1.Group("/images")
		{
			images.POST("", imageH.Upload)
			images.PUT("/:id", imageH.Replace)
			images.DELETE("/:id", imageH.Delete)
			images.GET("/:id", imageH.GetByID)
			images.GET("/product/:productId", imageH.GetByProductID)
		}
	}

	// Сервер с graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		slog.Info("starting server", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server exited")
}