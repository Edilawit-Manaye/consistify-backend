package platform_api

import (
	"context"
	"consistent_1/Domain" 
	"time"
)
type LeetCodeAPI interface {
	FetchUserDailyActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error)
}
type CodeforcesAPI interface {
	FetchUserDailyActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error)
}

