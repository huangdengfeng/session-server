package test

import (
	"github.com/alicebob/miniredis/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"os"
	"session-server/entity/config"
	c "session-server/entity/grpc/client"
	"session-server/entity/grpc/server"
	"session-server/entity/pb"
	"session-server/logic"
	"session-server/repo/cache"
	"session-server/service"
	"testing"
	"time"
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
	kacp := keepalive.ClientParameters{
		Time:                1 * time.Minute,  // 客户端每隔1min发送一次心跳ping
		Timeout:             10 * time.Second, // 如果没有收到服务端的心跳响应，认为连接失败的超时时间
		PermitWithoutStream: true,             // 即使没有活动的RPC流，也允许发送心跳
	}
	opts := []grpc.DialOption{grpc.WithKeepaliveParams(kacp), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 1 * time.Second})}
	opts = append(opts, c.CreateDefaultInterceptor())
	conn, err := grpc.NewClient(config.Global.Server.Listen, opts...)

	if err != nil {
		log.Fatalf("connect error [%s]", err)
	}
	client = pb.NewSessionClient(conn)

	log.Infof("[test] set up init success")
}

func teardown() {
	config.Shutdown()
	log.Infof("[test] tear down success")
}
