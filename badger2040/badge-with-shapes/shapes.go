package main

import (
	"image/color"
	"math"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinydraw"
)

type ShapeType int

// ShapeType constants
const (
	CIRCLE ShapeType = iota
	RECTANGLE
	TRIANGLE
)

type ShapePlotter interface {
	Shape
	Hitbox(Shape) Hitbox
	Plot(drivers.Displayer, Shape)
}

type Shape struct {
	Type    ShapeType
	X       int16
	Y       int16
	Outline color.RGBA
	Fill    color.RGBA
	Shape   interface{}
	hitbox  Hitbox
}

type Hitbox struct {
	RegistrationX int16
	RegistrationY int16
	MinX          int16
	MinY          int16
	MaxX          int16
	MaxY          int16
	SizeX         int16
	SizeY         int16
}

func (shape *Shape) Hitbox() Hitbox {
	hitbox := Hitbox{}
	switch shape.Type {
	case CIRCLE:
		hitbox = shape.Shape.(Circle).Hitbox(shape)
		break
	case RECTANGLE:
		hitbox = shape.Shape.(Rectangle).Hitbox(shape)
		break
	case TRIANGLE:
		hitbox = shape.Shape.(Triangle).Hitbox(shape)
		break
	}
	return hitbox
}

func (shape *Shape) Plot(display drivers.Displayer) {
	switch shape.Type {
	case CIRCLE:
		shape.Shape.(Circle).Plot(display, shape)
		break
	case RECTANGLE:
		shape.Shape.(Rectangle).Plot(display, shape)
		break
	case TRIANGLE:
		shape.Shape.(Triangle).Plot(display, shape)
		break
	}
}

type Circle struct {
	Radius int16
}

func NewCircle(x, y, size int16, outline, fill color.RGBA) Shape {
	// treating `size` as diameter, not radius
	const minimumRadius = 5
	radius := int16(math.Round(math.Max(minimumRadius, float64(size)/2)))

	return Shape{
		Type:    CIRCLE,
		X:       x,
		Y:       y,
		Outline: outline,
		Fill:    fill,
		Shape: Circle{
			Radius: radius,
		},
	}
}

func (circle Circle) Hitbox(shape *Shape) Hitbox {
	shape.hitbox.RegistrationX = circle.Radius
	shape.hitbox.RegistrationY = circle.Radius
	shape.hitbox.MinX = shape.X - shape.hitbox.RegistrationX
	shape.hitbox.MaxX = shape.X + shape.hitbox.RegistrationX
	shape.hitbox.MinY = shape.Y - shape.hitbox.RegistrationY
	shape.hitbox.MaxY = shape.Y + shape.hitbox.RegistrationY
	shape.hitbox.SizeX = shape.hitbox.MaxX - shape.hitbox.MinX
	shape.hitbox.SizeY = shape.hitbox.MaxY - shape.hitbox.MinY
	return shape.hitbox
}

func (circle Circle) Plot(display drivers.Displayer, shape *Shape) {
	circle.Hitbox(shape)

	tinydraw.FilledCircle(
		display,
		shape.X,
		shape.Y,
		circle.Radius,
		shape.Fill,
	)
	tinydraw.Circle(
		display,
		shape.X,
		shape.Y,
		circle.Radius,
		shape.Outline,
	)
}

type Rectangle struct {
	Width  int16
	Height int16
}

func NewRectangle(x, y, width, height int16, outline, fill color.RGBA) Shape {
	return Shape{
		Type:    RECTANGLE,
		X:       x,
		Y:       y,
		Outline: outline,
		Fill:    fill,
		Shape: Rectangle{
			Width:  width,
			Height: height,
		},
	}
}

func (rectangle Rectangle) Hitbox(shape *Shape) Hitbox {
	shape.hitbox.RegistrationX = 0
	shape.hitbox.RegistrationY = 0
	shape.hitbox.MinX = shape.X
	shape.hitbox.MaxX = shape.X + rectangle.Width
	shape.hitbox.MinY = shape.Y
	shape.hitbox.MaxY = shape.Y + rectangle.Height
	shape.hitbox.SizeX = shape.hitbox.MaxX - shape.hitbox.MinX
	shape.hitbox.SizeY = shape.hitbox.MaxY - shape.hitbox.MinY
	return shape.hitbox
}

func (rectangle Rectangle) Plot(display drivers.Displayer, shape *Shape) {
	rectangle.Hitbox(shape)

	tinydraw.FilledRectangle(
		display,
		shape.X,
		shape.Y,
		rectangle.Width,
		rectangle.Height,
		shape.Fill,
	)
	tinydraw.Rectangle(
		display,
		shape.X,
		shape.Y,
		rectangle.Width,
		rectangle.Height,
		shape.Outline,
	)
}

type Triangle struct {
	Size int16
}

func NewTriangle(x, y, size int16, outline, fill color.RGBA) Shape {
	return Shape{
		Type:    TRIANGLE,
		X:       x,
		Y:       y,
		Outline: outline,
		Fill:    fill,
		Shape: Triangle{
			Size: size,
		},
	}
}

func (triangle Triangle) Hitbox(shape *Shape) Hitbox {
	shape.hitbox.RegistrationX = int16(triangle.Size / 2)
	shape.hitbox.RegistrationY = 0
	shape.hitbox.MinX = shape.X - int16(triangle.Size/2)
	shape.hitbox.MaxX = shape.X + int16(triangle.Size/2)
	shape.hitbox.MinY = shape.Y
	shape.hitbox.MaxY = shape.Y + int16(triangle.Size/2)
	shape.hitbox.SizeX = shape.hitbox.MaxX - shape.hitbox.MinX
	shape.hitbox.SizeY = shape.hitbox.MaxY - shape.hitbox.MinY
	return shape.hitbox
}

func (triangle Triangle) Plot(display drivers.Displayer, shape *Shape) {
	triangle.Hitbox(shape)

	tinydraw.FilledTriangle(
		display,
		shape.X,
		shape.Y,
		shape.X-int16(triangle.Size/2),
		shape.Y+int16(triangle.Size/2),
		shape.X+int16(triangle.Size/2),
		shape.Y+int16(triangle.Size/2),
		shape.Fill,
	)
	tinydraw.Triangle(
		display,
		shape.X,
		shape.Y,
		shape.X-int16(triangle.Size/2),
		shape.Y+int16(triangle.Size/2),
		shape.X+int16(triangle.Size/2),
		shape.Y+int16(triangle.Size/2),
		shape.Outline,
	)
}
