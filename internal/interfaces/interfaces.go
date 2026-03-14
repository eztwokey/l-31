package interfaces

import (
	"context"

	"github.com/eztwokey/l3-serv/internal/models"
)

type Storage interface {
	Create(ctx context.Context, n models.Notification) (models.Notification, error)
	Get(ctx context.Context, id string) (models.Notification, error)
	Update(ctx context.Context, n models.Notification) (models.Notification, error)
	Delete(ctx context.Context, id string) error
}
