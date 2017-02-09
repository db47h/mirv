package gdb_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/db47h/mirv/gdb"
)

func Test_StartGDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := gdb.StartGDBAgent(ctx, ":42424", nil)
	if err != nil {
		t.Fatalf("Cannot create listener: %v", err)
	}
	client, err := net.Dial("tcp", "localhost:42424")
	if err != nil {
		t.Fatalf("Cannot create client: %v", err)
	}
	go func() {
		for {
			client.Write([]byte(""))
		}
	}()
	time.Sleep(200 * time.Millisecond)
	cancel()
	time.Sleep(200 * time.Millisecond)
}
