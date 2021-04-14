package transportgrpc

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LogRequest is a gRPC UnaryServerInterceptor that will log the API call to stdOut
func LogRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, err error) {

	fmt.Printf("Request for : %s\n", info.FullMethod)
	// Last but super important, execute the handler so that the actualy gRPC request is also performed
	// return handler(ctx, req)

	authorize(ctx)

	response, err = handler(ctx, req)

	return
}

// authorize function authorizes the token received from Metadata
func authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	// }

	logrus.Info(md, ok)

	return nil
}
