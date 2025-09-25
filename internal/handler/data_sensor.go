package handler

import (
	"backend/internal/consts"
	"backend/internal/repo"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDataSensor(c *fiber.Ctx) error {
	var (
		err        error
		entry      repo.DataSensor
		entries    repo.List_DataSensor
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
	if pagination.Search != "" {
		conditions = append(conditions, "device_name LIKE ?")
		args = append(args, "%"+pagination.Search+"%")
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

	sensorType := c.Query("sensor_type")
	var responseData interface{}

	if sensorType == "humidity" {
		// Chỉ trả về humidity data
		humidityData := make([]map[string]interface{}, len(entries))
		for i, entry := range entries {
			humidityData[i] = map[string]interface{}{
				"id":       entry.ID,
				"humidity": entry.Humidity,
				"time":     entry.Time,
			}
		}
		responseData = humidityData
	} else if sensorType == "temp" {
		// Chỉ trả về temp data
		tempData := make([]map[string]interface{}, len(entries))
		for i, entry := range entries {
			tempData[i] = map[string]interface{}{
				"id":   entry.ID,
				"temp": entry.Temp,
				"time": entry.Time,
			}
		}
		responseData = tempData
	} else if sensorType == "lux" {
		// Chỉ trả về lux data
		lightData := make([]map[string]interface{}, len(entries))
		for i, entry := range entries {
			lightData[i] = map[string]interface{}{
				"id":   entry.ID,
				"lux":  entry.Lux,
				"time": entry.Time,
			}
		}
		responseData = lightData
	} else {
		responseData = entries
	}

	pagination.Total = entry.Count(query, args)
	return ResponseSuccess(c, fiber.StatusOK, consts.GetSuccess, fiber.Map{
		"data":        responseData,
		"pagination":  pagination,
		"sensor_type": sensorType,
	})
}
