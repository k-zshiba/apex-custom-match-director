package liveapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/koshi/apex-custom-match-director/internal/domain"
)

type Handler func(domain.LiveEvent)

type Receiver struct {
	server  *http.Server
	address string
}

func NewReceiver(port int, handler Handler) *Receiver {
	mux := http.NewServeMux()
	receiver := &Receiver{address: fmt.Sprintf("127.0.0.1:%d", port)}
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/liveapi", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "read event", http.StatusBadRequest)
			return
		}
		event, err := ParseEvent(body)
		if err != nil {
			http.Error(w, "parse event", http.StatusBadRequest)
			return
		}
		handler(event)
		w.WriteHeader(http.StatusNoContent)
	})
	receiver.server = &http.Server{
		Addr:              receiver.address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return receiver
}

func (r *Receiver) Start() error {
	listener, err := net.Listen("tcp", r.address)
	if err != nil {
		return fmt.Errorf("start liveapi receiver: %w", err)
	}
	go func() {
		if err := r.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// The owning service exposes request-time errors through state; startup errors are returned synchronously.
		}
	}()
	return nil
}

func (r *Receiver) Shutdown(ctx context.Context) error {
	if r == nil || r.server == nil {
		return nil
	}
	return r.server.Shutdown(ctx)
}

func (r *Receiver) Address() string {
	if r == nil {
		return ""
	}
	return r.address
}
