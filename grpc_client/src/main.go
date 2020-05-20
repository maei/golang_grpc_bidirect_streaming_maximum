package main

import (
	"github.com/maei/golang_grpc_bidirect_streaming_maximum/grpc_client/src/service"
	"github.com/maei/shared_utils_go/logger"
	"math/rand"
	"time"
)

func main() {
	logger.Info("gRPC-Client: Starting gRPC-Client")
	rand.Seed(time.Now().UnixNano())
	service.MaximumService.GetMaximum()
}
