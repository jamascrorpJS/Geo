package interactors

import (
	"jamascrorpJS/gwatch/internal/models"
	"jamascrorpJS/gwatch/pkg/redis"
)

func ToGeoPosition(geo models.Position) redis.GeoPos {
	return redis.GeoPos{
		DeviceID:  geo.DeviceID,
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
		Timestamp: geo.Timestamp,
	}
}

func ToPosition(geo redis.GeoPos) models.Position {
	return models.Position{
		DeviceID:  geo.DeviceID,
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
		Timestamp: geo.Timestamp,
	}
}
