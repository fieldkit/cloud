package main

import (
	"time"
)

type Location struct {
	UpdatedAt   time.Time
	Coordinates []float32
}

type StreamState struct {
	Location Location
}
