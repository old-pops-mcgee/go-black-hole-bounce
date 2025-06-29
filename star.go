package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const STAR_RENDER_SCALE float32 = 1.5

type Star struct {
	game              *Game
	pos               rl.Vector2
	radius            float32
	angle             float32
	turningDirection  bool
	timeToDetonation  int32
	detonationCounter int32
}

func initStar(g *Game, p rl.Vector2, r float32) Star {
	detonationVal := int32(rand.Intn(300)) + 300
	return Star{
		game:              g,
		pos:               p,
		radius:            r,
		angle:             rand.Float32() * 2 * math.Pi,
		turningDirection:  rand.Intn(2) > 0,
		timeToDetonation:  detonationVal,
		detonationCounter: detonationVal,
	}
}

func (s *Star) update() {
	if s.turningDirection {
		s.angle += math.Pi / 120
	} else {
		s.angle -= math.Pi / 120
	}

	s.detonationCounter -= 1
}

func (s *Star) render() {
	var color rl.Color
	if float32(s.detonationCounter) < (float32(s.timeToDetonation) * float32(0.33)) {
		color = rl.Red
	} else if float32(s.detonationCounter) < (float32(s.timeToDetonation) * float32(0.67)) {
		color = rl.Orange
	} else {
		color = rl.Yellow
	}

	t := s.game.starTexture

	rl.DrawTexturePro(
		s.game.starTexture,
		rl.NewRectangle(0, 0, float32(t.Width), float32(t.Height)),
		rl.NewRectangle(s.pos.X, s.pos.Y, float32(t.Width)*STAR_RENDER_SCALE, float32(t.Height)*STAR_RENDER_SCALE),
		rl.Vector2{
			X: (float32(t.Width) / 2) * STAR_RENDER_SCALE,
			Y: (float32(t.Height) / 2) * STAR_RENDER_SCALE,
		},
		s.angle*(180/math.Pi),
		color,
	)
}
