package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"RWB_L0/config"
	"RWB_L0/internal/adapters/cache"
	"RWB_L0/internal/adapters/postgres"
	httpcontroller "RWB_L0/internal/controllers/http"
	"RWB_L0/internal/controllers/http/v1"
	natscontroller "RWB_L0/internal/controllers/nats"
	"RWB_L0/internal/usecase"
	"RWB_L0/pkg/logger"
	pkgnats "RWB_L0/pkg/nats"
	pkgpostgres "RWB_L0/pkg/postgres"
)

// App - главная структура приложения
type App struct {
	cfg          *config.Config
	log          logger.Logger
	httpServer   *httpcontroller.Server
	natsConsumer *natscontroller.Consumer
	natsClient   *pkgnats.Client
	db           *sql.DB
}

// New - создание приложения
func New() (*App, error) {
	// 1. Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Инициализируем логер
	log := logger.New(cfg.Logging.Level)
	log.Info("Starting Order Service...")
	log.Info("Configuration loaded: server_port=%d, db_host=%s", cfg.Server.Port, cfg.Database.Host)

	return &App{
		cfg: cfg,
		log: log,
	}, nil
}

// Run - запуск приложения
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Инициализируем компоненты
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	if err := a.initNATS(); err != nil {
		return fmt.Errorf("failed to init NATS: %w", err)
	}

	// 2. Создаём Use Cases
	orderRepo := postgres.NewOrderRepository(a.db)
	orderCache := cache.NewMemoryCache()
	orderUseCase := usecase.NewOrderUseCase(orderRepo, orderCache)

	// 3. Восстанавливаем кэш из БД
	a.log.Info("Restoring cache from database...")
	if err := orderUseCase.RestoreCache(ctx); err != nil {
		a.log.Warn("Failed to restore cache: %v", err)
	} else {
		stats := orderUseCase.GetCacheStats()
		a.log.Info("Cache restored successfully: %d orders", stats["cached_orders"])
	}

	// 4. Инициализируем HTTP сервер
	a.initHTTPServer(orderUseCase)

	// 5. Инициализируем NATS consumer
	a.initNATSConsumer(orderUseCase)

	// 6. Запускаем серверы в горутинах
	errChan := make(chan error, 2)

	// HTTP Server
	go func() {
		a.log.Info("Starting HTTP server on %s", a.httpServer.GetAddress())
		if err := a.httpServer.Start(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// NATS Consumer
	go func() {
		if err := a.natsConsumer.Start(ctx, a.cfg.NATS.Subject, a.cfg.NATS.DurableName); err != nil {
			errChan <- fmt.Errorf("NATS consumer error: %w", err)
		}
	}()

	// 7. Ждём сигнала остановки
	a.log.Info("Order Service started successfully!")
	return a.waitForShutdown(ctx, cancel, errChan)
}

// initDatabase - инициализация PostgreSQL
func (a *App) initDatabase() error {
	a.log.Info("Connecting to PostgreSQL: %s:%d", a.cfg.Database.Host, a.cfg.Database.Port)

	db, err := pkgpostgres.New(&pkgpostgres.Config{
		Host:         a.cfg.Database.Host,
		Port:         a.cfg.Database.Port,
		User:         a.cfg.Database.User,
		Password:     a.cfg.Database.Password,
		DBName:       a.cfg.Database.DBName,
		SSLMode:      a.cfg.Database.SSLMode,
		MaxOpenConns: a.cfg.Database.MaxOpenConns,
		MaxIdleConns: a.cfg.Database.MaxIdleConns,
	})
	if err != nil {
		return err
	}

	a.db = db
	a.log.Info("PostgreSQL connected successfully")
	return nil
}

// initNATS - инициализация NATS Streaming
func (a *App) initNATS() error {
	a.log.Info("Connecting to NATS: %s", a.cfg.NATS.URL)

	client, err := pkgnats.New(&pkgnats.Config{
		URL:       a.cfg.NATS.URL,
		ClusterID: a.cfg.NATS.ClusterID,
		ClientID:  a.cfg.NATS.ClientID,
	})
	if err != nil {
		return err
	}

	a.natsClient = client
	a.log.Info("NATS connected successfully")
	return nil
}

// initHTTPServer - инициализация HTTP сервера
func (a *App) initHTTPServer(orderUseCase *usecase.OrderUseCase) {
	// Создаём handlers
	orderHandler := v1.NewOrderHandler(orderUseCase)
	webHandler := v1.NewWebHandler(orderUseCase)

	// Создаём middleware
	mw := httpcontroller.NewMiddleware()

	// Создаём router
	router := httpcontroller.NewRouter(orderHandler, webHandler, mw)

	// Создаём сервер
	a.httpServer = httpcontroller.NewServer(
		a.cfg.Server.Host,
		a.cfg.Server.Port,
		router,
	)
}

// initNATSConsumer - инициализация NATS consumer
func (a *App) initNATSConsumer(orderUseCase *usecase.OrderUseCase) {
	subscriber := pkgnats.NewSubscriber(a.natsClient)
	handler := natscontroller.NewHandler(orderUseCase, a.log)
	a.natsConsumer = natscontroller.NewConsumer(subscriber, handler, a.log)
}

// waitForShutdown - ожидание сигнала остановки
func (a *App) waitForShutdown(_ context.Context, cancel context.CancelFunc, errChan chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		a.log.Error("Application error: %v", err)
		cancel()
		return err
	case sig := <-quit:
		a.log.Info("Received shutdown signal: %v", sig)
		cancel()
	}

	// Graceful shutdown
	return a.shutdown()
}

// shutdown - корректное завершение работы
func (a *App) shutdown() error {
	a.log.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if a.httpServer != nil {
		a.log.Info("Stopping HTTP server...")
		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("HTTP server shutdown error: %v", err)
		}
	}

	// Останавливаем NATS consumer
	if a.natsConsumer != nil {
		a.log.Info("Stopping NATS consumer...")
		if err := a.natsConsumer.Stop(); err != nil {
			a.log.Error("NATS consumer shutdown error: %v", err)
		}
	}

	// Закрываем NATS клиент
	if a.natsClient != nil {
		a.log.Info("Closing NATS connection...")
		if err := a.natsClient.Close(); err != nil {
			a.log.Error("NATS close error: %v", err)
		}
	}

	// Закрываем БД
	if a.db != nil {
		a.log.Info("Closing database connection...")
		if err := a.db.Close(); err != nil {
			a.log.Error("Database close error: %v", err)
		}
	}

	a.log.Info("Shutdown completed successfully")
	return nil
}
