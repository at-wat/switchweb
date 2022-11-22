package main

import (
	"context"
	"log"
)

func main() {
	log.SetFlags(0)

	cli := newClient()
	s := newServer(cli)
	s.updateDevices(context.TODO())

	s.start()
}
