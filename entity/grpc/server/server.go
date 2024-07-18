package server

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"session-server/entity/config"
	"session-server/entity/pb"
)

func Start(s pb.SessionServer) *grpc.Server {

	var grpcServer = grpc.NewServer(CreateDefaultInterceptor())
	pb.RegisterSessionServer(grpcServer, s)

	go func() {
		listen, err := net.Listen("tcp", config.Global.Server.Listen)
		if err != nil {
			log.Fatalf("listen error [%s]", err)
		}
		err = grpcServer.Serve(listen)
		if err != nil {
			log.Fatalf("server serve error [%s]", err)
		}
	}()
	return grpcServer
}

func Stop(grpcServer *grpc.Server) {
	grpcServer.GracefulStop()
}
