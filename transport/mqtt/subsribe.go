package mqtt

import (
	"strings"
	"time"
)

type PublishMQTTFn func(topic string, qos byte, retain bool, payload []byte) error

func (s *Server) MakePublishMQTTFn(timeout time.Duration) PublishMQTTFn {
	return func(topic string, qos byte, retain bool, payload []byte) error {
		c := s.MQTTClient()
		t := c.Publish(topic, qos, retain, payload)
		t.WaitTimeout(timeout)
		if t.Error() != nil {
			return t.Error()
		}
		return nil
	}
}

type SubscribeMQTTFn func(topic string, qos byte) error

func (s *Server) MakeSubscribeMQTTFn() SubscribeMQTTFn {
	return func(topic string, qos byte) error {
		subscribeTopic := getSubscribeTopic(topic)
		c := s.MQTTClient()
		t := c.Subscribe(subscribeTopic, qos, s.MQTTHandler())
		t.WaitTimeout(s.timeout)
		if t.Error() != nil {
			return t.Error()
		}
		return nil
	}
}

func getSubscribeTopic(topic string) string {
	dirs := strings.Split(topic, "/")
	for i, dir := range dirs {
		if dir == "" {
			continue
		}
		if dir[0] == '{' && dir[len(dir)-1] == '}' {
			dirs[i] = "+"
		}
		if dir[0] == '*' {
			dirs[i] = "#"
		}
	}
	return strings.Join(dirs, "/")
}
