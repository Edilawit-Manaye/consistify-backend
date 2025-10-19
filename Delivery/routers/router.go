// package routers

// import (
// 	"consistent_1/Delivery/controllers"
// 	"consistent_1/Delivery/middleware"
// 	"consistent_1/Infrastructure/auth"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )


// func SetupRouter(
// 	userController *controllers.UserController,
// 	consistencyController *controllers.ConsistencyController,
// 	jwtService auth.JWTService,
// ) *gin.Engine {
// 	router := gin.Default()

	
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"} // React/Flutter web defaults or your specific frontend URL
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
// 	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
// 	config.ExposeHeaders = []string{"Content-Length"}
// 	config.AllowCredentials = true
// 	router.Use(cors.New(config))
// 	publicRoutes := router.Group("/api/v1")
// 	{
// 		publicRoutes.POST("/register", userController.RegisterUser)
// 		publicRoutes.POST("/login", userController.LoginUser)
// 	}
// 	authenticatedRoutes := router.Group("/api/v1")
// 	authenticatedRoutes.Use(middleware.AuthMiddleware(jwtService))
// 	{
		
// 		authenticatedRoutes.GET("/profile", userController.GetUserProfile)
// 		authenticatedRoutes.PATCH("/profile", userController.UpdateUserProfile)
// 		authenticatedRoutes.GET("/consistency", consistencyController.GetDailyConsistency) // Can take 'date' query param
// 		authenticatedRoutes.GET("/consistency/history", consistencyController.GetConsistencyHistory) // Takes 'startDate', 'endDate' query params
// 		authenticatedRoutes.GET("/consistency/streaks", consistencyController.GetUserStreaks)
// 		authenticatedRoutes.POST("/consistency/check", consistencyController.TriggerDailyConsistencyCheck) // Manual trigger for debugging
// 	}

// 	return router
// }



package routers

import (
	"consistent_1/Delivery/controllers"
	"consistent_1/Delivery/middleware"
	"consistent_1/Infrastructure/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController *controllers.UserController,
	consistencyController *controllers.ConsistencyController,
	jwtService auth.JWTService,
) *gin.Engine {
	// --- START MODIFICATION ---
	// Create a new Gin engine *without* default middleware (Logger and Recovery)
	router := gin.New()

	// Optionally add back Logger if you still want request logging, but exclude Recovery for now
	router.Use(gin.Logger())
	// --- END MODIFICATION ---

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"} // React/Flutter web defaults or your specific frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	publicRoutes := router.Group("/api/v1")
	{
		publicRoutes.POST("/register", userController.RegisterUser)
		publicRoutes.POST("/login", userController.LoginUser)
	}

	authenticatedRoutes := router.Group("/api/v1")
	authenticatedRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		authenticatedRoutes.GET("/profile", userController.GetUserProfile)
		authenticatedRoutes.PATCH("/profile", userController.UpdateUserProfile)
		authenticatedRoutes.GET("/consistency", consistencyController.GetDailyConsistency) // Can take 'date' query param
		authenticatedRoutes.GET("/consistency/history", consistencyController.GetConsistencyHistory) // Takes 'startDate', 'endDate' query params
		authenticatedRoutes.GET("/consistency/streaks", consistencyController.GetUserStreaks)
		authenticatedRoutes.POST("/consistency/check", consistencyController.TriggerDailyConsistencyCheck) // Manual trigger for debugging
	}

	return router
}