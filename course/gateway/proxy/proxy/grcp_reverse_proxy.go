package proxy

import (
	"context"
	"log"

	"go-examples/course/gateway/proxy/load_balance"

	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			return ctx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
