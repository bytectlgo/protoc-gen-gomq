package mq

type Client interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) error
}
