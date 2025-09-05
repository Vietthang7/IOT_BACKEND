package model

import "time"

type DeviceHistory struct {
	ID         int       `json:"id"`
	DeviceName string    `json:"device_name"`
	Action     string    `json:"action"`
	Time       time.Time `json:"time"`
}

type DeviceControlRequest struct {
	DeviceName string `json:"device_name"`
	Action     string `json:"action"`
}
