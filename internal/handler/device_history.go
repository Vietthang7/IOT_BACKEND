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

	if c.Query("start_time") != "" {
		startTime, err := time.Parse("2006-01-02T15:04:05.000Z", c.Query("start_time"))
		if err == nil {
			conditions = append(conditions, "time >= ?")
			args = append(args, startTime)
		}
	}

	if c.Query("end_time") != "" {
		endTime, err := time.Parse("2006-01-02T15:04:05.000Z", c.Query("end_time"))
		if err == nil {
			conditions = append(conditions, "time <= ?")
			args = append(args, endTime)
		}
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
