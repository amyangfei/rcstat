package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/koding/multiconfig"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var Config *ServerConfig

func initConfig(configFile string) error {
	m := multiconfig.NewWithPath(configFile)
	base := &BaseConfig{}
	m.MustLoad(base)
	Config = &ServerConfig{}
	Config.Base = *base

	f, err := os.OpenFile(base.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	Config.LogCfg.Output = f
	Config.LogCfg.Format = &logrus.JSONFormatter{}
	Config.LogCfg.Level = LogString2Level(base.LogLevel)

	return nil
}

func singalHandle(s os.Signal, server *CollectServer) {
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
		fmt.Printf("recevie signal %v, exit.", s)
		server.Shutdown()
		os.Exit(0)
	case syscall.SIGHUP:
		// TODO: reload
		fmt.Printf("server reload...")
	}
}

func main() {
	var configFile string
	var printVersion bool

	flag.BoolVar(&printVersion, "version", false, "print version")
	flag.StringVar(&configFile, "c", "config.toml", "path to config file")
	flag.Parse()

	if printVersion {
		PrintVersion()
		os.Exit(0)
	}

	if err := initConfig(configFile); err != nil {
		panic(err)
	}

	exit := make(chan struct{})
	server := NewCollectServer(Config)
	server.Run(exit)

	stop := make(chan os.Signal)
	signal.Notify(
		stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var wg sync.WaitGroup

	select {
	case s := <-stop:
		singalHandle(s, server)
	case <-exit:
		fmt.Println("collect server exit.")
		os.Exit(1)
	}
	wg.Wait()
}
