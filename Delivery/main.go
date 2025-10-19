// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"consistent_1/Delivery/controllers"
// 	"consistent_1/Delivery/routers"
// 	"consistent_1/Infrastructure/auth"
// 	"consistent_1/Infrastructure/database"
// 	"consistent_1/Infrastructure/notifications"
// 	"consistent_1/Infrastructure/platform_api"
// 	"consistent_1/Infrastructure/scheduler"
// 	"consistent_1/Repositories"
// 	"consistent_1/Usecases"

// 	firebase "firebase.google.com/go" 
// 	"github.com/spf13/viper"
// 	"google.golang.org/api/option"    
// )

// func main() {
// 	viper.SetConfigFile(".env") 
// 	viper.AutomaticEnv()       
// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	serverPort := viper.GetString("SERVER_PORT")
// 	if serverPort == "" {
// 		serverPort = ":8080" 
// 	}
// 	jwtSecret := viper.GetString("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET not set in environment variables")
// 	}

	
// 	firebaseServiceAccountPath := viper.GetString("FIREBASE_SERVICE_ACCOUNT_PATH")
// 	if firebaseServiceAccountPath == "" {
// 		log.Fatal("FIREBASE_SERVICE_ACCOUNT_PATH not set in environment variables. Push notifications will not work.")
// 	}

	
// 	firebaseApp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(firebaseServiceAccountPath))
// 	if err != nil {
// 		log.Fatalf("Error initializing Firebase app: %v", err)
// 	}
// 	log.Println("Firebase Admin SDK initialized successfully.")
// 	mongoClient, err := database.NewMongoClient()
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}
// 	defer func() {
// 		if err := mongoClient.Disconnect(context.Background()); err != nil {
// 			log.Printf("Error disconnecting from MongoDB: %v", err)
// 		}
// 	}()

// 	passwordService := auth.NewPasswordService()
// 	jwtService := auth.NewJWTService(jwtSecret)
// 	fcmService := notifications.NewFCMService(firebaseApp) 
// 	leetcodeAPI := platform_api.NewLeetCodeAPI(viper.GetString("LEETCODE_API_BASE_URL"))
// 	codeforcesAPI := platform_api.NewCodeforcesAPI(viper.GetString("CODEFORCES_API_BASE_URL"))
// 	userRepo := repositories.NewUserRepository(mongoClient.DB)
// 	consistencyRepo := repositories.NewConsistencyRepository(mongoClient.DB)
// 	platformUsecase := usecases.NewPlatformUsecase(userRepo, leetcodeAPI, codeforcesAPI)
// 	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService)
// 	consistencyUsecase := usecases.NewConsistencyUsecase(userRepo, consistencyRepo, platformUsecase, fcmService)
// 	userController := controllers.NewUserController(userUsecase)
// 	consistencyController := controllers.NewConsistencyController(consistencyUsecase)
// 	router := routers.SetupRouter(userController, consistencyController, jwtService)
// 	consistencyScheduler := scheduler.NewConsistencyScheduler(consistencyUsecase, userUsecase)
// 	consistencyScheduler.ScheduleDailyConsistencyCheck()
// 	consistencyScheduler.ScheduleNotificationReminders()
// 	consistencyScheduler.Start()
// 	defer consistencyScheduler.Stop() 
// 	go func() {
// 		log.Printf("Server starting on %s", serverPort)
// 		if err := router.Run(serverPort); err != nil {
// 			log.Fatalf("Server failed to start: %v", err)
// 		}
// 	}()
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	log.Println("Shutting down server...")
// 	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	log.Println("Server gracefully stopped.")
// }




package main

import (
	"context"
	"log"
	"os" // Ensure "os" is imported
	"os/signal"
	"syscall"
	"time"

	"consistent_1/Delivery/controllers"
	"consistent_1/Delivery/routers"
	"consistent_1/Infrastructure/auth"
	"consistent_1/Infrastructure/database"
	"consistent_1/Infrastructure/notifications"
	"consistent_1/Infrastructure/platform_api"
	"consistent_1/Infrastructure/scheduler"
	"consistent_1/Repositories"
	"consistent_1/Usecases"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serverPort := viper.GetString("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set in environment variables")
	}

	// --- START MODIFIED FIREBASE INITIALIZATION ---

	// Read the service account file path from .env (which Dockerfile creates and points to the dynamic file)
	firebaseServiceAccountPath := viper.GetString("FIREBASE_SERVICE_ACCOUNT_PATH")
	if firebaseServiceAccountPath == "" {
		log.Fatal("FIREBASE_SERVICE_ACCOUNT_PATH not set in environment variables. Push notifications will not work.")
	}

	// Read the project ID from a Render environment variable (set on Render's UI)
	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID") // Using os.Getenv to directly read Render env var
	if firebaseProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID environment variable not set on Render. Push notifications will not work.")
	}

	// Create a Firebase config with the explicit project ID. This is the fix.
	config := &firebase.Config{
		ProjectID: firebaseProjectID,
	}

	// Initialize Firebase app using the explicit config and the credentials file.
	// The 'nil' from your original code is replaced by 'config'.
	firebaseApp, err := firebase.NewApp(context.Background(), config, option.WithCredentialsFile(firebaseServiceAccountPath))
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}
	log.Println("Firebase Admin SDK initialized successfully.")

	// --- END MODIFIED FIREBASE INITIALIZATION ---


	mongoClient, err := database.NewMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJWTService(jwtSecret)
	fcmService := notifications.NewFCMService(firebaseApp)
	leetcodeAPI := platform_api.NewLeetCodeAPI(viper.GetString("LEETCODE_API_BASE_URL"))
	codeforcesAPI := platform_api.NewCodeforcesAPI(viper.GetString("CODEFORCES_API_BASE_URL"))
	userRepo := repositories.NewUserRepository(mongoClient.DB)
	consistencyRepo := repositories.NewConsistencyRepository(mongoClient.DB)
	platformUsecase := usecases.NewPlatformUsecase(userRepo, leetcodeAPI, codeforcesAPI)
	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService)
	consistencyUsecase := usecases.NewConsistencyUsecase(userRepo, consistencyRepo, platformUsecase, fcmService)
	userController := controllers.NewUserController(userUsecase)
	consistencyController := controllers.NewConsistencyController(consistencyUsecase)
	router := routers.SetupRouter(userController, consistencyController, jwtService)
	consistencyScheduler := scheduler.NewConsistencyScheduler(consistencyUsecase, userUsecase)
	consistencyScheduler.ScheduleDailyConsistencyCheck()
	consistencyScheduler.ScheduleNotificationReminders()
	consistencyScheduler.Start()
	defer consistencyScheduler.Stop()
	go func() {
		log.Printf("Server starting on %s", serverPort)
		if err := router.Run(serverPort); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Server gracefully stopped.")
}






