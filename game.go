package main

import (
	"embed"
	_ "embed"
	"fmt"
	"image/color"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets
var ASSETS embed.FS

const WindowHeight int = 840
const WindowWidth int = 1540
const MaxStars int = 5

type Game struct {
	// Game state
	gameState State

	// Textures
	backgroundTexture rl.Texture2D
	shipTexture       rl.Texture2D
	asteroidTexture   rl.Texture2D
	blackHoleTexture  rl.Texture2D
	starTexture       rl.Texture2D

	// Sounds
	music          rl.Sound
	explosionSound rl.Sound
	engineSound    rl.Sound

	// Game components
	ship                   Ship
	blackHoleList          []BlackHole
	starList               []Star
	asteroidList           []Asteroid
	explosionClusterList   []ExplosionCluster
	addedFinalExplosion    bool
	restartCounter         int32
	asteroidCountdownRange rl.Vector2
	starAdditionCountdown  int32
	starMultiplier         int32
	asteroidCountdown      int32
	score                  int32
}

type State int

const (
	Start State = iota
	Play
	Restart
)

func init() {
	rl.AddFileSystem(ASSETS)
}

func initGame() Game {
	// Init game contexts
	rl.InitWindow(int32(WindowWidth), int32(WindowHeight), "Black Hole Bounce")
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	// Load the assets and return game element
	return Game{
		gameState:         Start,
		backgroundTexture: rl.LoadTexture("assets/images/background.png"),
		shipTexture:       rl.LoadTexture("assets/images/ship.png"),
		asteroidTexture:   rl.LoadTexture("assets/images/asteroid.png"),
		blackHoleTexture:  rl.LoadTexture("assets/images/black_hole.png"),
		starTexture:       rl.LoadTexture("assets/images/star.png"),
		music:             rl.LoadSound("assets/sound/scifi_background.wav"),
		explosionSound:    rl.LoadSound("assets/sound/explosion.wav"),
		engineSound:       rl.LoadSound("assets/sound/engine.wav"),
	}
}

// Reload game components (resets to starting state)
func (g *Game) reloadGameComponents() {
	g.ship = initShip(g, rl.Vector2{X: float32(WindowWidth) / 2, Y: float32(WindowHeight) / 2}, 5, g.shipTexture)
	g.blackHoleList = []BlackHole{}
	g.asteroidList = []Asteroid{}
	g.explosionClusterList = []ExplosionCluster{}
	g.addedFinalExplosion = false
	g.restartCounter = 120
	g.asteroidCountdownRange = rl.Vector2{X: 180, Y: 300}
	g.starAdditionCountdown = 1800
	g.starMultiplier = 1
	g.asteroidCountdown = int32(rand.Intn(int(g.asteroidCountdownRange.Y-g.asteroidCountdownRange.X))) + int32(g.asteroidCountdownRange.X)
	starList := []Star{}
	for range MaxStars {
		starList = append(starList, g.generateRandomStar())
	}
	g.starList = starList
	g.score = 0
}

// Unload the loaded assets before closing the game
func (g *Game) unload() {
	rl.UnloadTexture(g.backgroundTexture)
	rl.UnloadTexture(g.shipTexture)
	rl.UnloadTexture(g.asteroidTexture)
	rl.UnloadTexture(g.blackHoleTexture)
	rl.UnloadTexture(g.starTexture)
	rl.UnloadSound(g.music)
	rl.UnloadSound(g.explosionSound)
	rl.UnloadSound(g.engineSound)
}

// Main loop of the game
func (g *Game) run() {
	// Boot initial tasks
	rl.PlaySound(g.music)
	g.reloadGameComponents()

	var update = func() {
		g.handleInput()
		g.update()
		g.render()
	}
	rl.SetMainLoop(update)
	for !rl.WindowShouldClose() {
		update()
	}

	g.unload()
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

// Render textures to screen
func (g *Game) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	// Draw the background
	rl.DrawTexture(g.backgroundTexture, 0, 0, color.RGBA{255, 255, 255, 200})

	switch g.gameState {
	case Play:
		// Render stars
		for _, star := range g.starList {
			star.render()
		}

		// Render black holes
		for _, blackHole := range g.blackHoleList {
			blackHole.render()
		}

		// Render asteroids
		for _, asteroid := range g.asteroidList {
			fmt.Println("Afdsa")
			asteroid.render()
		}

		// Render explosions
		for _, cluster := range g.explosionClusterList {
			cluster.render()
		}

		// Render the ship
		g.ship.render()

		// Draw UI elements
		rl.DrawText(fmt.Sprintf("Score: %d", g.score), 10, 10, 40, rl.RayWhite)
	case Start:
		rl.DrawText("Black Hole Bounce", 550, 300, 64, rl.RayWhite)
		rl.DrawText("Stay Alive as Long as You Can!", 350, 350, 64, rl.RayWhite)

		rl.DrawText("Left/Right arrows: Turn", 550, 450, 48, rl.RayWhite)
		rl.DrawText("Up/Down arrows: Accelerate/Decelerate", 350, 500, 48, rl.RayWhite)
		rl.DrawText("Press Space to Start!", 470, 600, 64, rl.RayWhite)
	case Restart:
		rl.DrawText(fmt.Sprintf("Final Score: %d", g.score), 570, 300, 64, rl.RayWhite)
		rl.DrawText("Play Again?", 635, 350, 64, rl.RayWhite)

		rl.DrawText("Left/Right arrows: Turn", 550, 450, 48, rl.RayWhite)
		rl.DrawText("Up/Down arrows: Accelerate/Decelerate", 350, 500, 48, rl.RayWhite)
		rl.DrawText("Press Space to Start!", 470, 600, 64, rl.RayWhite)
	}

	rl.EndDrawing()
}

// Process game logic updates
func (g *Game) update() {
	if !rl.IsSoundPlaying(g.music) {
		rl.PlaySound(g.music)
	}

	if g.gameState == Play {
		// Process end-game components
		if g.ship.isDead {
			if !g.addedFinalExplosion {
				g.createNewExplosion(g.ship.pos, 50)
				g.addedFinalExplosion = true
			}
			g.restartCounter -= 1
		}
		if g.restartCounter <= 0 {
			g.gameState = Restart
		}

		// Increase the score
		if !g.ship.isDead {
			g.score += 1
		}

		// Process asteroid event
		g.asteroidCountdown -= 1
		if g.asteroidCountdown <= 0 {
			g.createNewAsteroid()
		}

		// Process star event
		g.starAdditionCountdown -= 1
		if g.starAdditionCountdown <= 0 {
			g.starMultiplier = int32(math.Min(5, float64(g.starMultiplier+1)))
			g.starAdditionCountdown = 1800
		}

		// Update the ship
		g.ship.update()

		// Update the black holes
		newBlackholeList := []BlackHole{}
		for _, blackHole := range g.blackHoleList {
			blackHole.update()

			if blackHole.radius > DECAY_RATE {
				newBlackholeList = append(newBlackholeList, blackHole)
			} else {
				for range g.starMultiplier {
					g.starList = append(g.starList, g.generateRandomStar())
				}
			}
		}
		g.blackHoleList = newBlackholeList

		// Update the stars
		newStarList := []Star{}
		for _, star := range g.starList {
			star.update()

			if star.detonationCounter > 0 {
				newStarList = append(newStarList, star)
			} else {
				g.addBlackHole(star.pos)
			}
		}
		g.starList = newStarList

		// Update the asteroids
		newAsteroidList := []Asteroid{}
		for _, asteroid := range g.asteroidList {
			asteroid.update()
			if asteroid.isAlive {
				newAsteroidList = append(newAsteroidList, asteroid)
			}
		}
		g.asteroidList = newAsteroidList

		// Update the explosions
		newExplosionClusterList := []ExplosionCluster{}
		for _, cluster := range g.explosionClusterList {
			cluster.update()

			if len(cluster.explosions) > 0 {
				newExplosionClusterList = append(newExplosionClusterList, cluster)
			}
		}
		g.explosionClusterList = newExplosionClusterList

	}
}

// Handle user input
func (g *Game) handleInput() {
	switch g.gameState {
	case Play:
		if rl.IsKeyDown(rl.KeyRight) {
			g.ship.angle += math.Pi / 60
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			g.ship.angle -= math.Pi / 60
		}
		if rl.IsKeyDown(rl.KeyUp) {
			g.ship.increaseSpeed()
		}
		if rl.IsKeyDown(rl.KeyDown) {
			g.ship.decreaseSpeed()
		}
	case Start, Restart:
		if rl.IsKeyPressed(rl.KeySpace) {
			g.reloadGameComponents()
			g.gameState = Play
		}
	}
}

func (g *Game) addBlackHole(p rl.Vector2) {
	g.blackHoleList = append(g.blackHoleList, initBlackHole(g, p, 45.0))
}

func (g *Game) generateRandomStar() Star {
	return initStar(g, rl.Vector2{X: float32(rand.Intn(WindowWidth-40) + 20), Y: float32(rand.Intn(WindowHeight-40) + 20)}, 5)
}

func (g *Game) createNewAsteroid() {
	// Determine side first
	side := rand.Intn(4)
	initialVelocity := rl.Vector2{}
	directionOfFreeSide := rand.Float32() - 0.5
	if directionOfFreeSide < 0 {
		directionOfFreeSide = -1
	} else {
		directionOfFreeSide = 1
	}
	pos := rl.Vector2{}
	velocityScale := rand.Float32() * 20.0

	switch side {
	case 0:
		// Top side
		pos.X = float32(rand.Intn(WindowWidth-80) + 40)
		pos.Y = float32(WindowHeight - 20)
		initialVelocity.X = rand.Float32() * directionOfFreeSide * velocityScale
		initialVelocity.Y = rand.Float32() * velocityScale
	case 1:
		// Right side
		pos.X = float32(WindowWidth + 20)
		pos.Y = float32(rand.Intn(WindowHeight-80) + 40)
		initialVelocity.X = rand.Float32() * -velocityScale
		initialVelocity.Y = rand.Float32() * directionOfFreeSide * velocityScale
	case 2:
		// Bottom side
		pos.X = float32(rand.Intn(WindowWidth-80) + 40)
		pos.Y = float32(WindowHeight + 20)
		initialVelocity.X = rand.Float32() * directionOfFreeSide * velocityScale
		initialVelocity.Y = rand.Float32() * -velocityScale
	case 3:
		// Left side
		pos.X = float32(WindowWidth - 20)
		pos.Y = float32(rand.Intn(WindowHeight-80) + 40)
		initialVelocity.X = rand.Float32() * velocityScale
		initialVelocity.Y = rand.Float32() * directionOfFreeSide * velocityScale
	}

	g.asteroidList = append(g.asteroidList, initAsteroid(g, pos, 10.0, initialVelocity))
	fmt.Println(len(g.asteroidList))
	g.asteroidCountdownRange = rl.Vector2{
		X: float32(math.Max(20, float64(g.asteroidCountdownRange.X)-10)),
		Y: float32(math.Max(40, float64(g.asteroidCountdownRange.Y)-10)),
	}
	g.asteroidCountdown = int32(rand.Intn(int(g.asteroidCountdownRange.Y)) + int(g.asteroidCountdownRange.X))
}

func (g *Game) createNewExplosion(p rl.Vector2, e int32) {
	g.explosionClusterList = append(g.explosionClusterList, initExplosionCluster(g, p, e))
}
