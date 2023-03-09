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
	// initializers.LoadEnvVar()
	initializers.ConnectDb()
	initializers.SyncDatabase()
	initializers.ConnectSupabase()
}

func main() {
	user_Repository := repository.NewUserRepository(initializers.DB)
	user_Service := service.NewUserService(user_Repository)
	user_Handler := handlers.NewUserHandler(user_Service)

	category_Repository := repository.NewCategoryRepository(initializers.DB)
	category_Service := service.NewCategoryService(category_Repository)
	category_Handler := handlers.NewCategoryHandler(category_Service)

	post_Repository := repository.NewPostRepository(initializers.DB)
	post_Service := service.NewPostService(post_Repository)
	post_Handler := handlers.NewPostHandler(post_Service)

	attachment_Repository := repository.NewAttachmentRepository(initializers.SupabaseClient)
	attachment_Service := service.NewAttachmentService(attachment_Repository)
	attachment_Handler := handlers.NewAttachmentHandler(attachment_Service)

	saving_Repository := repository.NewSavingRepository(initializers.DB)
	saving_Service := service.NewSavingService(saving_Repository)
	saving_Handler := handlers.NewSavingHandler(saving_Service)

	saving_Category_Repository := repository.NewSavingCategoryRepository(initializers.DB)
	saving_Category_Service := service.NewSavingCategoryService(saving_Category_Repository)
	saving_Category_Handler := handlers.NewSavingCategoryHandler(saving_Category_Service)

	payment_Repository := repository.NewPaymentRepository(initializers.DB, "gojek")
	payment_Service := service.NewPaymentService(payment_Repository)
	payment_Handler := handlers.NewPaymentHandler(payment_Service)

	// category, insert ,get, update, delete

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/api/v1")

	v1.POST("/login", user_Handler.LoginHandler)

	v1.POST("/register", user_Handler.RegisterHandler)
	v1.GET("/register/validate", user_Handler.ValidateHandler)
	v1.DELETE("/logout", middlewares.RequireAuth, handlers.LogoutHandler)

	v1.GET("/user", middlewares.RequireAuth, handlers.GetUser)
	v1.GET("/user/subscribe", middlewares.RequireAuth, payment_Handler.PremiumPayment)
	v1.POST("/user/subscribe/validate", payment_Handler.PremiumPaymentValidate)

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

	v1.POST("/posts/:id/attachment", middlewares.RequireAuth, attachment_Handler.UploadFile)
	v1.DELETE("/posts/:id/attachment/:attid", middlewares.RequireAuth, attachment_Handler.DeleteFile)

	v1.GET("/saving", middlewares.RequireAuth, saving_Handler.GetSavingHandler)
	v1.GET("/saving/:id", middlewares.RequireAuth, saving_Handler.GetSavingByIdHandler)
	v1.GET("/saving/totalamount", middlewares.RequireAuth, saving_Handler.GetTotalAmountHandler)
	v1.POST("/saving", middlewares.RequireAuth, saving_Handler.CreateSavingHandler)
	v1.PUT("/saving/:id", middlewares.RequireAuth, saving_Handler.UpdateSavingHandler)
	v1.DELETE("/saving/:id", middlewares.RequireAuth, saving_Handler.DeleteSavingHandler)

	v1.GET("/savingcategory", middlewares.RequireAuth, saving_Category_Handler.GetSavingCategoryHandler)
	v1.GET("/savingcategory/:id", middlewares.RequireAuth, saving_Category_Handler.GetSavingCategoryByIdHandler)
	v1.POST("/savingcategory", middlewares.RequireAuth, saving_Category_Handler.CreateSavingCategoryHandler)
	v1.PUT("/savingcategory/:id", middlewares.RequireAuth, saving_Category_Handler.UpdateSavingCategoryHandler)
	v1.DELETE("/savingcategory/:id", middlewares.RequireAuth, saving_Category_Handler.DeleteSavingCategoryHandler)

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
