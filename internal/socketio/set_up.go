package socketio

import (
	"backend/internal/broker"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

var server *socketio.Server

func Setup() {
	var err error

	server = socketio.NewServer(nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// Đăng ký broadcaster để tránh circular dependency
	broker.SetBroadcaster(&SocketIOBroadcaster{})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logrus.Infof("socket.io: client %s connected", s.ID())
		return nil
	})

	server.OnEvent("/", "device_control", func(s socketio.Conn, msg map[string]interface{}) {
		deviceName, _ := msg["device"].(string)
		action, _ := msg["action"].(string)

		logrus.Infof("Received device control: device=%s, action=%s", deviceName, action)

		// Broadcast về tất cả client khác
		BroadcastData("device_status_changed", map[string]interface{}{
			"device": deviceName,
			"action": action,
			"time":   time.Now(),
		})
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		logrus.Errorf("socket.io error: %v", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logrus.Infof("socket.io: client %s disconnected: %s", s.ID(), reason)
	})

	go server.Serve()
	logrus.Info("*************** SOCKET.IO SETUP FINISHED ***************")
}

func ServeHTTP() http.Handler {
	return server
}

// Implement Broadcaster interface để tránh circular dependency
type SocketIOBroadcaster struct{}

func (s *SocketIOBroadcaster) BroadcastData(event string, data interface{}) {
	if server != nil {
		server.BroadcastToNamespace("/", event, data)
		logrus.Debugf("Broadcasting event: %s", event)
	}
}

func (s *SocketIOBroadcaster) BroadcastSensorData(data interface{}) {
	s.BroadcastData("sensor_data", data)
}

func (s *SocketIOBroadcaster) BroadcastDeviceStatus(deviceName, action string) {
	s.BroadcastData("device_status", map[string]interface{}{
		"device": deviceName,
		"action": action,
		"time":   time.Now(),
	})
}

// Local broadcast functions (vẫn giữ để backward compatibility)
func BroadcastData(event string, data interface{}) {
	if server != nil {
		server.BroadcastToNamespace("/", event, data)
	}
}

func BroadcastSensorData(data interface{}) {
	BroadcastData("sensor_data", data)
}

func BroadcastDeviceStatus(deviceName, action string) {
	BroadcastData("device_status", map[string]interface{}{
		"device": deviceName,
		"action": action,
		"time":   time.Now(),
	})
}
