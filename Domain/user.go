package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email                     string             `bson:"email" json:"email"`
	PasswordHash              string             `bson:"passwordHash" json:"-"` 
	Username                  string             `bson:"username" json:"username"`
	PlatformUsernames         map[string]string  `bson:"platformUsernames" json:"platformUsernames"` 
	NotificationTime          string             `bson:"notificationTime" json:"notificationTime"`  
	Timezone                  string             `bson:"timezone" json:"timezone"`                   
	FCMTokens                 []string           `bson:"fcmTokens,omitempty" json:"fcmTokens,omitempty"` 
	CreatedAt                 time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt                 time.Time          `bson:"updatedAt" json:"updatedAt"`
	LeetCodeLastTotalSolved int       `bson:"leetcodeLastTotalSolved,omitempty" json:"leetcodeLastTotalSolved,omitempty"`
	LeetCodeLastCheckDate   time.Time `bson:"leetcodeLastCheckDate,omitempty" json:"leetcodeLastCheckDate,omitempty"`
	
}
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type UserRegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	ConfirmPassword  string `json:"confirmPassword" binding:"required,min=6"`
	Username         string `json:"username" binding:"required,min=3"`
	NotificationTime string `json:"notificationTime" binding:"required"` 
	Timezone         string `json:"timezone" binding:"required"`         
}
type UserProfileUpdateRequest struct {
	Username          *string            `json:"username,omitempty"`
	NotificationTime  *string            `json:"notificationTime,omitempty"`
	Timezone          *string            `json:"timezone,omitempty"`
	PlatformUsernames map[string]string `json:"platformUsernames,omitempty"`
	FCMToken          *string            `json:"fcmToken,omitempty"` 
}
type FCMNotification struct {
	To           string            `json:"to"`                 
	Priority     string            `json:"priority,omitempty"` 
	Notification *struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"notification,omitempty"`
	Data map[string]string `json:"data,omitempty"` 
}
type FCMResponse struct {
	MulticastID  int64 `json:"multicast_id"`
	Success      int   `json:"success"`
	Failure      int   `json:"failure"`
	CanonicalIDs int   `json:"canonical_ids"`
	Results      []struct {
		MessageID    string `json:"message_id,omitempty"`
		Error        string `json:"error,omitempty"`
		RegistrationID string `json:"registration_id,omitempty"` 
	} `json:"results"`
}



