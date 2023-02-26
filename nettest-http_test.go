package main

import (
	"net"
	"testing"
	"time"
)

func TestOpenPort(t *testing.T) {
	go main()

	conn, err := net.DialTimeout("tcp", "localhost:8081", 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to localhost:8081: %v", err)
	}
	defer conn.Close()
}
