package mqtt

import (
	"fmt"
	"net/http"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	MQTT_REPLY_QOS_HEADER    = "X-MQTT-Reply-QoS"
	MQTT_REPLY_RETAIN_HEADER = "X-MQTT-Reply-Retain"
	MQTT_REPLY_TOPIC_HEADER  = "X-MQTT-Reply-Topic"
)

type MQTTResponseWriter struct {
	header http.Header
	client mqtt.Client
}

func (rw *MQTTResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *MQTTResponseWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	topic := rw.getTopic()
	if topic == "" {
		log.Warnf("topic is empty, body: %s", string(b))
		return 0, nil
	}
	t := rw.client.Publish(topic, rw.getQoS(), rw.getRetain(), b)
	t.Wait()
	if t.Error() != nil {
		return 0, t.Error()
	}
	return len(b), nil
}
func (rw *MQTTResponseWriter) getQoS() byte {
	qos, err := strconv.ParseInt(rw.header.Get(MQTT_REPLY_QOS_HEADER), 10, 64)
	if err != nil {
		return 0
	}
	return byte(qos)
}
func (rw *MQTTResponseWriter) getRetain() bool {
	return rw.header.Get(MQTT_REPLY_RETAIN_HEADER) == "true"
}
func (rw *MQTTResponseWriter) getTopic() string {
	return rw.header.Get(MQTT_REPLY_TOPIC_HEADER)
}

func (rw *MQTTResponseWriter) WriteHeader(statusCode int) {
	rw.header.Set("Status", fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)))
}
