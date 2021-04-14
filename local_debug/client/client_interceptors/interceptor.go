package clientinterceptors

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Logic before invoking the invoker
	start := time.Now()
	// Calls the invoker to execute RPC

	// m := map[string]string{
	// 	"interceptorkey": "3711",
	// }

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return errors.New("deu ruim aqui")
	}

	logrus.Info(md)
	md.Append("interceptorkeyauthorization", "3711")
	logrus.Info(md)

	// header := metadata.New(m)
	// ctx = metadata.NewOutgoingContext(ctx, header)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)
	// Logic after invoking the invoker
	logrus.Infof("Invoked RPC method=%s; Duration=%s; Error=%v", method, time.Since(start), err)
	return err
}
