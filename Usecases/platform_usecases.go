


package usecases

import (
	"context"
	"fmt"
	"time"

	"consistent_1/Domain"
	"consistent_1/Infrastructure/platform_api" 
	"consistent_1/Repositories"

)


type PlatformUsecase interface {
	
	FetchUserDailyActivity(ctx context.Context, userID string, date time.Time) ([]domain.PlatformActivity, error)

	
	FetchLeetCodeActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error)
	FetchCodeforcesActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error)
}

type platformUsecase struct {
	userRepo      repositories.UserRepository
	leetcodeAPI   platform_api.LeetCodeAPI    
	codeforcesAPI platform_api.CodeforcesAPI  
	
}


func NewPlatformUsecase(
	userRepo repositories.UserRepository,
	leetcodeAPI platform_api.LeetCodeAPI,
	codeforcesAPI platform_api.CodeforcesAPI,
	
) PlatformUsecase {
	return &platformUsecase{
		userRepo:      userRepo,
		leetcodeAPI:   leetcodeAPI,
		codeforcesAPI: codeforcesAPI,
		
	}
}


func (uc *platformUsecase) FetchUserDailyActivity(ctx context.Context, userID string, date time.Time) ([]domain.PlatformActivity, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allActivities []domain.PlatformActivity
	
	queryDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)


	
	if lcUsername, ok := user.PlatformUsernames["leetcode"]; ok && lcUsername != "" {
		
		activity, err := uc.FetchLeetCodeActivity(ctx, lcUsername, queryDate)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch LeetCode activity: %w", err)
		}
		allActivities = append(allActivities, activity)
	}

	
	if cfUsername, ok := user.PlatformUsernames["codeforces"]; ok && cfUsername != "" {
		
		activity, err := uc.FetchCodeforcesActivity(ctx, cfUsername, queryDate)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Codeforces activity: %w", err)
		}
		allActivities = append(allActivities, activity)
	}



	return allActivities, nil
}


func (uc *platformUsecase) FetchLeetCodeActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error) {
	
	return uc.leetcodeAPI.FetchUserDailyActivity(ctx, username, date)
}


func (uc *platformUsecase) FetchCodeforcesActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error) {
	return uc.codeforcesAPI.FetchUserDailyActivity(ctx, username, date)
}

