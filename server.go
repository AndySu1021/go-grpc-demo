package main

import (
	"context"
	"fmt"
	pb "go-grpc-demo/proto"
	"go-grpc-demo/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {}

func (s server) Square(ctx context.Context, request *pb.ProtoRequest) (response *pb.ProtoResponse, err error) {
	fmt.Printf("Square function is invoked with %v \n", request)
	number := request.GetNumber()
	response = &pb.ProtoResponse{
		Result: number * number,
	}
	return
}

func (s server) Sum(sumServer pb.ProtoService_SumServer) error {
	done := make(chan bool)
	var total int64
	go func() {
		for {
			request, err := sumServer.Recv()
			fmt.Printf("Sum function is invoked with %v \n", request)
			if err == io.EOF {
				done <- true
				return
			}
			utils.CheckErr(err)
			total += request.Number
		}
	}()

	<-done

	fmt.Printf("Total is %d \n", total)

	return nil
}

func (s server) Loop(request *pb.ProtoRequest, loopServer pb.ProtoService_LoopServer) error {
	fmt.Printf("Loop function is invoked with %v \n", request)

	var wg sync.WaitGroup
	var i int64
	for i = 0; i < request.Number; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			resp := pb.ProtoResponse{Result: count}
			err := loopServer.Send(&resp)
			utils.CheckErr(err)
			log.Printf("count value: %d", count)
		}(i)
	}

	wg.Wait()
	return nil
}

func (s server) SumAndReturn(returnServer pb.ProtoService_SumAndReturnServer) error {
	total := make(chan int64, 1)
	done := make(chan bool)
	total <- 0
	go func() {
		for {
			request, err := returnServer.Recv()
			fmt.Printf("SumAndReturn function is invoked with %v \n", request)
			if err == io.EOF {
				done <- true
				return
			}
			utils.CheckErr(err)
			tmp := <-total
			tmp += request.Number
			total <- tmp
		}
	}()

	<-done

	response := &pb.ProtoResponse{Result: <-total}
	err := returnServer.Send(response)
	utils.CheckErr(err)

	return nil
}

func main() {
	fmt.Println("starting gRPC server...")

	listener, err := net.Listen("tcp", "localhost:50051")
	utils.CheckErr(err)

	grpcServer := grpc.NewServer()
	pb.RegisterProtoServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	utils.CheckErr(err)
}
