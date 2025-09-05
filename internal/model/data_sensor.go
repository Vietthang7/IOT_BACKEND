package model

import "time"

type DataSensor struct {
	ID       int       `json:"id"`
	Temp     float64   `json:"temp"`
	Humidity float64   `json:"humidity" db:"humidity"`
	Lux      int       `json:"lux" db:"lux"`
	Time     time.Time `json:"time" db:"time"`
}

type DataSensorRequest struct {
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Lux      int     `json:"lux"`
}
