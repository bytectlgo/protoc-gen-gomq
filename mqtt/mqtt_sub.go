package mqtt

import (
	"net/http"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	mqttClient pmqtt.Client
	Handler    http.Handler
}
