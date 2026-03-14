package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eztwokey/l3-serv/internal/api"
	"github.com/eztwokey/l3-serv/internal/config"
	"github.com/eztwokey/l3-serv/internal/logic"
	"github.com/eztwokey/l3-serv/internal/models"
	"github.com/eztwokey/l3-serv/internal/sender"
	"github.com/eztwokey/l3-serv/internal/storage"
	"github.com/eztwokey/l3-serv/internal/worker"
	"github.com/wb-go/wbf/logger"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

func main() {
	cfg := new(config.Config)
	if err := cfg.Read(config.LocalPath); err != nil {
		log.Fatal(err)
	}

	wbLog := logger.NewSlogAdapter("l3-serv", "local")

	redisClient := redis.New("", "", 0)

	store := storage.New(redisClient)

	rmqStrategy := retry.Strategy{
		Attempts: 5,
		Delay:    2 * time.Second,
		Backoff:  2,
	}

	rmqClient, err := rabbitmq.NewClient(rabbitmq.ClientConfig{
		URL:            cfg.RabbitMQ.URL,
		ConnectionName: "l3-serv",
		ConnectTimeout: 10 * time.Second,
		Heartbeat:      10 * time.Second,
		ReconnectStrat: rmqStrategy,
		ProducingStrat: rmqStrategy,
		ConsumingStrat: rmqStrategy,
	})

	if err != nil {
		log.Fatal("rabbitmq connect:", err)
	}

	wbLog.Info("rabbitmq connected", "healthy", rmqClient.Healthy())

	defer func() {
		if err := rmqClient.Close(); err != nil {
			wbLog.Error("rabbitmq close failed", "err", err)
		}
	}()

	if err := rmqClient.DeclareExchange(cfg.RabbitMQ.Exchange, "direct", true, false, false, nil); err != nil {
		log.Fatal("rabbitmq declare exchange:", err)
	}

	if err := rmqClient.DeclareQueue(
		cfg.RabbitMQ.Queue,
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.RoutingKey,
		true,
		false,
		true,
		nil,
	); err != nil {
		log.Fatal("rabbitmq declare queue:", err)
	}

	publisher := rabbitmq.NewPublisher(rmqClient, cfg.RabbitMQ.Exchange, "application/json")

	logic := logic.New(store, publisher, wbLog)

	senders := map[models.NotifyChannel]sender.Sender{
		models.ChannelTelegram: sender.NewTelegram(cfg.Telegram.BotToken),
		models.ChannelLog:      sender.NewLog(wbLog),
	}

	w := worker.New(store, senders, wbLog)

	consumer := rabbitmq.NewConsumer(rmqClient, rabbitmq.ConsumerConfig{
		Queue:         cfg.RabbitMQ.Queue,
		PrefetchCount: 10,
		Workers:       3,
		Nack: rabbitmq.NackConfig{
			Requeue: true,
		},
	}, w.Handle)

	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	go func() {
		if err := consumer.Start(consumerCtx); err != nil {
			wbLog.Error("consumer stopped", "err", err)
		}
	}()

	wbLog.Info("consumer started", "queue", cfg.RabbitMQ.Queue, "workers", 3)

	http := api.New(cfg, logic, wbLog)

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errChan <- http.Run()
	}()

	select {
	case signal := <-termChan:
		wbLog.Warn("got term signal chan", "signal", signal)
	case err := <-errChan:
		wbLog.Warn("got error chan", "error", err)
	}

	consumerCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	done := make(chan error, 1)
	go func() {
		done <- http.Shutdown(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		wbLog.Info("Server shutdown gracefully")
	case <-ctx.Done():
		wbLog.Info("Shutdown timeout exceeded, some connections were forcefully closed")
	}
}
