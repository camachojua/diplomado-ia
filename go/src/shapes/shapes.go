package main

import (
	"fmt"
	"math"
)

func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

type Circle struct {
	x, y, r float64
}

type Triangle struct {
	a, b, c float64
}

type Rectangle struct {
	x1, y1, x2, y2 float64
}

type Shape interface {
	area() float64
}

func (o *Circle) area() float64 {
	return math.Pi * o.r * o.r
}

// Usa la fórmula de Heron para calcular el área
func (o *Triangle) area() float64 {
	p := (o.a + o.b + o.c) / 2
	return math.Sqrt(p * (p - o.a) * (p - o.b) * (p - o.c))
}

func (o *Rectangle) area() float64 {
	l := distance(o.x1, o.y1, o.x1, o.y2)
	w := distance(o.x1, o.y1, o.x2, o.y1)
	return l * w
}

func totalArea(shapes ...Shape) float64 {
	var area float64
	for _, s := range shapes {
		area += s.area()
	}

	return area
}

func main() {
	c := Circle{0, 0, 5}
	r := Rectangle{0, 0, 5, 5}
	t := Triangle{5, 10, 15}

	fmt.Println("Total area is ", totalArea(&c, &r, &t))

}
