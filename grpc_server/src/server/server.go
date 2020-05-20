package server

import (
	"fmt"
	"github.com/maei/golang_grpc_bidirect_streaming_maximum/grpc_server/src/domain/maximumpb"
	"github.com/maei/shared_utils_go/logger"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type server struct{}

var (
	s = grpc.NewServer()
)

func (*server) GetMaximum(stream maximumpb.MaximumService_GetMaximumServer) error {
	//Result var
	var maxNumber int32

	// channel for receiving jobs
	jobs := make(chan int32)
	// channel for closing jobs to done
	done := make(chan bool)

	// receive routine for gRPC-Data from Client
	go func() {
		for {
			req, reqErr := stream.Recv()
			if reqErr == io.EOF {
				break
			}
			if reqErr != nil {
				logger.Error("gRPC-Server: Error while receiving Data from gRPC-Client", reqErr)
				break
			}
			jobs <- req.GetNumber()
		}
		close(jobs)
	}()
	// sender routine to send results to Client
	go func(number int32) {
		for {
			j, status := <-jobs
			if status {
				if number < j {
					number = j
				}
				res := &maximumpb.MaximumResponse{Result: number}
				sendErr := stream.Send(res)
				if sendErr != nil {
					logger.Error("gRPC-Server: Error while sending Data to gRPC-Client", sendErr)
					break
				}
				time.Sleep(2 * time.Second)

			} else {
				logger.Info("gRPC-Server: Received all jobs, closing channel")
				done <- true
				return
			}
		}
	}(maxNumber)

	<-done
	logger.Info(fmt.Sprintf("gRPC-Server: Finishing gRPC-Client Request with result: %v", maxNumber))
	return nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("error while listening gRPC Server", err)
	}

	maximumpb.RegisterMaximumServiceServer(s, &server{})

	errServer := s.Serve(lis)
	if errServer != nil {
		logger.Error("error while serve gRPC Server", errServer)
	}
}
