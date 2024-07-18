package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"session-server/entity/config"
	"session-server/entity/grpc/server"
	"session-server/entity/pb"
	"session-server/logic"
	"session-server/repo/cache"
	"session-server/service"
	"syscall"
)

func main() {
	config.Init()
	defer config.Shutdown()

	s := server.Start(createSessionServer())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	o := <-sig
	log.Printf("recieve signal %s ,server will stop gracefully", o.String())
	server.Stop(s)
}

func createSessionServer() pb.SessionServer {
	dao := cache.NewRedisDao(config.RedisClient)
	sessionService := logic.NewSessionService(dao)
	return service.NewSessionServer(sessionService)
}
