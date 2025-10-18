// package domain

// // PlatformActivity represents a user's activity on a single coding platform for a day.
// type PlatformActivity struct {
// 	PlatformName   string `bson:"platformName" json:"platformName"`    // e.g., "leetcode", "codeforces"
// 	ProblemsSolved int    `bson:"problemsSolved" json:"problemsSolved"` // Number of problems solved
// 	IsConsistent   bool   `bson:"isConsistent" json:"isConsistent"`    // True if daily goal met for this platform (e.g., solved >= 1 problem)
// }



package domain

import (
	"time"
)


type PlatformActivity struct {
	Platform       string    `bson:"platform" json:"platform"`            
	Username       string    `bson:"username" json:"username"`           
	Date           time.Time `bson:"date" json:"date"`                    
	ProblemsSolved int       `bson:"problemsSolved" json:"problemsSolved"`
	IsConsistent   bool      `bson:"isConsistent" json:"isConsistent"`    
}

