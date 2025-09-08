package handler

import (
	"backend/internal/consts"
	"backend/internal/repo"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetAllDeviceHistory(c *fiber.Ctx) error {
	var (
		err        error
		entry      repo.DeviceHistory
		entries    repo.List_DeviceHistory
		query      = ""
		args       = []interface{}{}
		pagination = consts.BindRequestTable(c, "time")
	)
	if pagination.Search != "" {
		query += "device_name LIKE ?"
		args = append(args, "%"+pagination.Search+"%")
	}
	if c.Query("start_time") != "" {
		startTime, _ := time.Parse("02-01-2006", c.Query("startTime"))
		query += " AND time > ?"
		args = append(args, startTime)
	}
	if c.Query("end_time") != "" {
		endTime, _ := time.Parse("02-01-2006", c.Query("endTime"))
		query += " AND time < ?"
		args = append(args, endTime.Add(24*time.Hour))
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
