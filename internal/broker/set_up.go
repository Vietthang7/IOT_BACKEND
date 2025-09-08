package broker

// Tạo interface để broadcast data
type Broadcaster interface {
	BroadcastData(event string, data interface{})
	BroadcastSensorData(data interface{})
	BroadcastDeviceStatus(deviceName, action string)
}

var GlobalBroadcaster Broadcaster

func SetBroadcaster(b Broadcaster) {
	GlobalBroadcaster = b
}

func BroadcastData(event string, data interface{}) {
	if GlobalBroadcaster != nil {
		GlobalBroadcaster.BroadcastData(event, data)
	}
}

func BroadcastSensorData(data interface{}) {
	if GlobalBroadcaster != nil {
		GlobalBroadcaster.BroadcastSensorData(data)
	}
}
