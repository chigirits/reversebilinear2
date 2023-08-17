package main

import (
	"image/color"
	"math"
)

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if 1 < v {
		return 1
	}
	return v
}

type Pixel struct {
	r float64
	g float64
	b float64
}

func (p *Pixel) ToColor() color.Color {
	return color.RGBA{
		R: uint8(math.Round(clamp(p.r) * 255)),
		G: uint8(math.Round(clamp(p.g) * 255)),
		B: uint8(math.Round(clamp(p.b) * 255)),
		A: 255,
	}
}

func (p *Pixel) Mul(v float64) *Pixel {
	return &Pixel{
		r: p.r * v,
		g: p.g * v,
		b: p.b * v,
	}
}

func AddPixel4(a, b, c, d *Pixel) *Pixel {
	return &Pixel{
		r: a.r + b.r + c.r + d.r,
		g: a.g + b.g + c.g + d.g,
		b: a.b + b.b + c.b + d.b,
	}
}

func NewPixel(c color.Color) *Pixel {
	r, g, b, _ := c.RGBA()
	result := &Pixel{
		r: float64(r) / 0xffff,
		g: float64(g) / 0xffff,
		b: float64(b) / 0xffff,
	}
	return result
}
