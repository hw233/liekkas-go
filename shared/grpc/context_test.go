package grpc

import (
	"context"
	"testing"
)

func TestMetadata(t *testing.T) {
	ctx := context.Background()
	ctx = WithCtxMetadata(ctx, "test", []string{"server1"})
	md, ok := GetCtxMetadata(ctx, "test")
	if !ok {
		t.Errorf("GetCtxMetadata !ok")
	}

	t.Logf("servers: %v", md.servers)
}
