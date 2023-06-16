package main

import (
	"context"
	"fmt"
	stlog "log"
	"simple_distributed/grades"
	"simple_distributed/log"
	"simple_distributed/registry"
	"simple_distributed/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService}, // 依赖的服务
		ServiceUpdateURL: serviceAddress + "/services",
	}
	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	// 获取logService的依赖，然后对依赖进行遍历
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("logging service found at:%s\n", logProvider)
		// 把url和服务名称传递过去
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	// 等待信号
	<-ctx.Done()

	fmt.Println("shutting down grading service ")

}
