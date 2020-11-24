package main

import (
	"projetoapi/model"
	"projetoapi/routes"
	"projetoapi/services"

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

	var admin model.Worker
	admin.Username = "Admin"
	admin.Name = "Test Admin Account"
	admin.Password = "admin123"
	admin.Admin = true
	services.Db.Create(&admin)

	var user model.Worker
	user.Username = "TestUser"
	user.Name = "Test User Account"
	user.Password = "user123"
	user.Admin = false
	services.Db.Create(&user)

	defer services.Db.Close()
}

func main() {

	services.FormatSwagger()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// AUTH
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	zone := router.Group("/api/zones")
	zone.Use(services.AuthorizationRequired())
	{
		zone.GET("/all", routes.GetZones)
		zone.GET("/worker", routes.GetWorkerZones)
		zone.GET("/id/:id", routes.GetZone)
		zone.POST("/id/:id/add", routes.AddPerson)
		zone.POST("/id/:id/remove", routes.RemovePerson)
	}

	worker := router.Group("/api/workers")
	worker.Use(services.AuthorizationRequired())
	{
		//worker.GET("/:id", routes.GetWorkerZoneByID)
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
