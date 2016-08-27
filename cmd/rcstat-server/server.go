package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/ipc"
)

type (
	ServerConfig struct {
		Base   BaseConfig
		LogCfg LogConfig
	}

	BaseConfig struct {
		C        string
		LogFile  string `default:"./rcstat_server.log"`
		LogLevel string `default:"info"`
		IpcAddr  string `default:"ipc:///tmp/rcstat_server.ipc"`
		Plugins  []string
	}
)

// CollectServer defines a server collecting redis command stats
type CollectServer struct {
	Config *ServerConfig
	Logger *logrus.Logger
	Sock   mangos.Socket
}

// NewCollectServer returns a new CollectServer instance.
func NewCollectServer(cfg *ServerConfig) *CollectServer {
	server := &CollectServer{Config: cfg}
	logger := logrus.New()
	logger.Out = cfg.LogCfg.Output
	logger.Level = cfg.LogCfg.Level
	logger.Formatter = cfg.LogCfg.Format
	server.Logger = logger
	return server
}

// Run CollectServer
func (server *CollectServer) Run(exit chan struct{}) {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = pull.NewSocket(); err != nil {
		server.Logger.Fatalf("Create pull socket error: %s", err)
		close(exit)
	}
	sock.AddTransport(ipc.NewTransport())
	if err = sock.Listen(server.Config.Base.IpcAddr); err != nil {
		server.Logger.Fatalf("Listen on pull socket error: %s", err.Error())
		close(exit)
	}
	server.Sock = sock
	go func() {
		server.Logger.Info("rcstat connect server starts...")
		for {
			msg, err = server.Sock.Recv()
			server.Logger.Debugf("receive: %s\n", string(msg))
		}
	}()
}

func (server *CollectServer) Shutdown() {
	server.Sock.Close()
}
