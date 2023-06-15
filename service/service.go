package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start ctx 上下文 registerHandlersFunc要注册的方法，返回后注册完之后的上下文
func Start(ctx context.Context,
	serviceName, host, port string, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}
func startService(ctx context.Context,
	serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		// if 启动的时候出现了错误，就将打印
		log.Println(srv.ListenAndServe())
		// 然后执行取消
		cancel()
	}()

	go func() {
		// 用户可以自己停止服务
		fmt.Printf("%v started. press any key to stop. \n", serviceName)
		var s string
		_, _ = fmt.Scanln(&s)
		_ = srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
