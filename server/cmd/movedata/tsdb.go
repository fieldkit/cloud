package main

import "context"

type MoveDataToTimeScaleDBHandler struct {
}

func NewMoveDataToTimeScaleDBHandler() *MoveDataToTimeScaleDBHandler {
	return &MoveDataToTimeScaleDBHandler{}
}

func (h *MoveDataToTimeScaleDBHandler) MoveReadings(ctx context.Context, readings []*MovedReading) error {
	return nil
}

func (h *MoveDataToTimeScaleDBHandler) Close(ctx context.Context) error {
	return nil
}
