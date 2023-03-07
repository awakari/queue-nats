package main

import (
	"fmt"
	"github.com/awakari/queue-nats/api/grpc"
	"github.com/awakari/queue-nats/config"
	"github.com/awakari/queue-nats/service"
	"github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
	"os"
)

func main() {
	slog.Info("starting...")
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		panic(fmt.Sprintf("failed to load the config: %s", err))
	}
	opts := slog.HandlerOptions{
		Level: slog.Level(cfg.Log.Level),
	}
	log := slog.New(opts.NewTextHandler(os.Stdout))
	nc, err := nats.Connect(cfg.Nats.Uri, nats.UserInfo(cfg.Nats.User, cfg.Nats.Password))
	if err != nil {
		panic(fmt.Sprintf("failed to connect NATS: %s", err))
	}
	defer nc.Close()
	js, err := nc.JetStream()
	if err != nil {
		panic(fmt.Sprintf("failed to connect JetStream: %s", err))
	}
	svc := service.NewService(js)
	svc = service.NewLogging(svc, log)
	log.Info("connected, starting to listen for incoming requests...")
	if err = grpc.Serve(svc, cfg.Api.Port); err != nil {
		panic(err)
	}
}
