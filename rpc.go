package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Service struct {
	server *rpc.Server
}

func NewService(storage *Storage) (*Service, error) {
	server := rpc.NewServer()
	err := server.RegisterName("user", storage)
	if err != nil {
		return nil, err
	}
	return &Service{
		server,
	}, nil
}

func (service *Service) Start(addr string) error {
	service.server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go service.server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
