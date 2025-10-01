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
	if c.Query("sort_type") != "" && c.Query("sort_order") != "" {
		sortType := c.Query("sort_type")
		sortOrder := c.Query("sort_order")

		validSortColumns := map[string]bool{
			"temp":     true,
			"humidity": true,
			"lux":      true,
			"time":     true,
		}

		validSortOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		if validSortColumns[sortType] && validSortOrders[sortOrder] {
			pagination.Order = sortType
			pagination.Dir = sortOrder
			logrus.Infof("Custom sort applied: %s %s", sortType, sortOrder)
		} else {
			logrus.Warnf("Invalid sort parameters: sort_type=%s, sort_order=%s", sortType, sortOrder)
		}
	}
	sensorType := c.Query("sensor_type")
	if c.Query("search_data") != "" {
		searchData := c.Query("search_data")
		if sensorType == "all" {
			conditions = append(conditions, "(temp LIKE ? OR humidity LIKE ? OR lux LIKE ? OR DATE_FORMAT(time, '%Y-%m-%d %H:%i:%s') LIKE ?)")
			args = append(args, "%"+searchData+"%", "%"+searchData+"%", "%"+searchData+"%", "%"+searchData+"%")
		} else {
			switch sensorType {
			case "temp":
				conditions = append(conditions, "(temp LIKE ? OR DATE_FORMAT(time, '%Y-%m-%d %H:%i:%s') LIKE ?)")
				args = append(args, "%"+searchData+"%", "%"+searchData+"%")
			case "humidity":
				conditions = append(conditions, "(humidity LIKE ? OR DATE_FORMAT(time, '%Y-%m-%d %H:%i:%s') LIKE ?)")
				args = append(args, "%"+searchData+"%", "%"+searchData+"%")
			case "lux":
				conditions = append(conditions, "(lux LIKE ? OR DATE_FORMAT(time, '%Y-%m-%d %H:%i:%s') LIKE ?)")
				args = append(args, "%"+searchData+"%", "%"+searchData+"%")
			default:
				conditions = append(conditions, "DATE_FORMAT(time, '%Y-%m-%d %H:%i:%s') LIKE ?")
				args = append(args, "%"+searchData+"%")
			}
		}
	}
	if c.Query("sort_type") != "" {
		if c.Query("sort_order") != "" {

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
		"data":        entries,
		"pagination":  pagination,
		"sensor_type": sensorType,
	})
}
