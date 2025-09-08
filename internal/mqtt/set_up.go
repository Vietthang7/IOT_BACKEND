package mqtt

import (
	"backend/internal/broker"
	"backend/internal/consts"
	"backend/internal/repo"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var client mqtt.Client

func Setup() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.22.18.104:1883")
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

	logrus.Infof("Received device status confirmation: %s = %s", deviceName, payload)

	// CHỈ LƯU VÀO DATABASE KHI NHẬN ĐƯỢC XÁC NHẬN TỪ ESP32
	deviceHistory := repo.DeviceHistory{
		DeviceName: deviceName,
		Action:     payload, // ON hoặc OFF đã được xác nhận từ ESP32
		Time:       time.Now(),
	}

	if err := deviceHistory.Create(); err != nil {
		logrus.Errorf("Failed to save device history: %v", err)
	} else {
		logrus.Infof("Device history saved: %s %s", deviceName, payload)
	}

	// Broadcast trạng thái đã được xác nhận
	broker.BroadcastData("device_status_confirmed", map[string]interface{}{
		"device": deviceName,
		"action": payload,
		"time":   time.Now(),
	})
}

// Hàm gửi lệnh điều khiển đến thiết bị
func PublishCommand(deviceName, action string) error {
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

	token := client.Publish(topic, 0, false, action)
	token.Wait()

	return token.Error()
}
