package sensor

type Sensor struct {
	name  string
	addr  string
	value bool
}

func New(name, addr string) *Sensor {
	return &Sensor{
		name:  name,
		addr:  addr,
		value: false,
	}
}

func (s *Sensor) GetValue() bool {
	return s.value
}

func (s *Sensor) WriteValue(newValue bool) {
	s.value = newValue
}
