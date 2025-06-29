package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Asteroid struct {
	game       *Game
	pos        rl.Vector2
	radius     float32
	velocity   rl.Vector2
	vaporTrail []rl.Vector3
}

func initAsteroid(g *Game, p rl.Vector2, r float32, v rl.Vector2) Asteroid {
	return Asteroid{
		game:       g,
		pos:        p,
		radius:     r,
		velocity:   v,
		vaporTrail: []rl.Vector3{},
	}
}
