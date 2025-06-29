package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var colorChoices []rl.Color = []rl.Color{rl.Red, rl.Orange, rl.Yellow, rl.Gold}

type Explosion struct {
	cluster *ExplosionCluster
	pos     rl.Vector2
	angle   float32
	speed   float32
	color   rl.Color
}

func initExplosion(c *ExplosionCluster, p rl.Vector2, a float32, s float32) Explosion {
	return Explosion{
		cluster: c,
		pos:     p,
		angle:   a,
		speed:   s,
		color:   colorChoices[rand.Intn(len(colorChoices))],
	}
}

func (e *Explosion) render() {
	renderPoints := []rl.Vector2{
		// Point 0
		rl.Vector2Add(e.pos, rl.Vector2{
			X: float32(math.Cos(float64(e.angle)) * float64(e.speed) * (rand.Float64()*3 + 2)),
			Y: float32(math.Sin(float64(e.angle)) * float64(e.speed) * (rand.Float64()*3 + 2)),
		}),
		// Point 1
		rl.Vector2Add(e.pos, rl.Vector2{
			X: float32(math.Cos(float64(e.angle+math.Pi*0.5)) * float64(e.speed) * (rand.Float64()*3 + 2)),
			Y: float32(math.Sin(float64(e.angle+math.Pi*0.5)) * float64(e.speed) * (rand.Float64()*3 + 2)),
		}),
		// Point 2
		rl.Vector2Add(e.pos, rl.Vector2{
			X: float32(math.Cos(float64(e.angle+math.Pi)) * float64(e.speed) * (rand.Float64()*3 + 2)),
			Y: float32(math.Sin(float64(e.angle+math.Pi)) * float64(e.speed) * (rand.Float64()*3 + 2)),
		}),
		// Point 3
		rl.Vector2Add(e.pos, rl.Vector2{
			X: float32(math.Cos(float64(e.angle+math.Pi*1.5)) * float64(e.speed) * (rand.Float64()*3 + 2)),
			Y: float32(math.Sin(float64(e.angle+math.Pi*1.5)) * float64(e.speed) * (rand.Float64()*3 + 2)),
		}),
	}

	rl.DrawTriangle(renderPoints[0], renderPoints[3], renderPoints[1], e.color)
	rl.DrawTriangle(renderPoints[2], renderPoints[1], renderPoints[3], e.color)
}

func (e *Explosion) update() {
	e.pos = rl.Vector2Add(e.pos, rl.Vector2{X: float32(math.Cos(float64(e.angle)) * float64(e.speed)), Y: float32(math.Sin(float64(e.angle)) * float64(e.speed))})
	e.speed -= 0.1
}
