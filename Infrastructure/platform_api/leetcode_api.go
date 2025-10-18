// package platform_api

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"
// 	"bytes"

// 	"consistent_1/Domain" // Assuming PlatformActivity is in Domain
// )

// // LeetCodeAPIClient implements the LeetCodeAPI interface.
// type LeetCodeAPIClient struct {
// 	baseURL    string
// 	httpClient *http.Client
// }

// // NewLeetCodeAPI creates a new LeetCodeAPIClient.
// func NewLeetCodeAPI(baseURL string) LeetCodeAPI {
// 	return &LeetCodeAPIClient{
// 		baseURL:    baseURL,
// 		httpClient: &http.Client{Timeout: 10 * time.Second},
// 	}
// }

// // GraphQL query for LeetCode user activity
// const leetcodeGraphQLQuery = `
// query userSolutionNum($username: String!) {
//     matchedUser(username: $username) {
//         submitStats {
//             acSubmissionNum {
//                 difficulty
//                 count
//             }
//         }
//     }
// }
// `

// // LeetCodeGraphQLResponse defines the structure of the LeetCode GraphQL response.
// type LeetCodeGraphQLResponse struct {
// 	Data struct {
// 		MatchedUser struct {
// 			SubmitStats struct {
// 				AcSubmissionNum []struct {
// 					Difficulty string `json:"difficulty"`
// 					Count      int    `json:"count"`
// 				} `json:"acSubmissionNum"`
// 			} `json:"submitStats"`
// 		} `json:"matchedUser"`
// 	} `json:"data"`
// 	Errors []struct {
// 		Message string `json:"message"`
// 	} `json:"errors"`
// }

// // FetchUserDailyActivity fetches the daily problems solved for a LeetCode user.
// // NOTE: LeetCode's public API does not provide *daily* submission counts directly.
// // This implementation will fetch *total* accepted submissions and assume a "consistency"
// // based on whether they have *any* accepted submissions, or if you had a previous
// // mechanism to fetch daily changes (which is usually not publicly available without scraping/private APIs).
// // For a true "daily solved" count, you'd need to:
// // 1. Store previous day's total.
// // 2. Fetch current day's total.
// // 3. Subtract to get daily solved.
// //
// // For simplicity and as a demonstration, we will consider a user "consistent" if their total
// // accepted submissions count is greater than 0, implying some activity.
// // A more accurate approach would involve fetching the user's recent submissions feed.
// func (api *LeetCodeAPIClient) FetchUserDailyActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error) {
// 	requestBody := map[string]interface{}{
// 		"query": leetcodeGraphQLQuery,
// 		"variables": map[string]interface{}{
// 			"username": username,
// 		},
// 		"operationName": "userSolutionNum",
// 	}
// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return domain.PlatformActivity{}, fmt.Errorf("failed to marshal LeetCode GraphQL request: %w", err)
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		return domain.PlatformActivity{}, fmt.Errorf("failed to create LeetCode GraphQL request: %w", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("User-Agent", "Consistify-Backend/1.0") // Good practice to set a User-Agent

// 	resp, err := api.httpClient.Do(req)
// 	if err != nil {
// 		return domain.PlatformActivity{}, fmt.Errorf("failed to make LeetCode GraphQL request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		respBody, _ := ioutil.ReadAll(resp.Body)
// 		log.Printf("LeetCode API error response (%d): %s", resp.StatusCode, string(respBody))
// 		return domain.PlatformActivity{}, fmt.Errorf("LeetCode API responded with status %d", resp.StatusCode)
// 	}

// 	var graphQLResp LeetCodeGraphQLResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
// 		return domain.PlatformActivity{}, fmt.Errorf("failed to decode LeetCode GraphQL response: %w", err)
// 	}

// 	if len(graphQLResp.Errors) > 0 {
// 		return domain.PlatformActivity{}, fmt.Errorf("LeetCode GraphQL error: %s", graphQLResp.Errors[0].Message)
// 	}

// 	if graphQLResp.Data.MatchedUser.SubmitStats.AcSubmissionNum == nil {
// 		return domain.PlatformActivity{}, fmt.Errorf("LeetCode user '%s' not found or no submission data", username)
// 	}

// 	totalProblemsSolved := 0
// 	for _, stat := range graphQLResp.Data.MatchedUser.SubmitStats.AcSubmissionNum {
// 		totalProblemsSolved += stat.Count
// 	}

// 	// Determine daily consistency: For LeetCode, this is challenging without daily history.
// 	// As a placeholder, we'll mark consistent if *any* problems have been solved overall.
// 	// For a real-world app, you'd need to track daily changes or scrape submission history.
// 	// Let's assume for *this demonstration* that any `totalProblemsSolved > 0`
// 	// means they are "active" for the purpose of the initial fetch.
// 	// The true daily check would need more sophisticated logic (e.g., check last submission timestamp).
// 	isConsistent := totalProblemsSolved > 0

// 	// THIS IS THE CRITICAL PART: `problemsSolvedToday` must be defined.
// 	// Since LeetCode GraphQL doesn't give us daily count directly, we'll make a simplifying assumption.
// 	// For a more accurate "ProblemsSolved" for the day, you'd need to:
// 	// 1. Fetch the user's recent submissions.
// 	// 2. Filter those submissions to only include ones made on `date`.
// 	// 3. Count them.
// 	// For now, let's represent `problemsSolvedToday` as a symbolic 1 if `isConsistent`, else 0.
// 	problemsSolvedToday := 0
// 	if isConsistent {
// 		problemsSolvedToday = 1 // Simplified: If they have any submissions, count 1 for today (needs improvement for accuracy)
// 	}


// 	return domain.PlatformActivity{
// 		Platform:       "leetcode",
// 		Username:       username,
// 		Date:           date.Truncate(24 * time.Hour), // Store date at start of day UTC
// 		IsConsistent:   isConsistent,
// 		ProblemsSolved: problemsSolvedToday, // Correctly assigned here
// 	}, nil
// }




package platform_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"bytes"

	"consistent_1/Domain"
)


