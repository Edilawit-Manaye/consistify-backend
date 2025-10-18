



package usecases

import (
	"context"
	"fmt"
	"log" 
	"time"

	"consistent_1/Domain"
	"consistent_1/Infrastructure/notifications"
	"consistent_1/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ConsistencyUsecase interface {
	CheckDailyConsistency(ctx context.Context, userID string) (*domain.DailyConsistency, error)
	GetDailyConsistency(ctx context.Context, userID string, date time.Time) (*domain.DailyConsistency, error)
	GetConsistencyHistory(ctx context.Context, userID string, startDate, endDate *time.Time) ([]domain.DailyConsistency, error)
	GetStreaks(ctx context.Context, userID string) (*domain.StreakInfo, error)
	SendConsistencyReminder(ctx context.Context, userID string) error
	TriggerDailyConsistencyCheck(ctx context.Context)
}

type consistencyUsecase struct {
	userRepo        repositories.UserRepository
	consistencyRepo repositories.ConsistencyRepository
	platformUsecase PlatformUsecase
	fcmService      notifications.FCMService
}
func NewConsistencyUsecase(
	userRepo repositories.UserRepository,
	consistencyRepo repositories.ConsistencyRepository,
	platformUsecase PlatformUsecase,
	fcmService notifications.FCMService,
) ConsistencyUsecase {
	return &consistencyUsecase{
		userRepo:        userRepo,
		consistencyRepo: consistencyRepo,
		platformUsecase: platformUsecase,
		fcmService:      fcmService,
	}
}


