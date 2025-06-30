package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Asteroid struct {
	game       *Game
	pos        rl.Vector2
	radius     float32
	velocity   rl.Vector2
	vaporTrail []rl.Vector3
	isAlive    bool
}

func initAsteroid(g *Game, p rl.Vector2, r float32, v rl.Vector2) Asteroid {
	return Asteroid{
		game:       g,
		pos:        p,
		radius:     r,
		velocity:   v,
		vaporTrail: []rl.Vector3{},
		isAlive:    true,
	}
}

func (a *Asteroid) render() {
	rl.DrawTexture(a.game.asteroidTexture, int32(a.pos.X), int32(a.pos.Y), rl.White)

	for _, dot := range a.vaporTrail {
		rl.DrawCircle(int32(dot.X), int32(dot.Y), dot.Z, color.RGBA{255, 190, 51, 100})
	}
}

func (a *Asteroid) update() {
	// Remove if we've gone out of bounds
	if a.pos.Y < -50 || a.pos.Y > float32(WindowHeight+50) || a.pos.X < -50 || a.pos.X > float32(WindowWidth+50) {
		a.isAlive = false
		return
	}

	// Check to see if we've been crushed
	for _, b := range a.game.blackHoleList {
		if rl.CheckCollisionCircles(a.pos, a.radius, b.pos, b.deathRadius) {
			a.isAlive = false
			a.game.createNewExplosion(a.pos, 15)
			return
		}

		// Calculate velocity updates from black hole
		force := b.calculateForceOnObject(a.pos)
		a.velocity = rl.Vector2Add(a.velocity, force)
	}

	// Calculate the new position
	a.pos = rl.Vector2Add(a.pos, a.velocity)

	// Update the vapor trail
	newVaporTrail := []rl.Vector3{}
	for _, dot := range a.vaporTrail {
		dot.Z -= 0.1
		if dot.Z > 0 {
			newVaporTrail = append(newVaporTrail, dot)
		}
	}

	// Add a new dot to the trail
	t := a.game.asteroidTexture
	newVaporTrail = append(newVaporTrail, rl.Vector3{
		X: a.pos.X + float32(t.Width)/2,
		Y: a.pos.Y + float32(t.Height)/2,
		Z: 3.0,
	})
	a.vaporTrail = newVaporTrail

}

func (a *Asteroid) getCollisionCircle() rl.Vector3 {
	t := a.game.asteroidTexture
	return rl.Vector3{
		X: a.pos.X + (float32(t.Width) / 2),
		Y: a.pos.Y + (float32(t.Height) / 2),
		Z: a.radius,
	}
}