type LeetCodeAPIClient struct {
	baseURL    string
	httpClient *http.Client
}


func NewLeetCodeAPI(baseURL string) LeetCodeAPI {
	return &LeetCodeAPIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}


const leetcodeGraphQLQuery = `
query userSolutionNum($username: String!) {
    matchedUser(username: $username) {
        submitStats {
            acSubmissionNum {
                difficulty
                count
            }
        }
    }
}
`
type LeetCodeGraphQLResponse struct {
	Data struct {
		MatchedUser struct {
			SubmitStats struct {
				AcSubmissionNum []struct {
					Difficulty string `json:"difficulty"`
					Count      int    `json:"count"`
				} `json:"acSubmissionNum"`
			} `json:"submitStats"`
		} `json:"matchedUser"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (api *LeetCodeAPIClient) FetchUserDailyActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error) {
	requestBody := map[string]interface{}{
		"query": leetcodeGraphQLQuery,
		"variables": map[string]interface{}{
			"username": username,
		},
		"operationName": "userSolutionNum",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to marshal LeetCode GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to create LeetCode GraphQL request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Consistify-Backend/1.0")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to make LeetCode GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("LeetCode API error response (%d) for user %s: %s", resp.StatusCode, username, string(respBody))
		return domain.PlatformActivity{}, fmt.Errorf("LeetCode API responded with status %d", resp.StatusCode)
	}

	var graphQLResp LeetCodeGraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to decode LeetCode GraphQL response: %w", err)
	}

	if len(graphQLResp.Errors) > 0 {
		return domain.PlatformActivity{}, fmt.Errorf("LeetCode GraphQL error: %s", graphQLResp.Errors[0].Message)
	}

	if graphQLResp.Data.MatchedUser.SubmitStats.AcSubmissionNum == nil {
		return domain.PlatformActivity{}, fmt.Errorf("LeetCode user '%s' not found or no submission data", username)
	}

	totalProblemsSolved := 0
	for _, stat := range graphQLResp.Data.MatchedUser.SubmitStats.AcSubmissionNum {
		totalProblemsSolved += stat.Count
	}
	return domain.PlatformActivity{
		Platform:       "leetcode",
		Username:       username,
		Date:           date.Truncate(24 * time.Hour),
		IsConsistent:   false, 
		ProblemsSolved: totalProblemsSolved, 
	}, nil
}