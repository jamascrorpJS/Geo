package models

import "errors"

type Response struct {
	Code int
	Mes  string
	Data interface{}
}

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal err")
	ErrDeviceId = errors.New("no correct id")
	ErrPosition = errors.New("no correct position")
)
