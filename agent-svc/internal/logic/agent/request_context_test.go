package agent

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestConversationPersistenceContextDetachesCancellation(t *testing.T) {
	t.Parallel()

	parent := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-user-id", "42", "x-request-id", "req-1"))
	parent, cancel := context.WithCancel(parent)

	detached := conversationPersistenceContext(parent)
	cancel()

	if err := detached.Err(); err != nil {
		t.Fatalf("expected detached context to ignore cancellation, got %v", err)
	}

	md, ok := metadata.FromIncomingContext(detached)
	if !ok {
		t.Fatal("expected detached context to preserve incoming metadata")
	}
	if got := md.Get("x-user-id"); len(got) != 1 || got[0] != "42" {
		t.Fatalf("expected x-user-id metadata to be preserved, got %v", got)
	}
	if got := md.Get("x-request-id"); len(got) != 1 || got[0] != "req-1" {
		t.Fatalf("expected x-request-id metadata to be preserved, got %v", got)
	}
}
