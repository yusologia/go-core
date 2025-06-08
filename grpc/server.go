package logiagrpc

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

type GRPCServer struct {
	listener net.Listener
	Server   *grpc.Server
}

type GRPCServerRegister interface {
	Register(*grpc.Server)
}

func (srv *GRPCServer) NewServer(address string) *GRPCServer {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	srv.listener = lis
	srv.Server = s

	return srv
}

func (srv *GRPCServer) Register(interfaces ...GRPCServerRegister) *GRPCServer {
	for _, itf := range interfaces {
		itf.Register(srv.Server)
	}

	return srv
}

func (srv *GRPCServer) Serve() error {
	return srv.Server.Serve(srv.listener)
}
