package seed

import "image/color"
import "encoding/hex"

const White Hex = "#ffffff"
const Black Hex = "#000000"
const Red Hex = "#ff0000"
const Green Hex = "#00ff00"
const Blue Hex = "#0000ff"

func RGB(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func RGBA(r, g, b, a uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

type Hex string

func (h Hex) RGBA() (r, g, b, a uint32) {
	var c [4]byte
	c[3] = 255
	hex.Decode(c[:], []byte(h[1:]))
	r, g, b, a = uint32(c[0]), uint32(c[1]), uint32(c[2]), uint32(c[3])
	if a != 255 {
		r = uint32((float64(r) / 255) * 65535)
		g = uint32((float64(g) / 255) * 65535)
		b = uint32((float64(b) / 255) * 65535)
		a = uint32((float64(a) / 255) * 65535)
	}
	return
}
