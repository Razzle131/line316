package actuator

type Actuator struct {
	name        string
	addr        string
	isActivated bool
}

func New(name, addr string) *Actuator {
	return &Actuator{
		name:        name,
		addr:        addr,
		isActivated: false,
	}
}

func (s *Actuator) IsActivated() bool {
	return s.isActivated
}

func (s *Actuator) Activate() {
	s.isActivated = true
}

func (s *Actuator) Deactivate() {
	s.isActivated = false
}
