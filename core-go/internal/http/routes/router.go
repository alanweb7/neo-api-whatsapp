package routes

import (
	"net/http"

	"github.com/alan/baileys-saas/core-go/internal/http/controllers"
	"github.com/alan/baileys-saas/core-go/internal/http/middleware"
	"github.com/alan/baileys-saas/core-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth    *controllers.AuthController
	Tenant  *controllers.TenantController
	User    *controllers.UserController
	APIKey  *controllers.APIKeyController
	Session *controllers.SessionController
	Message *controllers.MessageController
	Webhook *controllers.WebhookController
}

func Build(tokens *service.TokenService, c Controllers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())

	r.GET("/healthz", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "ready"}) })

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/login", c.Auth.Login)
		auth.POST("/refresh", c.Auth.Refresh)
		auth.GET("/me", middleware.Auth(tokens), c.Auth.Me)

		protected := v1.Group("")
		protected.Use(middleware.Auth(tokens))
		{
			protected.POST("/tenants", c.Tenant.Create)
			protected.GET("/tenants", c.Tenant.List)
			protected.GET("/tenants/:tenantId", c.Tenant.Get)
			protected.PUT("/tenants/:tenantId", c.Tenant.Update)

			protected.POST("/users", c.User.Create)
			protected.GET("/users", c.User.List)
			protected.POST("/users/attach", c.User.Attach)

			protected.POST("/api-keys", c.APIKey.Create)
			protected.GET("/api-keys", c.APIKey.List)
			protected.POST("/api-keys/:apiKeyId/revoke", c.APIKey.Revoke)

			protected.POST("/sessions", c.Session.Create)
			protected.POST("/sessions/:sessionId/start", c.Session.Start)
			protected.GET("/sessions", c.Session.List)
			protected.GET("/sessions/:sessionId", c.Session.Get)
			protected.GET("/sessions/:sessionId/qr", c.Session.QR)
			protected.GET("/sessions/:sessionId/status", c.Session.Status)
			protected.POST("/sessions/:sessionId/reconnect", c.Session.Reconnect)
			protected.POST("/sessions/:sessionId/disconnect", c.Session.Disconnect)
			protected.DELETE("/sessions/:sessionId", c.Session.Remove)

			protected.POST("/messages/text", c.Message.SendText)
			protected.POST("/messages/image", c.Message.SendImage)
			protected.POST("/messages/document", c.Message.SendDocument)
			protected.POST("/messages/audio", c.Message.SendAudio)
			protected.GET("/messages/logs", c.Message.ListLogs)

			protected.POST("/webhooks", c.Webhook.Create)
			protected.GET("/webhooks", c.Webhook.List)
			protected.PUT("/webhooks/:webhookId", c.Webhook.Update)
			protected.DELETE("/webhooks/:webhookId", c.Webhook.Delete)
			protected.GET("/webhooks/deliveries", c.Webhook.ListDeliveries)
		}
	}

	return r
}
