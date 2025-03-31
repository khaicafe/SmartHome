package routes

import (
	"go-react-app/controllers"
	"go-react-app/middlewares"
	"go-react-app/models"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", models.DB)
		c.Next()
	})

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middlewares.AuthMiddleware())

	var user models.User
	if err := models.DB.Where("role = ?", "admin").First(&user).Error; err != nil {
		r.POST("/api/setup", controllers.InitialSetup)
	}

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// user system
		api.POST("/auth/signup", controllers.Signup)
		api.POST("/auth/verify-signup-otp", controllers.VerifySignupOTP)
		api.POST("/auth/send-otp", controllers.SendOTP)
		api.POST("/auth/resend-otp", controllers.ResendOTP)
		api.POST("/auth/login", controllers.Login)
		api.POST("/auth/reset-password", controllers.ResetPassword)

		// api.GET("/device/on", controllers.TurnOn)
		// api.GET("/device/off", controllers.TurnOff)
		api.GET("/devices", controllers.GetDevices)
		api.GET("/device/:id/functions", controllers.GetDeviceFunctions)
		api.POST("/device/:id/command", controllers.ActionSendCommand)

		api.POST("/map-switch", controllers.MapSwitch)
		api.GET("/mapped-switches", controllers.GetMappedSwitches)
		api.PUT("/map-switch/:id", controllers.UpdateMappedSwitch)
		api.DELETE("/map-switch/:id", controllers.DeleteMappedSwitch)
		api.POST("/reset-switch-state", controllers.ResetSwitchStateHandler)

		api.GET("/settings", controllers.GetSettings)
		api.POST("/settings", controllers.UpdateSetting)

	}

	// Serve frontend build (React)
	r.Static("/assets", filepath.Join("..", "frontend", "dist", "assets"))
	r.StaticFile("/vite.svg", filepath.Join("..", "frontend", "dist", "vite.svg"))
	r.StaticFile("/", filepath.Join("..", "frontend", "dist", "index.html"))

	// Fallback to index.html for SPA routes
	r.NoRoute(func(c *gin.Context) {
		uri := c.Request.RequestURI
		if filepath.Ext(uri) == "" {
			c.File(filepath.Join("..", "frontend", "dist", "index.html"))
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	return r
}