func (uc *consistencyUsecase) CheckDailyConsistency(ctx context.Context, userID string) (*domain.DailyConsistency, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user %s for daily consistency check: %w", userID, err)
	}

	todayUTC := time.Now().UTC().Truncate(24 * time.Hour) 

	var platformActivities []domain.PlatformActivity
	overallConsistent := false
	if leetcodeUsername, ok := user.PlatformUsernames["leetcode"]; ok && leetcodeUsername != "" {
		leetcodeCurrentActivity, err := uc.platformUsecase.FetchLeetCodeActivity(ctx, leetcodeUsername, todayUTC)
		if err != nil {
			log.Printf("Error fetching LeetCode activity for user %s (%s): %v", userID, leetcodeUsername, err) // Keep error log
			platformActivities = append(platformActivities, domain.PlatformActivity{
				Platform:       "leetcode",
				Username:       leetcodeUsername,
				Date:           todayUTC,
				IsConsistent:   false,
				ProblemsSolved: 0,
			})
		} else {
			problemsSolvedTodayLeetCode := 0
			isLeetCodeConsistent := false

			
			if user.LeetCodeLastCheckDate.IsZero() && user.LeetCodeLastTotalSolved == 0 {
				if leetcodeCurrentActivity.ProblemsSolved > 0 {
					problemsSolvedTodayLeetCode = leetcodeCurrentActivity.ProblemsSolved 
					isLeetCodeConsistent = true
				} else {
					problemsSolvedTodayLeetCode = 0
					isLeetCodeConsistent = false
				}
			
			} else if user.LeetCodeLastCheckDate.Before(todayUTC) {
				problemsSolvedTodayLeetCode = leetcodeCurrentActivity.ProblemsSolved - user.LeetCodeLastTotalSolved
				if problemsSolvedTodayLeetCode > 0 {
					isLeetCodeConsistent = true
				}

			} else if user.LeetCodeLastCheckDate.Equal(todayUTC) {
				if leetcodeCurrentActivity.ProblemsSolved > user.LeetCodeLastTotalSolved {
					problemsSolvedTodayLeetCode = leetcodeCurrentActivity.ProblemsSolved - user.LeetCodeLastTotalSolved
					isLeetCodeConsistent = problemsSolvedTodayLeetCode > 0
				} else {
					existingDailyCons, getErr := uc.consistencyRepo.GetDailyConsistency(ctx, objUserID, todayUTC)
					if getErr == nil && existingDailyCons != nil {
						activityForPlatform := existingDailyCons.GetPlatformActivity("leetcode")
						isLeetCodeConsistent = activityForPlatform.IsConsistent
						problemsSolvedTodayLeetCode = activityForPlatform.ProblemsSolved
					} else {
						isLeetCodeConsistent = false
						problemsSolvedTodayLeetCode = 0
						if getErr != nil && getErr != domain.ErrConsistencyNotFound {
							log.Printf("Warning: Error retrieving existing LeetCode consistency for user %s on %s: %v", userID, todayUTC.Format(time.RFC3339), getErr) // Keep warning log
						}
					}
				}
			}
			err = uc.userRepo.UpdateUserLeetCodeStats(ctx, objUserID, leetcodeCurrentActivity.ProblemsSolved, todayUTC)
			if err != nil {
				log.Printf("Warning: Failed to update LeetCode stats for user %s (%s): %v", userID, leetcodeUsername, err) // Keep warning log
			}

			platformActivities = append(platformActivities, domain.PlatformActivity{
				Platform:       "leetcode",
				Username:       leetcodeUsername,
				Date:           todayUTC,
				IsConsistent:   isLeetCodeConsistent,
				ProblemsSolved: problemsSolvedTodayLeetCode,
			})
			if isLeetCodeConsistent {
				overallConsistent = true
			}
		}
	}
	if codeforcesUsername, ok := user.PlatformUsernames["codeforces"]; ok && codeforcesUsername != "" {
		codeforcesActivity, err := uc.platformUsecase.FetchCodeforcesActivity(ctx, codeforcesUsername, todayUTC)
		if err != nil {
			log.Printf("Error fetching Codeforces activity for user %s (%s): %v", userID, codeforcesUsername, err) // Keep error log
			platformActivities = append(platformActivities, domain.PlatformActivity{
				Platform:       "codeforces",
				Username:       codeforcesUsername,
				Date:           todayUTC,
				IsConsistent:   false,
				ProblemsSolved: 0,
			})
		} else {
			platformActivities = append(platformActivities, codeforcesActivity)
			if codeforcesActivity.IsConsistent {
				overallConsistent = true
			}
		}
	}
	dailyConsistency, err := uc.consistencyRepo.GetDailyConsistency(ctx, objUserID, todayUTC)
	if err != nil && err != domain.ErrConsistencyNotFound {
		return nil, fmt.Errorf("database error getting daily consistency: %w", err)
	}

	if dailyConsistency == nil {
		dailyConsistency = &domain.DailyConsistency{
			UserID:             objUserID,
			Date:               todayUTC,
			PlatformActivities: platformActivities,
			OverallConsistent:  overallConsistent,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
	} else {
		dailyConsistency.PlatformActivities = platformActivities
		dailyConsistency.OverallConsistent = overallConsistent
		dailyConsistency.UpdatedAt = time.Now()
	}

	
	if err := uc.consistencyRepo.SaveDailyConsistency(ctx, dailyConsistency); err != nil {
		return nil, fmt.Errorf("failed to save daily consistency for user %s: %w", userID, err)
	}

	return dailyConsistency, nil
}
func (uc *consistencyUsecase) GetDailyConsistency(ctx context.Context, userID string, date time.Time) (*domain.DailyConsistency, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return uc.consistencyRepo.GetDailyConsistency(ctx, objUserID, date.Truncate(24*time.Hour)) // Ensure date is truncated for consistent lookup
}
func (uc *consistencyUsecase) GetConsistencyHistory(ctx context.Context, userID string, startDate, endDate *time.Time) ([]domain.DailyConsistency, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	filter := domain.ConsistencyFilter{
		UserID:    objUserID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	return uc.consistencyRepo.GetConsistencyHistory(ctx, filter)
}
func (uc *consistencyUsecase) GetStreaks(ctx context.Context, userID string) (*domain.StreakInfo, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return uc.consistencyRepo.GetStreaks(ctx, objUserID)
}
func (uc *consistencyUsecase) SendConsistencyReminder(ctx context.Context, userID string) error {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user for reminder: %w", err)
	}

	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)
	dailyConsistency, err := uc.consistencyRepo.GetDailyConsistency(ctx, user.ID, todayUTC)
	if err != nil && err != domain.ErrConsistencyNotFound {
		return fmt.Errorf("error checking daily consistency for reminder: %w", err)
	}
	if dailyConsistency == nil || !dailyConsistency.OverallConsistent {
		if len(user.FCMTokens) == 0 {
			return nil
		}

		title := "Consistify Reminder! ⏰"
		body := fmt.Sprintf("Hey %s, you haven't solved today's challenge yet! Let's keep your streak alive 💪.", user.Username)
		data := map[string]string{"type": "consistency_reminder", "userId": userID}
		for _, token := range user.FCMTokens {
			err := uc.fcmService.SendNotification(ctx, token, title, body, data)
			if err != nil {
				log.Printf("Failed to send FCM notification to token %s for user %s: %v", token, userID, err) // Keep error log
				
			}
		}
		return nil
	}
	return nil
}


func (uc *consistencyUsecase) TriggerDailyConsistencyCheck(ctx context.Context) {
	log.Println("Scheduler: Triggering daily consistency check for all users...") 
	users, err := uc.userRepo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Scheduler: Failed to get all users for daily check: %v", err) 
		return
	}

	for _, user := range users {
		_, err := uc.CheckDailyConsistency(ctx, user.ID.Hex())
		if err != nil {
			log.Printf("Scheduler: Error checking consistency for user %s (ID: %s): %v", user.Email, user.ID.Hex(), err) // Keep error log
		}
	}
	log.Println("Scheduler: Daily consistency check completed for all users.") // Keep this high-level log
}