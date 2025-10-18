
package platform_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"consistent_1/Domain" 
)


type CodeforcesAPIClient struct {
	baseURL    string
	httpClient *http.Client
}
func NewCodeforcesAPI(baseURL string) CodeforcesAPI {
	return &CodeforcesAPIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
type CodeforcesSubmission struct {
	ID                  int         `json:"id"`
	ContestID           int         `json:"contestId"`
	CreationTimeSeconds int64       `json:"creationTimeSeconds"`
	Problem             struct {
		ContestID int    `json:"contestId"`
		Index     string `json:"index"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Points    float64 `json:"points"`
		Rating    int    `json:"rating"`
		Tags      []string `json:"tags"`
	} `json:"problem"`
	Author struct {
		ContestID int      `json:"contestId"`
		Members   []struct {
			Handle string `json:"handle"`
		} `json:"members"`
		ParticipantType string `json:"participantType"`
		Ghost           bool   `json:"ghost"`
		StartTimeSeconds int64  `json:"startTimeSeconds"`
	} `json:"author"`
	ProgrammingLanguage string `json:"programmingLanguage"`
	Verdict             string `json:"verdict"` 
	Testset             string `json:"testset"`
	PassedTestCount     int    `json:"passedTestCount"`
	TimeConsumedMillis  int    `json:"timeConsumedMillis"`
	MemoryConsumedBytes int    `json:"memoryConsumedBytes"`
}
type CodeforcesUserStatusResponse struct {
	Status string                 `json:"status"` 
	Result []CodeforcesSubmission `json:"result"`
	Comment string                `json:"comment,omitempty"`
}
func (api *CodeforcesAPIClient) FetchUserDailyActivity(ctx context.Context, username string, date time.Time) (domain.PlatformActivity, error) {
	url := fmt.Sprintf("%s/api/user.status?handle=%s&from=1&count=50", api.baseURL, username)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to create Codeforces request: %w", err)
	}
	req.Header.Set("User-Agent", "Consistify-Backend/1.0") 

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to make Codeforces request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Codeforces API error response (%d): %s", resp.StatusCode, string(respBody))
		return domain.PlatformActivity{}, fmt.Errorf("Codeforces API responded with status %d: %s", resp.StatusCode, string(respBody))
	}

	var cfResp CodeforcesUserStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return domain.PlatformActivity{}, fmt.Errorf("failed to decode Codeforces response: %w", err)
	}

	if cfResp.Status != "OK" {
		return domain.PlatformActivity{}, fmt.Errorf("Codeforces API error: %s", cfResp.Comment)
	}

	problemsSolvedToday := 0
	isConsistent := false
	uniqueProblemIDsToday := make(map[string]bool) 
	startOfDayUTC := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDayUTC := startOfDayUTC.Add(24 * time.Hour).Add(-time.Second) 

	for _, sub := range cfResp.Result {
		submissionTime := time.Unix(sub.CreationTimeSeconds, 0).UTC()
		if submissionTime.After(startOfDayUTC) && submissionTime.Before(endOfDayUTC) && sub.Verdict == "OK" {
			problemIdentifier := fmt.Sprintf("%d-%s", sub.Problem.ContestID, sub.Problem.Index)
			if !uniqueProblemIDsToday[problemIdentifier] {
				problemsSolvedToday++
				uniqueProblemIDsToday[problemIdentifier] = true
			}
			isConsistent = true 
		}
	}

	return domain.PlatformActivity{
		Platform:       "codeforces",
		Username:       username,
		Date:           date.Truncate(24 * time.Hour), 
		IsConsistent:   isConsistent,
		ProblemsSolved: problemsSolvedToday,
	}, nil
}