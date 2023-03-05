package main

import (
	"fmt"
	"log"
	"os"

	"bmacharia/aws_sam_blog_api/controller"
	"bmacharia/aws_sam_blog_api/database"
	"bmacharia/aws_sam_blog_api/model"
	"bmacharia/aws_sam_blog_api/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	serveApplication()
}

func loadDatabase() {
	database.InitDb()
	database.Db.AutoMigrate(&model.Role{})
	database.Db.AutoMigrate(&model.User{})
	database.Db.AutoMigrate(&model.Article{})
	seedData()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func seedData() {
	var roles = []model.Role{{Name: "admin", Description: "Administrator role"}, {Name: "author", Description: "Author role"}, {Name: "anonymous", Description: "Unauthenticated user role"}}
	var user = []model.User{{Username: os.Getenv("ADMIN_USERNAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), RoleID: 1}}
	database.Db.Save(&roles)
	database.Db.Save(&user)
}

func serveApplication() {
	router := gin.Default()

	authRoutes := router.Group("/auth/user")
	authRoutes.POST("/register", controller.Register)
	authRoutes.POST("/login", controller.Login)

	publicRoutes := router.Group("/api/view")
	publicRoutes.GET("/articles", controller.GetArticles)
	publicRoutes.GET("/article/:id", controller.GetArticle)

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(util.JWTAuth())
	adminRoutes.GET("/users", controller.GetUsers)
	adminRoutes.GET("/user/:id", controller.GetUser)
	adminRoutes.PUT("/user/:id", controller.UpdateUser)
	adminRoutes.POST("/user/role", controller.CreateRole)
	adminRoutes.GET("/user/roles", controller.GetRoles)
	adminRoutes.PUT("/user/role/:id", controller.UpdateRole)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(util.JWTAuthAuthor())
	protectedRoutes.POST("/article/create", controller.CreateArticle)
	protectedRoutes.PUT("/article/:id", controller.UpdateArticle)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}
