package routes

import (
    "github.com/gin-gonic/gin"
    "sanbercode-golang-batch-70-final-project/controllers"
    "sanbercode-golang-batch-70-final-project/middlewares"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // ===============================
    // ✅ SWAGGER - public (tanpa middleware)
    // ===============================
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // ===============================
    // ✅ API Group
    // ===============================
    api := r.Group("/api")
    {
        // ===============================
        // AUTH (tanpa middleware)
        // ===============================
        api.POST("/users/register", controllers.Register)
        api.POST("/users/login", controllers.Login)

        // ===============================
        // LETTERS (user & admin)
        // ===============================
        letters := api.Group("/letters")
        letters.Use(middlewares.AuthMiddleware("")) // "" = semua role boleh
        {
            letters.POST("", controllers.CreateLetter)
            letters.GET("", controllers.GetLetters)
            letters.GET("/:id", controllers.GetLetterByID)
            letters.PUT("/:id", controllers.UpdateLetter)
            letters.DELETE("/:id", controllers.DeleteLetter)
        }

        // ===============================
        // ADMIN (khusus role admin)
        // ===============================
        admin := api.Group("/")
        admin.Use(middlewares.AuthMiddleware("admin"))
        {
            // Users
            admin.POST("/users", controllers.CreateUser)
            admin.GET("/users", controllers.GetUsers)
            admin.GET("/users/:id", controllers.GetUserByID)
            admin.PUT("/users/:id", controllers.UpdateUser)
            admin.DELETE("/users/:id", controllers.DeleteUser)

            // Roles
            admin.POST("/roles", controllers.CreateRole)
            admin.GET("/roles", controllers.GetRoles)
            admin.GET("/roles/:id", controllers.GetRoleByID)
            admin.PUT("/roles/:id", controllers.UpdateRole)
            admin.DELETE("/roles/:id", controllers.DeleteRole)

            // Letter Types
            admin.POST("/letter_types", controllers.CreateLetterType)
            admin.GET("/letter_types", controllers.GetLetterTypes)
            admin.GET("/letter_types/:id", controllers.GetLetterTypeByID)
            admin.PUT("/letter_types/:id", controllers.UpdateLetterType)
            admin.DELETE("/letter_types/:id", controllers.DeleteLetterType)

            // Settings
            admin.POST("/settings", controllers.CreateSetting)
            admin.GET("/settings", controllers.GetSettings)
            admin.GET("/settings/:id", controllers.GetSettingByID)
            admin.PUT("/settings/:id", controllers.UpdateSetting)
            admin.DELETE("/settings/:id", controllers.DeleteSetting)
        }
    }

    return r
}
