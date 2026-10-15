package grpc

import (
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
)

type mockRegistrar struct {
	called atomic.Bool
}

func (m *mockRegistrar) RegisterGRPC(server *grpc.Server) {
	m.called.Store(true)
}

func TestGRPCServer_Lifecycle(t *testing.T) {
	logger := slog.Default()
	s := NewGRPCServer(logger)

	reg := &mockRegistrar{}
	s.Register(reg)

	errChan := make(chan error, 1)
	go func() {
		// Use port 0 to bind to any free port dynamically.
		errChan <- s.Start("127.0.0.1:0")
	}()

	// Wait a short duration to ensure server starts and calls register.
	time.Sleep(100 * time.Millisecond)

	if !reg.called.Load() {
		t.Error("Expected RegisterGRPC to be called during Start()")
	}

	// Stop the server.
	s.Stop()

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Server did not stop within 2 seconds")
	}
}
