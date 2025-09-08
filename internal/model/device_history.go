package model

import "time"

type DeviceHistory struct {
	Model      `gorm:"embedded"`
	DeviceName string    `json:"device_name"`
	Action     string    `json:"action"`
	Time       time.Time `json:"time"`
}
