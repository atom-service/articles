package main

import (
	"flag"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/yinxulai/goutils/grpc/interceptor"
	"github.com/yinxulai/grpc-module-articles/provider"
	"github.com/yinxulai/grpc-module-articles/standard"
	"google.golang.org/grpc"
)

var isDev bool

func init() {
	godotenv.Load()
	flag.BoolVar(&isDev, "dev", false, "运行模式，可选 dev")
}

func main() {
	flag.Parse() // 解析获取参数

	rpcListenAddress := os.Getenv("PRC_LISTEN_ADDRESS")
	httpDevListenAddress := os.Getenv("HTTP_DEV_LISTEN_ADDRESS")

	if isDev { // 开发模式 启动一个调试工具
		go standard.Serve(httpDevListenAddress, rpcListenAddress, standard.DefaultHtmlStringer, grpc.WithInsecure())
	}

	lis, err := net.Listen("tcp", rpcListenAddress)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(interceptor.NewCalllogs()...)
	standard.RegisterArticlesServer(grpcServer, provider.NewService())
	panic(grpcServer.Serve(lis))
}
