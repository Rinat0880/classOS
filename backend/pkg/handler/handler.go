// pkg/handler/handler.go - обновленные роуты
package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rinat0880/classOS_backend/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Logger())

	router.GET("/ws", h.handleWebSocket)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity, h.adminOnly)
	{
		groups := api.Group("/groups")
		{
			groups.GET("/", h.getAllGroups)
			groups.POST("/", h.createGroup)
			groups.GET("/:id", h.getGroupById)
			groups.PATCH("/:id", h.updateGroup)
			groups.DELETE("/:id", h.deleteGroup)

			users := groups.Group(":id/users")
			{
				users.POST("/", h.createUser)
			}
		}

		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
			users.GET("/:id", h.getUserById)
			users.PATCH("/:id", h.updateUser)
			users.DELETE("/:id", h.deleteUser)
			users.POST("/:id/password", h.changePassword)
		}

		admin := api.Group("/admin")
		{
			admin.POST("/sync", h.syncFromAD)
			admin.GET("/ad/status", h.checkADConnection)
		}

		devices := api.Group("/devices")
		{
			devices.GET("/", h.getAllDevices)
			devices.GET("/online", h.getOnlineDevices)
			devices.GET("/:name", h.getDeviceByName)
			devices.DELETE("/:name", h.deleteDevice)
		}

		logs := api.Group("/logs")
		{
			logs.GET("/", h.getLogs)
			logs.GET("/user/:username", h.getLogsByUsername)
			logs.GET("/device/:device", h.getLogsByDevice)
		}
	}
	return router
}
