package service

import (
	"context"
	"fmt"
	"github.com/maei/golang_grpc_bidirect_streaming_maximum/grpc_client/src/client"
	"github.com/maei/golang_grpc_bidirect_streaming_maximum/grpc_client/src/domain/maximumpb"
	"github.com/maei/shared_utils_go/logger"
	"io"
	"log"
	"math/rand"
)

var MaximumService maximumServiceInterface = &maximumService{}

type maximumServiceInterface interface {
	GetMaximum()
}

type maximumService struct{}

func generateNumbers() int32 {
	min := 1
	max := 100

	n := min + rand.Intn(max-min+1)
	return int32(n)
}

func (*maximumService) GetMaximum() {
	numReq := int8(5)

	conn, clientErr := client.GRPCClient.SetClient()
	if clientErr != nil {
		logger.Error("gRPC-Client: Error while creating gRPC-Client", clientErr)
	}

	c := maximumpb.NewMaximumServiceClient(conn)

	stream, streamErr := c.GetMaximum(context.Background())
	if streamErr != nil {
		logger.Error("gRPC-Client: Error while creating stream-obj", streamErr)
	}

	// channel for closing routines
	done := make(chan bool)

	// go-routine for sending requests to gRPC-Server

	go func(n int8) {
		for i := int8(0); i <= n; i++ {
			number := generateNumbers()
			log.Println(fmt.Sprintf("generated number: %v", number))
			req := &maximumpb.MaximumRequest{Number: number}
			err := stream.Send(req)
			if err != nil {
				logger.Error("gRPC-Client: Error while sending Data to gRPC-Server", err)
			}
			//time.Sleep(2 * time.Second)
		}
		stream.CloseSend()

	}(numReq)

	// go-routine for receiving responses from gRPC-Server
	go func() {
		for {
			res, resErr := stream.Recv()
			if resErr == io.EOF {
				break
			}
			if resErr != nil {
				logger.Error("gRPC-Client: Error while receiving Data from gRPC-Server", resErr)
			}
			log.Println(fmt.Sprintf("Result from GetMaximum: %v", res.GetResult()))
		}
		done <- true
	}()

	<-done
	logger.Info("gRPC-Client: Finished GetMaximum request")
}
