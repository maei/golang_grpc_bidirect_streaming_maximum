package main

import (
	"github.com/maei/golang_grpc_bidirect_streaming_maximum/grpc_server/src/server"
	"github.com/maei/shared_utils_go/logger"
)

func main() {
	logger.Info("Starting gRPC-Server for BiDi streaming for finding Maximum in stream of values")
	server.StartGRPCServer()
}
