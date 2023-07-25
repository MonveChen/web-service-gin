/*
 * @Author: Monve
 * @Date: 2023-07-25 07:22:50
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 08:21:29
 * @FilePath: /web-service-gin/controllers/record/record.go
 */
package record

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"web-service-gin/utils/redis"

	"github.com/gin-gonic/gin"
)

type StatsRequest struct {
	StartDate string `json:"start_date" example:"2023-07-24"`
	EndDate   string `json:"end_date" example:"2023-07-26"`
}

var ctx = context.Background()

func getAccessCountInDateRange(startDate, endDate string) (map[string]map[string]string, error) {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	dateRecord := make(map[string]map[string]string)
	for date := startTime; date.Before(endTime.Add(24 * time.Hour)); date = date.Add(24 * time.Hour) {
		key := fmt.Sprintf("access:%s", date.Format("2006-01-02"))
		data, err := redis.Db.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		dateRecord[date.Format("2006-01-02")] = data
	}

	return dateRecord, nil
}

// @BasePath /api/v1

// @Summary 查询访问情况
// @Schemes
// @Description 按（时间区间）查询访问情况
// @Tags Record
// @Accept json
// @Produce json
// @Param record query StatsRequest true "time Range"
// @Router /record/stats [get]
func StatsHandler(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Date range invail",
		})
		return
	}

	stats, err := getAccessCountInDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}
