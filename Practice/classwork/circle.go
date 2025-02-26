package main

import (
	"fmt"
	"math"
)

type Shape interface{
	Area() float64
	Perimeter() float64
}

type rectangle struct{
	height float64
	width float64
}

type circle struct{
	radius  float64
}

func (c *circle) Area() float64{
	
	return math.Pi * c.radius * c.radius
}
func (c *circle) Perimeter() float64{
	return 2 * math.Pi * c.radius
}

func (r *rectangle) Area() float64{
	return r.height * r.width
}
func (r *rectangle) Perimeter() float64{
	return 2*(r.height + r.width)
}

func main(){


	var s Shape
	s=circle{radius:5}
	fmt.Println("area:", s.Area())
	fmt.Println("perimeter", s.Perimeter())


	s=rectangle{height:10,width:20}
	fmt.Println("area:",s.Area())
	fmt.Println("perimeter:",s.Perimeter())

}