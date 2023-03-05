package main

import (
	"os"
	"quell-api/handlers"
	"quell-api/initializers"
	"quell-api/middlewares"
	"quell-api/repository"
	"quell-api/service"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDb()
	initializers.SyncDatabase()
	initializers.ConnectSupabase()
}

func main() {
	// login, register, validate
	user_Repository := repository.NewUserRepository(initializers.DB)
	user_Service := service.NewUserService(user_Repository)
	user_Handler := handlers.NewUserHandler(user_Service)

	// category, insert ,get, update, delete
	category_Repository := repository.NewCategoryRepository(initializers.DB)
	category_Service := service.NewCategoryService(category_Repository)
	category_Handler := handlers.NewCategoryHandler(category_Service)

	// post, insert, get, update, delete
	post_Repository := repository.NewPostRepository(initializers.DB)
	post_Service := service.NewPostService(post_Repository)
	post_Handler := handlers.NewPostHandler(post_Service)

	// upload, delete
	attachment_Repository := repository.NewAttachmentRepository(initializers.SupabaseClient)
	attachment_Service := service.NewAttachmentService(attachment_Repository)
	attachment_Handler := handlers.NewAttachmentHandler(attachment_Service)

	// payment_Repository := repository.NewPaymentRepository()
	// payment_Service := service.NewPaymentService(payment_Repository)
	// payment_Handler := handlers.NewPaymentHandler(payment_Service)

	// category, insert ,get, update, delete

	router := gin.Default()
	v1 := router.Group("/api/v1")

	v1.POST("/login", user_Handler.LoginHandler)

	v1.POST("/register", user_Handler.RegisterHandler)
	v1.GET("/register/validate", user_Handler.ValidateHandler)
	v1.DELETE("/logout", middlewares.RequireAuth, handlers.LogoutHandler)

	v1.GET("/user", middlewares.RequireAuth, handlers.GetUser)

	v1.GET("/category", category_Handler.GetCategoryHandler)
	v1.GET("/category/:id", category_Handler.GetCategoryByIdHandler)
	v1.POST("/category", middlewares.RequireAuth, category_Handler.CreateCategoryHandler)
	v1.PUT("/category/:id", middlewares.RequireAuth, category_Handler.UpdateCategoryHandler)
	v1.DELETE("/category/:id", middlewares.RequireAuth, category_Handler.DeleteCategoryHandler)

	v1.GET("/posts", post_Handler.GetPostHandler)
	v1.GET("/posts/:id", post_Handler.GetPostByIdHandler)
	v1.POST("/posts", middlewares.RequireAuth, post_Handler.CreatePostHandler)
	v1.PUT("/posts/:id", middlewares.RequireAuth, post_Handler.UpdatePostHandler)
	v1.DELETE("/posts/:id", middlewares.RequireAuth, post_Handler.DeletePostHandler)
	// v1.GET("/qris", middlewares.RequireAuth, payment_Handler.CreateQrisHandler)

	v1.POST("/posts/:id/attachment", middlewares.RequireAuth, attachment_Handler.UploadFile)
	v1.DELETE("/posts/:id/attachment/:attid", middlewares.RequireAuth, attachment_Handler.DeleteFile)

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})

	router.Run(":" + os.Getenv("PORT"))
}
