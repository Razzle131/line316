package core

import "time"

const (
	tickrate = 100 // number of ticks in second
)

const (
	gripperBaseLength        = 0.065 // m
	gripperLeftSensorLength  = 0.02  // m
	gripperRightSensorLength = 0.02  // m
	gripperRailLength        = 0.64  // m

	gripperVerticalSpeed = 0.1  // m/s
	gripperUpPos         = 0.09 // m
	gripperDownPos       = 0    // m

	gripperHorizontalSpeed = 0.1                                                                // m/s
	gripperCarouselPos     = gripperLeftSensorLength + gripperBaseLength/2                      // m
	gripperStartPos        = 0.2                                                                // m
	gripperPackagingPos    = 0.4                                                                // m
	gripperSortingPos      = gripperRailLength - gripperRightSensorLength - gripperBaseLength/2 // m
	gripperAbleMiss        = 0.02                                                               // m, +- from where gripper can still operate normal, like in normal position
)

const (
	carouselTotalSlots   = 6
	carouselInspectSlot  = 4 // numeration from zero in carousel gripper pos
	carouselDrillSlot    = 5 // numeration from zero in carousel gripper pos
	carouselNextSlotTime = time.Millisecond * 200
)

const (
	packagingTime = time.Millisecond * 1000
)

const (
	sortingTime = time.Millisecond * 1000
)
