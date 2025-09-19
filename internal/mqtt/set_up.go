package mqtt

import (
	"backend/internal/broker"
	"backend/internal/consts"
	"backend/internal/repo"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var client mqtt.Client

// Map để theo dõi các lệnh chờ phản hồi
var pendingCommands = make(map[string]chan string)
var pendingMutex = sync.Mutex{}

func Setup() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.22.18.102:1883")
	opts.SetUsername("user1")
	opts.SetPassword("123456")
	opts.SetClientID("backend-server")

	opts.SetDefaultPublishHandler(messageHandler)

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Error("MQTT connection error:", token.Error())
		return
	}

	// Subscribe vào các topics
	if token := client.Subscribe(consts.TOPIC_DATASENSOR, 0, nil); token.Wait() && token.Error() != nil {
		logrus.Error("MQTT subscription error:", token.Error())
		return
	}

	if token := client.Subscribe(consts.TOPIC_DEN_STATUS, 0, nil); token.Wait() && token.Error() != nil {
		logrus.Error("MQTT subscription error:", token.Error())
		return
	}

	if token := client.Subscribe(consts.TOPIC_QUAT_STATUS, 0, nil); token.Wait() && token.Error() != nil {
		logrus.Error("MQTT subscription error:", token.Error())
		return
	}

	if token := client.Subscribe(consts.TOPIC_DIEUHOA_STATUS, 0, nil); token.Wait() && token.Error() != nil {
		logrus.Error("MQTT subscription error:", token.Error())
		return
	}

	logrus.Info("MQTT client connected")
}

// Xử lý tin nhắn nhận được
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	fmt.Println(topic)
	jsonData, _ := json.MarshalIndent(payload, "", "  ")
	// In chuỗi JSON
	fmt.Println(string(jsonData))
	logrus.Infof("Received message: %s from topic: %s", payload, topic)

	if topic == consts.TOPIC_DATASENSOR {
		processSensorData(payload)
	} else if strings.Contains(topic, "Status") {
		processDeviceStatus(topic, payload)
	}

	broker.BroadcastData(topic, payload)
}

func processSensorData(payload string) {
	var temp, humidity float64
	var light int

	fmt.Sscanf(payload, "Temperature: %f *C, Humidity: %f %%, Light: %d", &temp, &humidity, &light)

	dataSensor := repo.DataSensor{
		Temp:     temp,
		Humidity: humidity,
		Lux:      light,
		Time:     time.Now(),
	}

	dataSensor.Create()
	broker.BroadcastSensorData(map[string]interface{}{
		"temp":     temp,
		"humidity": humidity,
		"lux":      light,
		"time":     time.Now(),
	})
}

func processDeviceStatus(topic, payload string) {
	var deviceName string

	if strings.Contains(topic, "denStatus") {
		deviceName = consts.DEVICE_DEN
	} else if strings.Contains(topic, "quatStatus") {
		deviceName = consts.DEVICE_QUAT
	} else if strings.Contains(topic, "dieuhoaStatus") {
		deviceName = consts.DEVICE_DIEUHOA
	}

	// Validation payload
	if payload != consts.ACTION_ON && payload != consts.ACTION_OFF {
		logrus.Errorf("Invalid device status: %s from topic: %s", payload, topic)
		return
	}

	// ✅ THÊM LOGIC GỬI PHẢN HỒI VỀ CHANNEL
	commandKey := fmt.Sprintf("%s:%s", deviceName, payload)

	pendingMutex.Lock()
	if responseChan, exists := pendingCommands[commandKey]; exists {
		// Gửi phản hồi qua channel
		select {
		case responseChan <- payload:
			logrus.Infof("Response sent to waiting command: %s", commandKey)
		default:
			logrus.Warnf("Response channel full for command: %s", commandKey)
		}
	}
	pendingMutex.Unlock()

	// Lưu vào database
	deviceHistory := repo.DeviceHistory{
		DeviceName: deviceName,
		Action:     payload,
		Time:       time.Now(),
	}

	if err := deviceHistory.Create(); err != nil {
		logrus.Errorf("Failed to save device history: %v", err)
	}

	// Broadcast trạng thái đã được xác nhận
	broker.BroadcastData("device_status_confirmed", map[string]interface{}{
		"device": deviceName,
		"action": payload,
		"time":   time.Now(),
	})
}

// Hàm gửi lệnh điều khiển đến thiết bị
func PublishCommandAndWait(ctx context.Context, deviceName, action string) error {
	var topic string

	switch deviceName {
	case consts.DEVICE_DEN:
		topic = consts.TOPIC_DEN
	case consts.DEVICE_QUAT:
		topic = consts.TOPIC_QUAT
	case consts.DEVICE_DIEUHOA:
		topic = consts.TOPIC_DIEUHOA
	case consts.DEVICE_ALL:
		topic = consts.TOPIC_TURNALL
	default:
		return fmt.Errorf("unknown device: %s", deviceName)
	}

	// Tạo key để theo dõi lệnh
	commandKey := fmt.Sprintf("%s:%s", deviceName, action)

	// Tạo channel để nhận phản hồi
	responseChan := make(chan string, 1)

	pendingMutex.Lock()
	pendingCommands[commandKey] = responseChan
	pendingMutex.Unlock()

	// Cleanup channel khi function kết thúc
	defer func() {
		pendingMutex.Lock()
		delete(pendingCommands, commandKey)
		pendingMutex.Unlock()
		close(responseChan)
	}()

	// Gửi lệnh
	token := client.Publish(topic, 0, false, action)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Chờ phản hồi hoặc timeout
	select {
	case response := <-responseChan:
		if response == action {
			logrus.Infof("Device confirmed: %s = %s", deviceName, response)
			return nil
		} else {
			return fmt.Errorf("device response mismatch: expected %s, got %s", action, response)
		}
	case <-ctx.Done():
		return ctx.Err() // Timeout hoặc cancel
	}
}
