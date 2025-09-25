package model

import "time"

type DataSensor struct {
	Model    `gorm:"embedded"`
	Temp     float64   `json:"temp,omitempty"`
	Humidity float64   `json:"humidity,omitempty" db:"humidity"`
	Lux      int       `json:"lux,omitempty" db:"lux"`
	Time     time.Time `json:"time" db:"time"`
}
