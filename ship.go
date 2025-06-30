package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SPEED float64 = 5
const MAX_ENGINE_SPEED float64 = 50

type Ship struct {
	game        *Game
	pos         rl.Vector2 // x, y
	radius      float32
	angle       float32
	velocity    rl.Vector2 // x velocity, y velocity
	engineSpeed float64
	vaporTrail  []rl.Vector3 // x position, y position, size
	isDead      bool
	texture     rl.Texture2D
}

func initShip(g *Game, p rl.Vector2, r float32, t rl.Texture2D) Ship {
	return Ship{
		game:        g,
		pos:         p,
		radius:      r,
		angle:       0,
		velocity:    rl.Vector2{X: 0, Y: 0},
		engineSpeed: 0,
		vaporTrail:  []rl.Vector3{},
		isDead:      false,
		texture:     t,
	}
}

func (s *Ship) render() {
	if !s.isDead {
		fTextureWidth := float32(s.texture.Width)
		fTextureHeight := float32(s.texture.Height)
		rl.DrawTexturePro(s.texture, rl.NewRectangle(0, 0, fTextureWidth, fTextureHeight), rl.NewRectangle(s.pos.X, s.pos.Y, fTextureWidth, fTextureHeight), rl.Vector2{X: fTextureWidth / 2, Y: fTextureHeight / 2}, s.angle*(180/math.Pi)+90, rl.White)
	}
	for _, dot := range s.vaporTrail {
		rl.DrawCircle(int32(dot.X), int32(dot.Y), dot.Z, color.RGBA{255, 95, 31, 100})
	}
}

func (s *Ship) update() {
	if !s.isDead {
		// Put a floor on the engine speed
		s.engineSpeed = math.Min(MAX_ENGINE_SPEED, math.Max(0, float64(s.engineSpeed)))

		// Handle engine sound
		if s.engineSpeed > 0 && !rl.IsSoundPlaying(s.game.engineSound) {
			rl.PlaySound(s.game.engineSound)
		}
		if s.engineSpeed <= 0 && rl.IsSoundPlaying(s.game.engineSound) {
			rl.StopSound(s.game.engineSound)
		}
		if rl.IsSoundPlaying(s.game.engineSound) {
			rl.SetSoundVolume(s.game.engineSound, float32(12.5*s.engineSpeed/MAX_ENGINE_SPEED))
		}

		// Blow up if we've gone out of bounds
		if s.pos.Y < 0 || s.pos.Y > float32(WindowHeight) || s.pos.X < 0 || s.pos.X > float32(WindowWidth) {
			s.isDead = true
			return
		}

		// Check asteroid collisions
		for _, asteroid := range s.game.asteroidList {
			asteroidCollisionCircle := asteroid.getCollisionCircle()
			if rl.CheckCollisionCircles(s.pos, s.radius, rl.Vector2{X: asteroidCollisionCircle.X, Y: asteroidCollisionCircle.Y}, asteroidCollisionCircle.Z) {
				s.isDead = true
				asteroid.isAlive = false
				s.game.createNewExplosion(asteroid.pos, 15)
				return
			}
		}

		// Check black hole collisions
		for _, blackHole := range s.game.blackHoleList {
			if rl.CheckCollisionCircles(s.pos, s.radius, blackHole.pos, blackHole.deathRadius) {
				fmt.Println("collided with black hole")
				fmt.Printf("s.pos: %v\ns.radius: %v\nb.pos: %v\nb.radius: %v\n", s.pos, s.radius, blackHole.pos, blackHole.deathRadius)
				fmt.Println(blackHole)
				s.isDead = true
				return
			}

			// Calculate velocity updates from black hole
			s.velocity = rl.Vector2Add(s.velocity, blackHole.calculateForceOnObject(s.pos))
		}

		// Calculate new position
		s.pos = rl.Vector2Add(s.pos, rl.Vector2{
			X: float32(math.Cos(float64(s.angle))*s.engineSpeed) + s.velocity.X,
			Y: float32(math.Sin(float64(s.angle))*s.engineSpeed) + s.velocity.Y,
		})

		// Cap velocities so we don't get too crazy
		s.velocity = rl.Vector2{
			X: float32(math.Min(MAX_SPEED, math.Max(-MAX_SPEED, float64(s.velocity.X)))),
			Y: float32(math.Min(MAX_SPEED, math.Max(-MAX_SPEED, float64(s.velocity.Y)))),
		}

		// Update the vapor trail
		newVaporTrail := []rl.Vector3{}
		for _, dot := range s.vaporTrail {
			dot.Z -= 0.1
			if dot.Z > 0 {
				newVaporTrail = append(newVaporTrail, dot)
			}
		}

		// Add vapor dots to the trail
		vaporFudgeFactor := float64(rand.Float32()-0.5) * 8.0
		theta := float64((math.Pi / 2) - s.angle)
		vaporDot := rl.Vector3{
			X: float32(float64(s.pos.X) - (float64(s.texture.Height)+vaporFudgeFactor)*math.Sin(theta)/2),
			Y: float32(float64(s.pos.Y) - (float64(s.texture.Height)+vaporFudgeFactor)*math.Cos(theta)/2),
			Z: float32(s.engineSpeed) / 2,
		}
		newVaporTrail = append(newVaporTrail, vaporDot)

		// Save new vapor trail to ship
		s.vaporTrail = newVaporTrail

	}
}

func (s *Ship) increaseSpeed() {
	s.engineSpeed += 0.1
}

func (s *Ship) decreaseSpeed() {
	s.engineSpeed -= 0.1
}
