package main

import (
	"errors"
	"log"
	"mcontext/internal/engine"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	e, err := engine.New()
	defer func() {
		if e != nil {
			e.Exit()
		}
	}()

	if err != nil {
		log.Printf("Engine new error: %v\n", err)
		return
	}

	err = e.Prepare()
	if err != nil {
		log.Printf("Engine prepare error: %v\n", err)
		return
	}

	// 开启监听协程
	go func() {
		if err := e.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Engine listen error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	// 接收到信号
}
