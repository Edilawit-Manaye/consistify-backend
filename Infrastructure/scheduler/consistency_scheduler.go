


package scheduler

import (
	"context"
	"log"
	"time"

	
	"consistent_1/Usecases"

	"github.com/robfig/cron/v3"
)
type ConsistencyScheduler struct {
	Cron *cron.Cron
	ConsistencyUsecase usecases.ConsistencyUsecase
	UserUsecase        usecases.UserUsecase
}
func NewConsistencyScheduler(
	consistencyUsecase usecases.ConsistencyUsecase,
	userUsecase usecases.UserUsecase,
) *ConsistencyScheduler {
	c := cron.New() 
	return &ConsistencyScheduler{
		Cron: c,
		ConsistencyUsecase: consistencyUsecase,
		UserUsecase:        userUsecase,
	}
}
func (s *ConsistencyScheduler) Start() {
	s.Cron.Start()
	log.Println("Consistency scheduler started.")
}
func (s *ConsistencyScheduler) Stop() {
	s.Cron.Stop()
	log.Println("Consistency scheduler stopped.")
}
func (s *ConsistencyScheduler) ScheduleDailyConsistencyCheck() {
	_, err := s.Cron.AddFunc("5 0 * * *", func() { 
		log.Println("Running daily consistency check for all users (server time)...")
		users, err := s.UserUsecase.GetAllUsers(context.Background())
		if err != nil {
			log.Printf("Error fetching all users for consistency check: %v", err)
			return
		}

		for _, user := range users {
			log.Printf("Checking consistency for user: %s (ID: %s)", user.Email, user.ID.Hex())
			_, err := s.ConsistencyUsecase.CheckDailyConsistency(context.Background(), user.ID.Hex())
			if err != nil {
				log.Printf("Error checking consistency for user %s: %v", user.ID.Hex(), err)
			}
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling daily consistency check: %v", err)
	}
	log.Println("Daily consistency check scheduled for 00:05 AM (server time).")
}
func (s *ConsistencyScheduler) ScheduleNotificationReminders() {
	_, err := s.Cron.AddFunc("0 * * * *", func() { 
		log.Println("Running hourly notification reminder check...")
		users, err := s.UserUsecase.GetAllUsers(context.Background())
		if err != nil {
			log.Printf("Error fetching users for notification check: %v", err)
			return
		}

		nowInUTC := time.Now().UTC() 

		for _, user := range users {
			if user.NotificationTime == "" || user.Timezone == "" {
				log.Printf("User %s (ID: %s) has incomplete notification settings, skipping reminder.", user.Email, user.ID.Hex())
				continue
			}

			loc, err := time.LoadLocation(user.Timezone)
			if err != nil {
				log.Printf("Invalid timezone '%s' for user %s (ID: %s): %v. Skipping reminder.", user.Timezone, user.Email, user.ID.Hex(), err)
				continue
			}
			nowInUserLocalTime := nowInUTC.In(loc)
			userPreferredTime := nowInUserLocalTime.Format("15:04") 
			if user.NotificationTime == userPreferredTime {
				log.Printf("It's %s in user's (%s) timezone (%s). Attempting to send notification.",
					userPreferredTime, user.Email, user.Timezone)
				err := s.ConsistencyUsecase.SendConsistencyReminder(context.Background(), user.ID.Hex())
				if err != nil {
					log.Printf("Error sending reminder to user %s: %v", user.ID.Hex(), err)
				}
			}
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling hourly notification reminder: %v", err)
	}
	log.Println("Hourly notification reminder check scheduled.")
}