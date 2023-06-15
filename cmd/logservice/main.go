package main

import (
	"context"
	"fmt"
	stlog "log"
	"simple_distributed/log"
	"simple_distributed/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(
		context.Background(),
		"log Service",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}

	// 等待信号
	<-ctx.Done()

	fmt.Println("shutting down log service ")

}
