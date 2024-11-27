package factory_pattern

type Car interface {
	Drive() string
}

type CarFactory struct {}

type Sendar struct {}

func (s *Sendar) Drive() string {
	return "I am driving Sendar car"
}

type SUV struct {}

func (s *SUV) Drive() string {
	return "I am driving SUV car"
}

