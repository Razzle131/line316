package core

import "errors"

var (
	ErrGripperAlreadyMoving = errors.New("gripper is moving already")
	ErrGripperWantToStop    = errors.New("gripper is want to stop")
)

var (
	ErrSlotOccupied = errors.New("this slot is busy")
	ErrSlotEmpty    = errors.New("this slot is empty")
)

var (
	ErrPuckPackaged = errors.New("puck is packaged")
)
