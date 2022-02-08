package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	pb "go-grpc-demo/proto"
	"go-grpc-demo/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type args struct {
	serverType string
	value int64
}

func getArgs() (args, error) {
	serverType := flag.String("type", "unary", "Server type [unary, server, client, bi-direct]")
	value := flag.Int64("value", 5, "Value for request")
	flag.Parse()

	if !checkServerType(*serverType) {
		return args{}, errors.New("wrong server type value")
	}

	return args{*serverType, *value}, nil
}

func checkServerType(serverType string) bool {
	serverTypeOptions := []string{"unary", "server", "client", "bi-direct"}
	for _, option := range serverTypeOptions {
		if serverType == option {
			return true
		}
	}
	return false
}

func main()  {
	// Showing useful information when the user enters the --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] \nOptions:\n", "./client")
		flag.PrintDefaults()
	}

	// get the arguments that was entered by the user
	args, err := getArgs()
	utils.CheckErr(err)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	utils.CheckErr(err)

	defer conn.Close()

	client := pb.NewProtoServiceClient(conn)

	switch args.serverType {
	case "server":
		doServer(client, args.value)
	case "client":
		doClient(client, args.value)
	case "bi-direct":
		doBiDirect(client, args.value)
	default:
		doUnary(client, args.value)
	}
}

func doUnary(client pb.ProtoServiceClient, value int64)  {
	fmt.Println("Staring to do a Unary RPC")
	request := &pb.ProtoRequest{Number: value}

	response, err := client.Square(context.Background(), request)
	utils.CheckErr(err)

	log.Printf("Response from ProtoService: %v", response.Result)
}

func doServer(client pb.ProtoServiceClient, value int64) {
	request := &pb.ProtoRequest{Number: value}
	stream, err := client.Loop(context.Background(), request)
	utils.CheckErr(err)

	done := make(chan bool)

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			utils.CheckErr(err)
			log.Printf("Response received: %d", response.Result)
		}
	}()

	<-done
}

func doClient(client pb.ProtoServiceClient, value int64) {
	requests := []*pb.ProtoRequest{
		{Number: value},
		{Number: value},
		{Number: value},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	stream, err := client.Sum(ctx)
	utils.CheckErr(err)

	for _, request := range requests {
		err := stream.Send(request)
		utils.CheckErr(err)
	}

	stream.CloseSend()

	time.Sleep(1 * time.Second)
}

func doBiDirect(client pb.ProtoServiceClient, value int64) {
	requests := []*pb.ProtoRequest{
		{Number: value},
		{Number: value},
		{Number: value},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	stream, err := client.SumAndReturn(ctx)
	utils.CheckErr(err)

	done := make(chan bool)

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			utils.CheckErr(err)
			log.Printf("Result total: %d", response.Result)
		}
	}()

	for _, request := range requests {
		err := stream.Send(request)
		utils.CheckErr(err)
	}

	stream.CloseSend()

	<-done
}