package main

import (
	"context"
	"fmt"
	stlog "log"
	"simple_distributed/grades"
	"simple_distributed/registry"
	"simple_distributed/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	r := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}

	// 等待信号
	<-ctx.Done()

	fmt.Println("shutting down grading service ")

}
