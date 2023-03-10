package balancer

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Modify(ctx context.Context) context.Context {
	return context.WithValue(ctx, "test", 100)
}

func TestGRPCContext(t *testing.T) {

	ctx := context.Background()
	// md := Metadata.MD{}
	// md.SetBalance()

	err := grpc.SetHeader(ctx, metadata.Pairs("test", "100"))
	if err != nil {
		t.Errorf("error: %v", err)
	}

	t.Logf("ctx value: %v", ctx.Value("test"))
}

func TestNormalContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test1", 100)
	ctx = context.WithValue(ctx, "test2", 200)
	t.Logf("ctx value: %v", ctx.Value("test1"))
	t.Logf("ctx value: %v", ctx.Value("test2"))
}

func TestContext(t *testing.T) {
	c1 := NewContextHandler("test")
	c2 := NewContextHandler("test2")

	// b := int64(99)
	//
	// c1.SetIDSource(func() int64 {
	// 	t.Logf("call b")
	// 	return b
	// })

	c1.SetIDOnce(10)
	// c1.SetID(10)

	ctx := context.Background()
	ctx = c1.WithContext(ctx)
	ctx = c2.WithContext(ctx)

	// c.SetID(10)

	ret, ok := GetContext(ctx, "test")
	if !ok {
		t.Logf("!OK1")
	}

	t.Logf("ret1: Server [%s], ID [%d]", ret.Server, ret.ID)

	ctx = c1.WithContext(ctx)

	// b = 988
	// c.SetID(10)
	// c.SetIDOnce(10)
	ret, _ = GetContext(ctx, "test")
	t.Logf("ret2: Server [%s], ID [%d]", ret.Server, ret.ID)
}
