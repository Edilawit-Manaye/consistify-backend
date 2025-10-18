package usecases

import (
	"context"
	"fmt"
	"time"

	"consistent_1/Domain"
	"consistent_1/Infrastructure/auth"
	"consistent_1/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type UserUsecase interface {
	RegisterUser(ctx context.Context, req *domain.UserRegisterRequest) (*domain.User, error)
	LoginUser(ctx context.Context, req *domain.UserLoginRequest) (string, error) 
	UpdateUserProfile(ctx context.Context, userID string, updates *domain.UserProfileUpdateRequest) error
	GetUserProfile(ctx context.Context, userID string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error) 
}

type userUsecase struct {
	userRepo      repositories.UserRepository
	passwordService auth.PasswordService
	jwtService    auth.JWTService
}
func NewUserUsecase(
	userRepo repositories.UserRepository,
	passwordService auth.PasswordService,
	jwtService auth.JWTService,
) UserUsecase {
	return &userUsecase{
		userRepo:      userRepo,
		passwordService: passwordService,
		jwtService:    jwtService,
	}
}
func (uc *userUsecase) RegisterUser(ctx context.Context, req *domain.UserRegisterRequest) (*domain.User, error) {
	if req.Password != req.ConfirmPassword {
		return nil, domain.ErrPasswordsDoNotMatch
	}
	_, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyExists 
	}
	if err != nil && err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("database error checking email: %w", err)
	}
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = time.Parse("15:04", req.NotificationTime)
	if err != nil {
		return nil, domain.ErrInvalidNotificationTime
	}
	_, err = time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone provided: %w", err)
	}
	user := &domain.User{
		Email:              req.Email,
		PasswordHash:       hashedPassword,
		Username:           req.Username,
		NotificationTime:   req.NotificationTime,
		Timezone:           req.Timezone,
		PlatformUsernames:  make(map[string]string), 
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		FCMTokens:          []string{},              
	}
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	return user, nil
}
func (uc *userUsecase) LoginUser(ctx context.Context, req *domain.UserLoginRequest) (string, error) {

	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("database error retrieving user: %w", err)
	}
	if err := uc.passwordService.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		return "", domain.ErrInvalidCredentials
	}
	token, err := uc.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}
func (uc *userUsecase) UpdateUserProfile(ctx context.Context, userID string, updates *domain.UserProfileUpdateRequest) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrUserNotFound 
	}

	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err 
	}

	if updates.Username != nil {
		user.Username = *updates.Username
	}
	if updates.NotificationTime != nil {

		_, err := time.Parse("15:04", *updates.NotificationTime)
		if err != nil {
			return domain.ErrInvalidNotificationTime
		}
		user.NotificationTime = *updates.NotificationTime
	}
	if updates.Timezone != nil {
		_, err := time.LoadLocation(*updates.Timezone)
		if err != nil {
			return fmt.Errorf("invalid timezone provided: %w", err)
		}
		user.Timezone = *updates.Timezone
	}
	if updates.PlatformUsernames != nil {
		user.PlatformUsernames = updates.PlatformUsernames
	}
	if updates.FCMToken != nil && *updates.FCMToken != "" {
		found := false
		for _, t := range user.FCMTokens {
			if t == *updates.FCMToken {
				found = true
				break
			}
		}
		if !found {
			user.FCMTokens = append(user.FCMTokens, *updates.FCMToken)
		}
	}
	user.ID = objID 

	return uc.userRepo.UpdateUser(ctx, user)
}
func (uc *userUsecase) GetUserProfile(ctx context.Context, userID string) (*domain.User, error) {
	return uc.userRepo.GetUserByID(ctx, userID)
}
func (uc *userUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.GetAllUsers(ctx)
}