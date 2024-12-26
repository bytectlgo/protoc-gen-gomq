package mqtt

import (
	"strings"
)

func MakeSubscribeFn(s *Server) func(topic string, qos byte) error {
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
