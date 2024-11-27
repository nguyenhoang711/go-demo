package pointertutor

type Person struct {
    Name string
    Age  int
}

func NewPerson(name string, age int) *Person {
    p := new(Person)
    p.Name = name
    p.Age = age
    return p
}

func (p *Person) IncrementAge() {
    p.Age++
}