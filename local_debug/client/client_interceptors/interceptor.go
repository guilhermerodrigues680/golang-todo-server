package clientinterceptors

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ClientInterceptor(logger *logrus.Logger) func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Logic before invoking the invoker
		start := time.Now()
		// Calls the invoker to execute RPC

		// m := map[string]string{
		// 	"interceptorkey": "3711",
		// }

		// md, ok := metadata.FromOutgoingContext(ctx)
		// if !ok {
		// 	return errors.New("deu ruim aqui")
		// }

		// logrus.Info(md)
		// md.Append("interceptorkeyauthorization", "3711")
		// logrus.Info(md)

		// header := metadata.New(m)
		// ctx = metadata.NewOutgoingContext(ctx, header)
		// ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		// Logic after invoking the invoker
		logger.Debugf("Invoked RPC Unary method=%s; Duration=%s; Error=%v", method, time.Since(start), err)
		return err
	}
}

func StreamClientInterceptor(logger *logrus.Logger) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		logger.Debugf("stream %s", method)
		// Logic before invoking the invoker
		start := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		// Logic after invoking the invoker
		logger.Debugf("Invoked RPC Stream method=%s; Duration=%s; Error=%v", method, time.Since(start), err)
		return clientStream, err
	}
}
