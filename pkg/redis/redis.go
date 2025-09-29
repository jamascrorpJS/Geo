package redis

import (
	"context"
	"errors"
	"fmt"
	"log"

	redis2 "github.com/redis/go-redis/v9"

	"jamascrorpJS/gwatch/pkg/config"
)

type Storage interface {
	AddGeoPos(ctx context.Context, key string, geoPos GeoPos) error
	SearchGeoPos(ctx context.Context, key string, member string) (GeoPos, error)
}

type redis struct {
	r      *redis2.Client
	config config.Config
}

func New(ctx context.Context, config config.Config) Storage {

	var (
		url  = config.GetString("redis.url")
		ping *redis2.StatusCmd
	)

	defer func() {
		err := ping.Err()
		if err == nil {
			log.Println("Success redis ping")
			return
		}

		log.Fatal(err)
	}()

	var (
		port     = config.GetString("redis.port")
		password = config.GetString("redis.password")
		db       = config.GetInt("redis.db")
	)

	client := redis2.NewClient(&redis2.Options{
		Addr:     url + ":" + port,
		Password: password,
		DB:       db,
	})

	ping = client.Ping(ctx)

	return &redis{
		r:      client,
		config: config,
	}
}

type GeoPos struct {
	DeviceID  string
	Latitude  float64
	Longitude float64
	Timestamp int
}

func (r *redis) AddGeoPos(ctx context.Context, key string, geoPos GeoPos) error {

	_, err := r.r.GeoAdd(ctx, key, &redis2.GeoLocation{
		Name:      geoPos.DeviceID,
		Longitude: geoPos.Longitude,
		Latitude:  geoPos.Latitude,
	}).Result()

	if err != nil {
		return fmt.Errorf("Errors from GeoAdd redis %w", err)
	}

	_, err = r.r.ZAdd(ctx, geoPos.DeviceID, redis2.Z{
		Score:  float64(geoPos.Timestamp),
		Member: fmt.Sprintf("%v,%v", geoPos.Latitude, geoPos.Longitude),
	}).Result()

	if err != nil {
		return fmt.Errorf("Errors from ZAdd redis %w", err)
	}

	return nil
}

var ErrNotFound = errors.New("GeoPos not found")

func (r *redis) SearchGeoPos(ctx context.Context, key string, member string) (GeoPos, error) {
	var geoPos GeoPos
	geo, err := r.r.GeoPos(ctx, key, member).Result()
	if err != nil {
		return geoPos, fmt.Errorf("Error from GeoPos%w", err)
	}
	ss, err := r.r.ZRevRangeWithScores(ctx, member, 0, 0).Result()
	if err != nil {
		return geoPos, fmt.Errorf("Error from ZRevRangeWithScores %w", err)
	}
	if len(geo) == 0 || geo[0] == nil || len(ss) == 0 {
		return geoPos, fmt.Errorf("Error %w", ErrNotFound)
	}
	for _, v := range geo {
		geoPos.Longitude = v.Longitude
		geoPos.Latitude = v.Latitude
	}
	for _, v := range ss {
		geoPos.Timestamp = int(v.Score)
	}
	geoPos.DeviceID = member
	return geoPos, nil
}
