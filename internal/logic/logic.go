package logic

import (
	"github.com/wb-go/wbf/logger"
	"github.com/wb-go/wbf/rabbitmq"

	"github.com/eztwokey/l3-serv/internal/storage"
)

type Logic struct {
	store     *storage.Storage
	publisher *rabbitmq.Publisher
	logger    logger.Logger
}

func New(store *storage.Storage, publisher *rabbitmq.Publisher, logger logger.Logger) *Logic {
	return &Logic{store: store, publisher: publisher, logger: logger}
}
