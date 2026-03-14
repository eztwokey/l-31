package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eztwokey/l3-serv/internal/models"
	"github.com/wb-go/wbf/redis"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Storage {
	return &Storage{rdb: rdb}
}

func notifyKey(id string) string { return "notify:" + id }

func (s *Storage) Create(ctx context.Context, n models.Notification) (models.Notification, error) {
	b, err := json.Marshal(n)
	if err != nil {
		return models.Notification{}, fmt.Errorf("marshal: %w", err)
	}
	if err := s.rdb.Set(ctx, notifyKey(n.ID), string(b)); err != nil {
		return models.Notification{}, err
	}
	return n, nil
}

func (s *Storage) Get(ctx context.Context, id string) (models.Notification, error) {
	g, err := s.rdb.Get(ctx, notifyKey(id))
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return models.Notification{}, ErrNotFound
		}
		return models.Notification{}, err
	}

	var n models.Notification
	if err := json.Unmarshal([]byte(g), &n); err != nil {
		return models.Notification{}, fmt.Errorf("unmarshal: %w", err)
	}
	return n, nil
}

func (s *Storage) Update(ctx context.Context, n models.Notification) (models.Notification, error) {
	if _, err := s.Get(ctx, n.ID); err != nil {
		return models.Notification{}, err
	}

	return s.Create(ctx, n)
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	if err := s.rdb.Del(ctx, notifyKey(id)); err != nil {
		return err
	}

	return nil
}
