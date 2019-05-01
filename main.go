package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const width int = 800
const height int = 600

// Position represents a 2D coordinate on the game screen
type Position struct {
	x, y float32
}

// Color represents an RGB value for a pixel
type Color struct {
	r, g, b byte
}

// Ball represents the ball being deflected
type Ball struct {
	Position
	radius int
	xv, yv float32
	color  Color
}

// Draw will draw the ball to the given screen
func (ball *Ball) Draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			// square to avoid sqrt
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

// Update the ball's position based on its x/y velocity
func (ball *Ball) Update() {
	ball.x += ball.xv
	ball.y += ball.yv

	if ball.y < 0 || int(ball.y) > height {
		ball.yv = -ball.yv
	}
}

// Paddle represents a deflection paddle on the screen
type Paddle struct {
	Position
	w, h  int
	color Color
}

// Draw will draw the paddle to the given screen
func (paddle *Paddle) Draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

// Update the Paddle's position based on its x/y velocity
func (paddle *Paddle) Update(arrowKey uint8) {
	switch arrowKey {
	case UpArrow:
		paddle.y--
	case DownArrow:
		paddle.y++
	default:
		return
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c Color, pixels []byte) {
	idx := (y*width + x) * 4

	if idx < len(pixels)-4 && idx >= 0 {
		pixels[idx] = c.r
		pixels[idx+1] = c.g
		pixels[idx+2] = c.b
		// pixels[idx+3] = c.a
	}
	// TODO Error case
}

// NoArrow was pressed
const NoArrow uint8 = 0

// UpArrow was pressed
const UpArrow uint8 = 1

// DownArrow was pressed
const DownArrow uint8 = 2

// LeftArrow was pressed
const LeftArrow uint8 = 3

// RightArrow was pressed
const RightArrow uint8 = 4

func getArrowPressed(keyState []uint8) uint8 {
	if keyState[sdl.SCANCODE_UP] != 0 {
		return UpArrow
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		return DownArrow
	} else if keyState[sdl.SCANCODE_LEFT] != 0 {
		return LeftArrow
	} else if keyState[sdl.SCANCODE_RIGHT] != 0 {
		return RightArrow
	} else {
		return NoArrow
	}
}

func main() {

	window, err := sdl.CreateWindow("Pongo", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, width*height*4) // *4 because RGBA

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			setPixel(x, y, Color{0, 0, 0}, pixels)
		}
	}

	// Game loop
	player1 := Paddle{Position{100, 100}, 20, 100, Color{255, 255, 255}}
	ball := Ball{Position{300, 300}, 20, 0, 0, Color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	for {

		for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			switch evt.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clear(pixels)

		player1.Update(getArrowPressed(keyState))

		player1.Draw(pixels)
		ball.Draw(pixels)

		tex.Update(nil, pixels, width*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}
