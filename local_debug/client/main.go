package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	clientinterceptors "todoapp/local_debug/client/client_interceptors"
	"todoapp/local_debug/client/commands"
	"todoapp/transport/grpc/pbtodoapp"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	// log tudo
	logger.SetLevel(logrus.TraceLevel)
	// Output to stdout instead of the default stderr
	logger.SetOutput(os.Stdout)
}

func main() {
	serverAddr := flag.String("server_addr", "localhost:10000", "Server address in the format of host:port")
	flag.Parse()

	if len(flag.Args()) == 0 {
		logger.Fatal("invalid call expected subcommand")
	}

	logger.Debug(os.Args)
	logger.Debug(flag.Args())

	//----gRPC Client
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(clientinterceptors.ClientInterceptor(logger)),
		grpc.WithStreamInterceptor(clientinterceptors.StreamClientInterceptor(logger)),
	}

	logger.Info("Creating a client connection...")
	// conn, err := grpc.DialContext(*serverAddr, opts...)
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}

	logger.Info("Connection created successfully!")
	defer conn.Close()
	client := pbtodoapp.NewTodoServiceClient(conn)
	//----gRPC Client

	subcommand := flag.Arg(0)
	switch subcommand {
	case commands.Create.Cmd():
		create(client)
	case commands.Readall.Cmd():
		readAll(client)
	case commands.Read.Cmd():
		read(client)
	case commands.DeleteMultiple.Cmd():
		deleteMultiple(client)
	default:
		logger.Fatalf("command '%s' not defined", subcommand)
	}

	os.Exit(0)
}

func create(client pbtodoapp.TodoServiceClient) {
	createCmd := flag.NewFlagSet(commands.Create.Cmd(), flag.ExitOnError)
	createCmd.Parse(flag.Args()[1:]) // ignora o primero arg que é o proprio comando

	logger.Trace(createCmd.Name(), createCmd.Args(), createCmd.Arg(0), createCmd.Arg(1), createCmd.Arg(2))

	description := createCmd.Arg(0)
	if description == "" {
		logger.Fatal("Necessario descricao")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := pbtodoapp.TodoCreateRequest{
		Description: description,
	}

	todo, err := client.Create(ctx, &req)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}

	logger.Info(todo)
}

func readAll(client pbtodoapp.TodoServiceClient) {
	// readCmd    := flag.NewFlagSet(commands.Readall.Cmd(), flag.ExitOnError)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream0, err := client.ReadAll(ctx, &pbtodoapp.ReadAllRequest{})
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
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
				log.Fatalf("%v, %v", client, err)
			}
		}

		log.Println(todo)
	}

	logrus.Info(stream0)
}

func read(client pbtodoapp.TodoServiceClient) {
	readCmd := flag.NewFlagSet(commands.Read.Cmd(), flag.ExitOnError)
	readCmd.Parse(flag.Args()[1:]) // ignora o primero arg que é o proprio comando

	logger.Trace(readCmd.Name(), readCmd.Args())

	idstr := readCmd.Arg(0)
	if idstr == "" {
		logger.Fatal("required id")
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		logger.Fatal("invalid id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reqId := pbtodoapp.Id{
		Id: uint64(id),
	}

	todo, err := client.Read(ctx, &reqId)
	if err != nil {
		logger.Fatalf("fail: %v", err)
	}

	logger.Info(todo)
}

func deleteMultiple(client pbtodoapp.TodoServiceClient) {
	deleteMultipleCmd := flag.NewFlagSet(commands.DeleteMultiple.Cmd(), flag.ExitOnError)
	deleteMultipleCmd.Parse(flag.Args()[1:]) // ignora o primero arg que é o proprio comando

	logger.Trace(deleteMultipleCmd.Name(), deleteMultipleCmd.Args())

	idStrList := deleteMultipleCmd.Args()
	if len(idStrList) == 0 {
		logger.Fatal("required ids")
	}

	idList := make([]int, 0)
	for _, idStr := range idStrList {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Fatalf("invalid id: '%s'", idStr)
		}
		idList = append(idList, id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := client.DeleteMultiple(ctx)
	if err != nil {
		logger.Fatalf("fail: %v", err)
	}

	for _, id := range idList {
		reqId := pbtodoapp.Id{
			Id: uint64(id),
		}

		err := stream.Send(&reqId)
		if err != nil {
			logger.Fatal(err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("OK!")
}
