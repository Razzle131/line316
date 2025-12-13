package core

type Actuator interface {
	IsActivated() bool
	Activate()
	Deactivate()
}

type Sensor interface {
	GetValue() bool
	WriteValue(newValue bool)
}

type Pucker interface {
	TakePuck() (Puck, error)
	PlacePuck(puck Puck) error
}
