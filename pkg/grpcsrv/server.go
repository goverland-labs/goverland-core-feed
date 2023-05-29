package grpcsrv

import (
	"google.golang.org/grpc"
)

func NewGrpcServer() *grpc.Server {
	server := grpc.NewServer(
		StdUnaryMiddleware(),
		StdStreamMiddleware(),
	)

	StdRegister(server)

	return server
}
