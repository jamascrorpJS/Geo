package models

type Position struct {
	DeviceID  string  `json:"deviceID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int     `json:"timestamp"`
	Speed     float32 `json:"-"`
}
