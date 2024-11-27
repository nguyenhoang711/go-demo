package embedded_struct

import (
	"fmt"
	"io"
	"math"
)

type Foo struct {
	io.ReadCloser // embbeded struct
	SomeMoreStuff int
}

type Person struct {
	Name string
}
func (p *Person) Talk() {
	fmt.Println("Hi, my name is", p.Name)
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Developer struct {
	Person
	Model string
}

type Circle struct {
	R float64
}

type Rectangle struct {
	X1, Y1, Z1 float64
}

func (r *Rectangle) Area() float64 {
	return r.X1 + r.Y1 + r.Z1
}

func (r *Rectangle) Perimeter() float64 {
	return r.X1 * r.Z1 * 1 / 2
}

func (r *Circle) Area() float64 {
	return 2 * math.Pi * r.R
}

func (r *Circle) Perimeter() float64 {
	return math.Pi * r.R * r.R
}

func TotalArea(shapes ...Shape) float64 {
	var area float64
	for _, s := range shapes {
	  area += s.Area()
	}
	return area
}

type MultiShape struct {
	Shapes []Shape
}

func (m *MultiShape) Area() float64 {
	var area float64
	for _, s := range m.Shapes {
	  area += s.Area()
	}
	return area
}

func (m *MultiShape) Perimeter() float64 {
	var perimeter float64
	for _, s := range m.Shapes {
	  perimeter += s.Perimeter()
	}
	return perimeter
}