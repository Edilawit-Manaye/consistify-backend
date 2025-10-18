package repositories

import (
	"context"
	"time"

	"consistent_1/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type ConsistencyRepository interface {
	SaveDailyConsistency(ctx context.Context, consistency *domain.DailyConsistency) error
	GetDailyConsistency(ctx context.Context, userID primitive.ObjectID, date time.Time) (*domain.DailyConsistency, error)
	GetConsistencyHistory(ctx context.Context, filter domain.ConsistencyFilter) ([]domain.DailyConsistency, error)
	GetStreaks(ctx context.Context, userID primitive.ObjectID) (*domain.StreakInfo, error)
}

type consistencyRepository struct {
	collection *mongo.Collection
}


func NewConsistencyRepository(db *mongo.Database) ConsistencyRepository {
	return &consistencyRepository{
		collection: db.Collection("daily_consistencies"),
	}
}


func (r *consistencyRepository) SaveDailyConsistency(ctx context.Context, consistency *domain.DailyConsistency) error {
	
	consistency.Date = time.Date(consistency.Date.Year(), consistency.Date.Month(), consistency.Date.Day(), 0, 0, 0, 0, time.UTC)

	filter := bson.M{"userId": consistency.UserID, "date": consistency.Date}
	update := bson.M{"$set": bson.M{
		"platformActivities": consistency.PlatformActivities,
		"overallConsistent":  consistency.OverallConsistent,
		"updatedAt":          time.Now(),
	}}
	opts := options.Update().SetUpsert(true) 

	if consistency.ID.IsZero() { 
		consistency.CreatedAt = time.Now()
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.UpsertedID != nil {
		
		consistency.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return nil
}
func (r *consistencyRepository) GetDailyConsistency(ctx context.Context, userID primitive.ObjectID, date time.Time) (*domain.DailyConsistency, error) {
	var consistency domain.DailyConsistency
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	err := r.collection.FindOne(ctx, bson.M{"userId": userID, "date": normalizedDate}).Decode(&consistency)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrConsistencyNotFound
	}
	return &consistency, err
}
func (r *consistencyRepository) GetConsistencyHistory(ctx context.Context, filter domain.ConsistencyFilter) ([]domain.DailyConsistency, error) {
	bsonFilter := bson.M{"userId": filter.UserID} 

	if filter.StartDate != nil && filter.EndDate != nil {
		startOfDay := time.Date(filter.StartDate.Year(), filter.StartDate.Month(), filter.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := time.Date(filter.EndDate.Year(), filter.EndDate.Month(), filter.EndDate.Day(), 23, 59, 59, 999999999, time.UTC) // End of day
		bsonFilter["date"] = bson.M{"$gte": startOfDay, "$lte": endOfDay}
	} else if filter.StartDate != nil {
		startOfDay := time.Date(filter.StartDate.Year(), filter.StartDate.Month(), filter.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		bsonFilter["date"] = bson.M{"$gte": startOfDay}
	} else if filter.EndDate != nil {
		endOfDay := time.Date(filter.EndDate.Year(), filter.EndDate.Month(), filter.EndDate.Day(), 23, 59, 59, 999999999, time.UTC)
		bsonFilter["date"] = bson.M{"$lte": endOfDay}
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}}) 

	var consistencies []domain.DailyConsistency
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &consistencies); err != nil {
		return nil, err
	}
	return consistencies, nil
}


func (r *consistencyRepository) GetStreaks(ctx context.Context, userID primitive.ObjectID) (*domain.StreakInfo, error) {

	consistencies, err := r.GetConsistencyHistory(ctx, domain.ConsistencyFilter{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	streakInfo := &domain.StreakInfo{}
	if len(consistencies) == 0 {
		return streakInfo, nil
	}


	var consistentDays []time.Time
	for _, dc := range consistencies {
		if dc.OverallConsistent {
			consistentDays = append(consistentDays, dc.Date)
		}
	}

	if len(consistentDays) == 0 {
		return streakInfo, nil 
	}



	var currentStreak int
	var longestStreak int
	var lastConsistentDay *time.Time

	today := time.Now().UTC().Truncate(24 * time.Hour) 
	yesterday := today.AddDate(0, 0, -1) 
	for i := len(consistentDays) - 1; i >= 0; i-- {
		day := consistentDays[i]
		if lastConsistentDay == nil {
			if day.Equal(today) || day.Equal(yesterday) {
				currentStreak = 1
				lastConsistentDay = &day
			} else {
			
				currentStreak = 0
				break
			}
		} else {
			
			if day.AddDate(0, 0, 1).Equal(*lastConsistentDay) {
				currentStreak++
				lastConsistentDay = &day 
			} else {
			
				break
			}
		}
	}


	
	if len(consistentDays) > 0 {
		currentCount := 0
		for i := 0; i < len(consistentDays); i++ {
			if i == 0 || consistentDays[i].Equal(consistentDays[i-1].AddDate(0, 0, 1)) {
				currentCount++
			} else {
				currentCount = 1 
			}
			if currentCount > longestStreak {
				longestStreak = currentCount
			}
		}
	}
	if len(consistentDays) > 0 {
		mostRecentConsistentDay := consistentDays[len(consistentDays)-1]
		streakInfo.LastConsistentDay = &mostRecentConsistentDay
	}

	actualCurrentStreak := 0
	if len(consistentDays) > 0 {
		mostRecentConsistent := consistentDays[len(consistentDays)-1]
		if mostRecentConsistent.Equal(today) {
			actualCurrentStreak = 1
			for i := len(consistentDays) - 2; i >= 0; i-- {
				if consistentDays[i].Equal(consistentDays[i+1].AddDate(0, 0, -1)) {
					actualCurrentStreak++
				} else {
					break
				}
			}
		} else if mostRecentConsistent.Equal(yesterday) {
			actualCurrentStreak = 1
			for i := len(consistentDays) - 2; i >= 0; i-- {
				if consistentDays[i].Equal(consistentDays[i+1].AddDate(0, 0, -1)) {
					actualCurrentStreak++
				} else {
					break
				}
			}
		} else {
			actualCurrentStreak = 0 
		}
	}


	streakInfo.CurrentStreak = actualCurrentStreak
	streakInfo.LongestStreak = longestStreak

	return streakInfo, nil
}