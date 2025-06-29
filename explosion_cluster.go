package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ExplosionCluster struct {
	game       *Game
	pos        rl.Vector2
	explosions []Explosion
}

func initExplosionCluster(g *Game, p rl.Vector2, explosionCount int32) ExplosionCluster {
	cluster := ExplosionCluster{
		game: g,
		pos:  p,
	}

	// Initialize the explosion cluster
	var explosions []Explosion
	for range explosionCount {
		pos := rl.Vector2{X: p.X + rand.Float32()*5 - 5, Y: p.Y + rand.Float32()*5 - 5}
		explosions = append(explosions, initExplosion(&cluster, pos, rand.Float32()*2*math.Pi, rand.Float32()*5.0))
	}
	cluster.explosions = explosions

	rl.PlaySound(g.explosionSound)

	return cluster
}

func (c *ExplosionCluster) render() {
	for _, explosion := range c.explosions {
		explosion.render()
	}
}

func (c *ExplosionCluster) update() {
	newExplosionList := []Explosion{}
	for _, explosion := range c.explosions {
		explosion.update()
		if explosion.speed > 0 {
			newExplosionList = append(newExplosionList, explosion)
		}
	}
	c.explosions = newExplosionList
}
