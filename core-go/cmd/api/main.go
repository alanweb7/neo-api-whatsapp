package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alan/baileys-saas/core-go/internal/config"
	"github.com/alan/baileys-saas/core-go/internal/domain"
	"github.com/alan/baileys-saas/core-go/internal/events"
	"github.com/alan/baileys-saas/core-go/internal/http/controllers"
	"github.com/alan/baileys-saas/core-go/internal/http/routes"
	dbinfra "github.com/alan/baileys-saas/core-go/internal/infra/db"
	"github.com/alan/baileys-saas/core-go/internal/infra/engineclient"
	redisinfra "github.com/alan/baileys-saas/core-go/internal/infra/redis"
	"github.com/alan/baileys-saas/core-go/internal/repository"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/alan/baileys-saas/core-go/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)
	db, err := dbinfra.Connect(cfg.DBURL, cfg.Env == "production")
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&domain.Plan{},
		&domain.Tenant{},
		&domain.User{},
		&domain.TenantUser{},
		&domain.ApiKey{},
		&domain.WhatsAppSession{},
		&domain.WebhookEndpoint{},
		&domain.WebhookDelivery{},
		&domain.MessageLog{},
		&domain.AuditLog{},
	); err != nil {
		panic(err)
	}

	redisClient, err := redisinfra.Connect(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		panic(err)
	}

	tokenSvc := service.NewTokenService(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.JWTAccessTTLMin, cfg.JWTRefreshTTLDays)
	engine := engineclient.New(cfg.EngineBaseURL, cfg.InternalAPIKey)

	tenantRepo := repository.NewTenantRepository(db)
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	authSvc := service.NewAuthService(userRepo, tokenSvc)
	tenantSvc := service.NewTenantService(tenantRepo)
	userSvc := service.NewUserService(userRepo)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo)
	sessionSvc := service.NewSessionService(sessionRepo, engine)
	messageSvc := service.NewMessageService(messageRepo, sessionRepo, engine)
	webhookSvc := service.NewWebhookService(webhookRepo)

	authController := controllers.NewAuthController(authSvc, userRepo)
	tenantController := controllers.NewTenantController(tenantSvc)
	userController := controllers.NewUserController(userSvc)
	apiKeyController := controllers.NewAPIKeyController(apiKeySvc)
	sessionController := controllers.NewSessionController(sessionSvc)
	messageController := controllers.NewMessageController(messageSvc)
	webhookController := controllers.NewWebhookController(webhookSvc)

	r := routes.Build(tokenSvc, routes.Controllers{
		Auth:    authController,
		Tenant:  tenantController,
		User:    userController,
		APIKey:  apiKeyController,
		Session: sessionController,
		Message: messageController,
		Webhook: webhookController,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eventConsumer := events.NewConsumer(redisClient, sessionRepo, messageRepo, log)
	eventConsumer.Start(ctx)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: r,
	}

	go func() {
		log.WithField("port", cfg.HTTPPort).Info("core api started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("server failed")
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = redisClient.Close()

	log.Info("core api stopped")
}
