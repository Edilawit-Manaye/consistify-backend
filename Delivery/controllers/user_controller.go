package controllers

import (
	"log"
	"net/http"

	domain "consistent_1/Domain"
	usecases "consistent_1/Usecases"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecases.UserUsecase
}

func NewUserController(userUsecase usecases.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

// func (ctrl *UserController) RegisterUser(c *gin.Context) {
// 	var req domain.UserRegisterRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user, err := ctrl.userUsecase.RegisterUser(c.Request.Context(), &req)
// 	if err != nil {
// 		switch err {
// 		case domain.ErrPasswordsDoNotMatch:
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		case domain.ErrEmailAlreadyExists:
// 			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
// 		case domain.ErrInvalidNotificationTime:
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		default:
// 			log.Printf("Error registering user: %v", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message":  "User registered successfully",
// 		"userID":   user.ID.Hex(),
// 		"username": user.Username,
// 		"email":    user.Email,
// 	})
// }

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	log.Printf("--- RegisterUser: Received request ---") // Log entry point

	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("RegisterUser: JSON binding error: %v", err) // Log binding error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Log success and request data. Be careful not to log sensitive info like raw passwords in production.
	log.Printf("RegisterUser: JSON bound successfully. Email: %s, Username: %s", req.Email, req.Username)

	user, err := ctrl.userUsecase.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		log.Printf("RegisterUser: Error from usecase: %v", err) // Log usecase error
		switch err {
		case domain.ErrPasswordsDoNotMatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case domain.ErrInvalidNotificationTime:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("RegisterUser: Unhandled error during registration: %v", err) // More specific message for default
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	log.Printf("RegisterUser: User registered successfully: %s", user.Email) // Log success
	c.JSON(http.StatusCreated, gin.H{
		"message":  "User registered successfully",
		"userID":   user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
	})
}

func (ctrl *UserController) LoginUser(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.userUsecase.LoginUser(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			log.Printf("Error logging in user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (ctrl *UserController) GetUserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	user, err := ctrl.userUsecase.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			log.Printf("Error getting user profile for %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		}
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var req domain.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userUsecase.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrInvalidNotificationTime:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("Error updating user profile for %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
