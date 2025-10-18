package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type DailyConsistency struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"userId" json:"userId"`
	Date               time.Time          `bson:"date" json:"date"`                        
	PlatformActivities []PlatformActivity `bson:"platformActivities" json:"platformActivities"` 
	OverallConsistent  bool               `bson:"overallConsistent" json:"overallConsistent"` // True if user met overall daily goal (e.g., solved at least one problem on any platform)
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}
func (dc *DailyConsistency) GetPlatformActivity(platformName string) PlatformActivity {
	for _, activity := range dc.PlatformActivities {
		if activity.Platform == platformName { 
			return activity
		}
	}
	
	return PlatformActivity{}
}



type StreakInfo struct {
	CurrentStreak     int        `json:"currentStreak"`
	LongestStreak     int        `json:"longestStreak"`
	LastConsistentDay *time.Time `json:"lastConsistentDay,omitempty"`
}


type ConsistencyFilter struct {
	UserID    primitive.ObjectID
	StartDate *time.Time
	EndDate   *time.Time
}