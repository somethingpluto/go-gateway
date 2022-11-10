package main

import (
	"go_gateway/common/lib"
	"go_gateway/initialize"
	"go_gateway/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = lib.InitModule("./conf/dev/")
	defer lib.Destroy()
	initialize.Init()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
