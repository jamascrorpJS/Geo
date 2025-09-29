package interactors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"jamascrorpJS/gwatch/internal/models"
	"jamascrorpJS/gwatch/pkg/cache"
	"jamascrorpJS/gwatch/pkg/redis"
)

type Coordinator interface {
	SendCoordinate(ctx context.Context, p models.Position) error
	ReceiveCoordinate(ctx context.Context, devId string) (models.Position, error)
}

type coordinate struct {
	cache   cache.InMemory
	storage redis.Storage
}

func New(cache cache.InMemory, storage redis.Storage) Coordinator {
	return &coordinate{cache: cache, storage: storage}
}

func (c *coordinate) SendCoordinate(ctx context.Context, p models.Position) error {
	if err := validatePosition(p); err != nil {
		return err
	}
	err := c.save(ctx, p)
	if err != nil {
		return err
	}
	return err
}

func (c *coordinate) ReceiveCoordinate(ctx context.Context, devId string) (models.Position, error) {
	position, ok := c.cache.Get(devId)
	if !ok {
		savedPosition, err := c.storage.SearchGeoPos(ctx, GeoKey, devId)
		if err != nil {
			fmt.Println(err.Error())
			if errors.Is(err, redis.ErrNotFound) {
				return models.Position{}, models.ErrNotFound
			}
			return models.Position{}, models.ErrInternal
		}
		c.cache.Set(savedPosition.DeviceID, savedPosition, cacheTTL)
		return ToPosition(savedPosition), nil
	}
	actualPosition, ok := position.(redis.GeoPos)
	if !ok {
		return models.Position{}, models.ErrInternal
	}
	return ToPosition(actualPosition), nil
}

func (c *coordinate) save(ctx context.Context, p models.Position) error {
	_, ok := c.cache.Get(p.DeviceID)
	if !ok {
		savedPos, err := c.storage.SearchGeoPos(ctx, GeoKey, p.DeviceID)
		if err != nil {
			if errors.Is(err, redis.ErrNotFound) {
				err := c.storage.AddGeoPos(ctx, GeoKey, ToGeoPosition(p))
				if err != nil {
					return models.ErrInternal
				}
				c.cache.Set(p.DeviceID, p, cacheTTL)
				return nil
			}
			return models.ErrInternal
		}
		c.cache.Set(p.DeviceID, savedPos, cacheTTL)
	}
	err := c.storage.AddGeoPos(ctx, GeoKey, ToGeoPosition(p))
	if err != nil {
		return models.ErrInternal
	}
	c.cache.Set(p.DeviceID, p, cacheTTL)
	return nil
}

func validatePosition(p models.Position) error {
	if p.DeviceID == "" {
		return models.ErrDeviceId
	}
	if p.Latitude < -90 || p.Latitude > 90 {
		return models.ErrPosition
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return models.ErrPosition
	}
	return nil
}

const GeoKey = "location"
const cacheTTL = 10 * time.Second
