package handler

import (
	"backend/internal/consts"
	"backend/internal/repo"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDeviceHistory(c *fiber.Ctx) error {
	var (
		err        error
		entry      repo.DeviceHistory
		entries    repo.List_DeviceHistory
		query      = ""
		args       = []interface{}{}
		pagination = consts.BindRequestTable(c, "time")
		conditions = []string{}
	)
	sortParam := c.Query("sort", "true")
	if sortParam == "true" {
		pagination.Dir = "desc"
	} else {
		pagination.Dir = "asc"
	}
	pagination.Order = "time"

	if c.Query("device_name") != "" {
		conditions = append(conditions, "device_name LIKE ?")
		args = append(args, c.Query("device_name"))
	}

	if c.Query("action") != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, c.Query("action"))
	}

	if c.Query("search_time") != "" {
		searchTime := c.Query("search_time")
		conditions = append(conditions, "time LIKE ?")
		args = append(args, "%"+searchTime+"%")
	}

	if len(conditions) > 0 {
		query = strings.Join(conditions, " AND ")
	}

	if entries, err = entry.Find(&pagination, query, args); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}
	pagination.Total = entry.Count(query, args)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":       entries,
		"pagination": pagination,
	})
}

func ListDevices(c *fiber.Ctx) error {
	var (
		err          error
		entry        repo.DeviceHistory
		deviceStatus []map[string]interface{}
	)

	if deviceStatus, err = entry.GetDevicesWithLatestStatus(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"devices": deviceStatus,
		"count":   len(deviceStatus),
	})
}
func GetDeviceStatsByDateRange(c *fiber.Ctx) error {
	var (
		err       error
		entry     repo.DeviceHistory
		stats     []map[string]interface{}
		startDate string
		endDate   string
	)

	// Lấy start_date và end_date từ query params
	startDate = c.Query("date")
	endDate = c.Query("date")

	// Nếu không có, mặc định lấy 7 ngày gần nhất
	if startDate == "" || endDate == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return ResponseError(c, fiber.StatusBadRequest,
			"Invalid start_date format. Use YYYY-MM-DD", consts.GetFailed)
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return ResponseError(c, fiber.StatusBadRequest,
			"Invalid end_date format. Use YYYY-MM-DD", consts.GetFailed)
	}

	if stats, err = entry.GetDeviceStatsByDateRange(startDate, endDate); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	// Group data by date for easier reading
	groupedByDate := make(map[string][]map[string]interface{})
	for _, stat := range stats {
		date := stat["date"].(time.Time).Format("2006-01-02")
		if groupedByDate[date] == nil {
			groupedByDate[date] = []map[string]interface{}{}
		}
		groupedByDate[date] = append(groupedByDate[date], map[string]interface{}{
			"device_name":   stat["device_name"],
			"on_count":      stat["on_count"],
			"off_count":     stat["off_count"],
			"total_actions": stat["total_actions"],
		})
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"start_date":    startDate,
		"end_date":      endDate,
		"stats_by_date": groupedByDate,
		"total_days":    len(groupedByDate),
	})
}

func GetTodayDeviceStats(c *fiber.Ctx) error {
	var (
		err   error
		entry repo.DeviceHistory
		stats []map[string]interface{}
	)

	if stats, err = entry.GetTodayDeviceStats(); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("%s: %s", consts.GetFail, err.Error()), consts.GetFailed)
	}

	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"date":  time.Now().Format("2006-01-02"),
		"stats": stats,
		"count": len(stats),
	})
}
