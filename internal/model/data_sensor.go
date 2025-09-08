package model

import "time"

type DataSensor struct {
	Model       `gorm:"embedded"`
	Temp     float64   `json:"temp"`
	Humidity float64   `json:"humidity" db:"humidity"`
	Lux      int       `json:"lux" db:"lux"`
	Time     time.Time `json:"time" db:"time"`
}
