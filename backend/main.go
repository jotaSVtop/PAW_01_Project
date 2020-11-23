package main

import (
	"projetoapi/model"
	"projetoapi/routes"
	"projetoapi/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var identityKey = "id"

func init() {
	services.OpenDatabase()
	services.Db.AutoMigrate(&model.Worker{})
	services.Db.AutoMigrate(&model.Zone{})
	services.Db.AutoMigrate(&model.WorkersZone{})

	defer services.Db.Close()
}

func main() {

	services.FormatSwagger()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"content-type", "authorization"}
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// AUTH
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	zone := router.Group("/api/zones")
	zone.Use(services.AuthorizationRequired())
	{
		zone.GET("/", routes.GetZones)
		zone.GET("/:id", routes.GetZone)
	}

	admin := router.Group("/api/admin")
	admin.Use(services.AuthorizationRequired())
	{
		admin.GET("/zones", routes.GetZones)
		admin.POST("/zones", routes.AddZone)
		admin.DELETE("/zones", routes.DeleteZone)
		admin.DELETE("/users", routes.DeleteUser)
		admin.POST("/users", routes.Register)
		admin.GET("/users", routes.GetUsers)
	}

	auth := router.Group("/api/")
	{
		auth.POST("/login", routes.GenerateToken)
		auth.PUT("/refresh_token", services.AuthorizationRequired(), routes.RefreshToken)
	}

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8081")
}