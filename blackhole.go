package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const STANDARD_FORCE float32 = 2500.0
const FUDGE_FACTOR float64 = 3.5
const DECAYING_FORCE_ADDER float32 = 20.0
const DECAY_RATE float32 = 0.25
const RENDER_SCALE float32 = 3.5

type BlackHole struct {
	game             *Game
	pos              rl.Vector2
	initialRadius    float32
	radius           float32
	force            float32
	level            int
	deathRadius      float32
	angle            float32
	turningDirection bool
	rotationSpeed    float32
}

func initBlackHole(g *Game, p rl.Vector2, r float32) BlackHole {
	return BlackHole{
		game:             g,
		pos:              p,
		initialRadius:    r,
		radius:           r,
		force:            STANDARD_FORCE,
		level:            1,
		deathRadius:      0.4 * r,
		angle:            rand.Float32() * 2 * math.Pi,
		turningDirection: rand.Intn(2) > 0,
		rotationSpeed:    rand.Float32() * math.Pi / 15,
	}
}

func (b *BlackHole) render() {
	scale := b.radius / b.initialRadius
	t := b.game.blackHoleTexture
	scaledTWidth := RENDER_SCALE * float32(t.Width) * scale
	scaledTHeight := RENDER_SCALE * float32(t.Height) * scale
	rl.DrawTexturePro(
		t,
		rl.NewRectangle(0, 0, float32(t.Width), float32(t.Height)),
		rl.NewRectangle(b.pos.X, b.pos.Y, scaledTWidth, scaledTHeight),
		rl.Vector2{
			X: scaledTWidth / 2,
			Y: scaledTHeight / 2,
		},
		b.angle*(180/math.Pi),
		rl.White,
	)
}

func (b *BlackHole) update() {
	b.radius -= DECAY_RATE
	b.deathRadius -= 0.4 * DECAY_RATE
	b.force += DECAYING_FORCE_ADDER
	if b.turningDirection {
		b.angle += b.rotationSpeed
	} else {
		b.angle -= b.rotationSpeed
	}
}

func (b *BlackHole) calculateForceOnObject(obj rl.Vector2) rl.Vector2 {
	angle := math.Atan2(float64(b.pos.Y-obj.Y), float64(b.pos.X-obj.X))
	dis := math.Sqrt(math.Pow(float64(b.pos.Y-obj.Y), 2) + math.Pow(float64(b.pos.X-obj.X), 2))

	gForce := float64(b.force) / (FUDGE_FACTOR * math.Pow(dis, 2))
	return rl.Vector2{X: float32(math.Cos(angle) * gForce), Y: float32(math.Sin(angle) * gForce)}
}
