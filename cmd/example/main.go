package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("hello")
	server := NewMQTTSever(
		&Service{},
	)
	server.Start(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	server.Stop(context.Background())
}
