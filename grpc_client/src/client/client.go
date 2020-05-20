package client

import (
	"google.golang.org/grpc"
	"log"
)

var GRPCClient grpcClientInterface = &grpcClient{}

type grpcClientInterface interface {
	SetClient() (*grpc.ClientConn, error)
}

type grpcClient struct{}

func (*grpcClient) SetClient() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return nil, err
	}

	return conn, nil
}
