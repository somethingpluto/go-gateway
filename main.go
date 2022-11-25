package main

import (
	"flag"
	"go_gateway/common/lib"
	"go_gateway/initialize"
	"go_gateway/router"
	"os"
	"os/signal"
	"syscall"
)

var (
	endpoint = flag.String("endpoint", "dashboard", "input endpoint dashboard or server")
	config   = flag.String("config", "./conf/dev/", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()

	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "dashboard" {
		err := lib.InitModule("./conf/dev/")
		if err != nil {
			panic(err)
		}
		defer lib.Destroy()
		initialize.Init()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	}

}
