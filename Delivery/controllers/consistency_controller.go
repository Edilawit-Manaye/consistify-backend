package controllers

import (
	"log"
	"net/http"
	"time"

	"consistent_1/Domain"
	 "consistent_1/Usecases"  

	"github.com/gin-gonic/gin"
)


type ConsistencyController struct {
	consistencyUsecase usecases.ConsistencyUsecase
}
func NewConsistencyController(consistencyUsecase usecases.ConsistencyUsecase) *ConsistencyController {
	return &ConsistencyController{
		consistencyUsecase: consistencyUsecase,
	}
}


func (ctrl *ConsistencyController) GetDailyConsistency(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	dateStr := c.Query("date") 

	var queryDate time.Time
	if dateStr == "" {
		queryDate = time.Now() 
	} else {
		var err error
		queryDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected YYYY-MM-DD"})
			return
		}
	}

	consistency, err := ctrl.consistencyUsecase.GetDailyConsistency(c.Request.Context(), userID, queryDate)
	if err != nil {
		switch err {
		case domain.ErrConsistencyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "No consistency record found for this date"})
		default:
			log.Printf("Error getting daily consistency for user %s, date %s: %v", userID, queryDate.Format("2006-01-02"), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily consistency"})
		}
		return
	}

	c.JSON(http.StatusOK, consistency)
}


func (ctrl *ConsistencyController) GetConsistencyHistory(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	startDateStr := c.Query("startDate") 
	endDateStr := c.Query("endDate")    

	var startDate, endDate *time.Time
	var err error

	if startDateStr != "" {
		t, e := time.Parse("2006-01-02", startDateStr)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Expected YYYY-MM-DD"})
			return
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, e := time.Parse("2006-01-02", endDateStr)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Expected YYYY-MM-DD"})
			return
		}
		endDate = &t
	}

	history, err := ctrl.consistencyUsecase.GetConsistencyHistory(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting consistency history for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve consistency history"})
		return
	}

	c.JSON(http.StatusOK, history)
}


func (ctrl *ConsistencyController) GetUserStreaks(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	streakInfo, err := ctrl.consistencyUsecase.GetStreaks(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting streaks for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve streak information"})
		return
	}

	c.JSON(http.StatusOK, streakInfo)
}


func (ctrl *ConsistencyController) TriggerDailyConsistencyCheck(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	consistency, err := ctrl.consistencyUsecase.CheckDailyConsistency(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error manually triggering consistency check for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger consistency check"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily consistency check triggered successfully", "consistency": consistency})
}