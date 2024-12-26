package mqtt

import (
	"fmt"
	"net/http"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	t := rw.client.Publish(rw.getTopic(), rw.getQoS(), rw.getRetain(), b)
	t.Wait()
	if t.Error() != nil {
		return http.StatusBadRequest, t.Error()
	}
	return http.StatusOK, nil
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
