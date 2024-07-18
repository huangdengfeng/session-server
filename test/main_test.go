package test

import (
	"github.com/alicebob/miniredis/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"session-server/entity/config"
	gclient "session-server/entity/grpc/client"
	"session-server/entity/grpc/server"
	"session-server/entity/pb"
	"session-server/logic"
	"session-server/repo/cache"
	"session-server/service"
	"testing"
)

var client pb.SessionClient

func TestMain(m *testing.M) {
	setup()
	// 运行测试
	exitCode := m.Run()
	// 退出测试
	teardown()
	os.Exit(exitCode)
}
func setup() {
	config.ServerConfigPath = "../conf"
	mr := miniredis.NewMiniRedis()
	err := mr.StartAddr(":6379")
	if err != nil {
		panic(err)
	}
	config.Init()

	var createSessionServer = func() pb.SessionServer {
		dao := cache.NewRedisDao(config.RedisClient)
		sessionService := logic.NewSessionService(dao)
		return service.NewSessionServer(sessionService)
	}

	server.Start(createSessionServer())

	// start client
	client = gclient.CreateClient(config.Global.Server.Listen)

	log.Infof("[test] set up init success")
}

func teardown() {
	config.Shutdown()
	log.Infof("[test] tear down success")
}
