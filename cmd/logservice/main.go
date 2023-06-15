package main

import (
	"context"
	"fmt"
	stlog "log"
	"simple_distributed/log"
	"simple_distributed/registry"
	"simple_distributed/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	r := registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}

	// 等待信号
	<-ctx.Done()

	fmt.Println("shutting down log service ")

}
