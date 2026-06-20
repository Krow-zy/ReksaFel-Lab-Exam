package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTelemetryClient_Dispatch(t *testing.T) {
	token := "correct-token"

	// 1. Create mock HTTP test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTelemetryClient(server.URL, token)

	t.Run("Successful Dispatch", func(t *testing.T) {
		err := client.Dispatch(context.Background(), "/submit", map[string]string{"foo": "bar"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("Unauthorized Dispatch", func(t *testing.T) {
		badClient := NewTelemetryClient(server.URL, "wrong-token")
		err := badClient.Dispatch(context.Background(), "/submit", map[string]string{"foo": "bar"})
		if err == nil {
			t.Error("expected error due to unauthorized token, got nil")
		}
	})

	t.Run("Context Timeout Cancellation", func(t *testing.T) {
		// Mock slow server response
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer slowServer.Close()

		slowClient := NewTelemetryClient(slowServer.URL, token)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		err := slowClient.Dispatch(ctx, "/submit", map[string]string{"foo": "bar"})
		if err == nil {
			t.Error("expected context deadline error, got nil")
		}
	})
}
