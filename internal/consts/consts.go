package consts

import "time"

const (
	// Device Names
	DEVICE_DEN     = "den"
	DEVICE_QUAT    = "quat"
	DEVICE_CUA     = "cua"
	DEVICE_CHUONG  = "chuong"
	DEVICE_DIEUHOA = "dieuhoa"
	DEVICE_ALL     = "all"

	// Actions
	ACTION_ON  = "ON"
	ACTION_OFF = "OFF"

	// MQTT Topics
	TOPIC_DATASENSOR     = "esp32/datasensor"
	TOPIC_DEN            = "esp32/den"
	TOPIC_CUA            = "esp32/cua"
	TOPIC_CHUONG         = "esp32/chuong"
	TOPIC_QUAT           = "esp32/quat"
	TOPIC_DIEUHOA        = "esp32/dieuhoa"
	TOPIC_TURNALL        = "esp32/turnall"
	TOPIC_DEN_STATUS     = "esp32/denStatus"
	TOPIC_CHUONG_STATUS  = "esp32/chuongStatus"
	TOPIC_CUA_STATUS     = "esp32/cuaStatus"
	TOPIC_QUAT_STATUS    = "esp32/quatStatus"
	TOPIC_DIEUHOA_STATUS = "esp32/dieuhoaStatus"
	MQTT_CONNECTED       = "connected"

	// TimeOut
	TIMEOUT = 2 * time.Second
)
