package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MAX_GRPC_SIZE = 1024 * 1024 * 100

func NewGrpcConnection(grpcUrl string) (*grpc.ClientConn, error) {
	return grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MAX_GRPC_SIZE),
			grpc.MaxCallSendMsgSize(MAX_GRPC_SIZE)), grpc.WithBlock())
}
