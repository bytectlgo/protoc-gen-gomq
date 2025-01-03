package mqtt

import (
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
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
	subscribeFn := func(topic string, qos byte) {
		for {
			interval := 5 * time.Second
			failCount := 0
			maxInterval := 60 * time.Second
			sussessResult := true
			subscribeTopic := getSubscribeTopic(topic)
			c := s.MQTTClient()
			t := c.Subscribe(subscribeTopic, qos, s.MQTTHandler())
			waitFlag := t.WaitTimeout(s.timeout)
			err := t.Error()
			if waitFlag {
				if err == nil {
					sToken := t.(*mqtt.SubscribeToken)
					for k, v := range sToken.Result() {
						if v != byte(qos) {
							sussessResult = false
							log.Errorf("subscribe topic(%s) failed result: %d", k, v)
						}
					}
				} else {
					sussessResult = false
				}
			} else {
				sussessResult = false
				log.Warnf("subscribe topic(%s) wait timeout", topic)
			}
			if sussessResult {
				// 订阅成功, 跳出循环
				log.Infof("subscribe topic(%s) success", topic)
				break
			}
			// 订阅失败, 重试
			if !c.IsConnected() {
				// 返回等下次重连
				log.Errorf("mqtt client is disconnected, subscribe topic(%s) error", topic)
				return
			}
			// 重试
			failCount++
			waitInterval := interval + time.Duration(failCount)*time.Second
			if waitInterval > maxInterval {
				waitInterval = maxInterval
			}
			time.Sleep(waitInterval)
		}
		return
	}

	return func(topic string, qos byte) error {
		// 异步订阅
		go subscribeFn(topic, qos)
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
