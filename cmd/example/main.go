package main

import (
	"context"
	"os"
	"os/signal"
)

func main() {
	server := NewMQTTSever(
		&Service{},
	)
	server.Start(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	server.Stop(context.Background())
}
