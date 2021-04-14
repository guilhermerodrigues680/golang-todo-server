package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"
	clientinterceptors "todoapp/local_debug/client/client_interceptors"
	"todoapp/transport/grpc/pbtodoapp"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "Server address in the format of host:port")
)

func main() {
	flag.Parse()

	tail := flag.Args()
	logrus.Info(tail, flag.Arg(0), flag.Arg(1))
	if flag.Arg(0) == "" {
		logrus.Fatal("Necessario descricao")
	}

	description := flag.Arg(0)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithUnaryInterceptor(clientinterceptors.ClientInterceptor))
	conn, err := grpc.Dial(*serverAddr, opts...)
	// conn, err := grpc.DialContext(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pbtodoapp.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	req := pbtodoapp.TodoCreateRequest{
		Description: description,
	}

	m := make(map[string]string)
	m["route"] = "66"
	m["route2"] = "66"
	header := metadata.New(m)
	// md := metadata.Pairs("authorization", "jwtToken")
	// ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = metadata.NewOutgoingContext(ctx, header)
	// todo, err := client.Create(ctx, &req, grpc.Header(&header)) // !!nao funciona
	todo, err := client.Create(ctx, &req)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	logrus.Info(todo)

	stream0, err := client.ReadAll(ctx, &pbtodoapp.ReadAllRequest{})
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	for {
		todo, err := stream0.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			// check if is status err
			if st, ok := status.FromError(err); ok {
				log.Fatalf("Status err = %s, %d", st.Message(), st.Code())
			} else {
				log.Fatalf("%v.GetNames(_) = _, %v", client, err)
			}
		}

		log.Println(todo)
	}

	logrus.Info(stream0)
}
