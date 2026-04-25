package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// PointOfInterest representa um setor ou objetivo no mapa
type PointOfInterest struct {
	X, Y float64
	Name string
}

// Star cria o efeito de fundo
type Star struct {
	X, Y, Speed float64
}

type Game struct {
	playerX, playerY float64
	stars            []Star
	sectors          []PointOfInterest
	currentMission   int
}

func (g *Game) Update() error {
	// 1. Controles de Movimentação (Aumentei a velocidade para facilitar a navegação)
	speed := 7.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerY += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerX += speed
	}

	// 2. Lógica de Missão: Checar se chegou no objetivo atual
	if g.currentMission < 20 {
		target := g.sectors[g.currentMission]
		dist := math.Sqrt(math.Pow(target.X-g.playerX, 2) + math.Pow(target.Y-g.playerY, 2))

		// Se estiver a menos de 60 pixels, missão cumprida!
		if dist < 60 {
			g.currentMission++
			fmt.Printf("Missão %d de 20 completa!\n", g.currentMission)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fundo Espacial
	screen.Fill(color.RGBA{5, 5, 15, 255})

	// 3. Desenhar Estrelas (Efeito Parallax Simples)
	for _, s := range g.stars {
		// As estrelas se movem levemente opostas ao jogador para dar profundidade
		sx := math.Mod(s.X-g.playerX*0.5, float64(screenWidth))
		if sx < 0 {
			sx += float64(screenWidth)
		}
		sy := math.Mod(s.Y-g.playerY*0.5, float64(screenHeight))
		if sy < 0 {
			sy += float64(screenHeight)
		}

		screen.Set(int(sx), int(sy), color.RGBA{150, 150, 150, 255})
	}

	// 4. Desenhar os 150 Setores
	for i, p := range g.sectors {
		relX := p.X - g.playerX + (screenWidth / 2)
		relY := p.Y - g.playerY + (screenHeight / 2)

		// Renderiza apenas se estiver na tela
		if relX > -100 && relX < screenWidth+100 && relY > -100 && relY < screenHeight+100 {
			// Diferencia o objetivo atual dos outros setores
			c := color.RGBA{0, 150, 255, 255} // Azul para setores comuns
			if i == g.currentMission {
				c = color.RGBA{255, 200, 0, 255} // Amarelo para o objetivo
			}
			ebitenutil.DrawRect(screen, relX-10, relY-10, 20, 20, c)
			ebitenutil.DebugPrintAt(screen, p.Name, int(relX)-20, int(relY)+15)
		}
	}

	// 5. Radar / Bússola (Aponta para a missão atual)
	if g.currentMission < 20 {
		target := g.sectors[g.currentMission]
		angle := math.Atan2(target.Y-g.playerY, target.X-g.playerX)

		// Desenha uma seta indicadora ao redor do jogador
		pointerX := float64(screenWidth/2) + math.Cos(angle)*50
		pointerY := float64(screenHeight/2) + math.Sin(angle)*50
		ebitenutil.DrawLine(screen, float64(screenWidth/2), float64(screenHeight/2), pointerX, pointerY, color.RGBA{255, 255, 0, 255})
	}

	// 6. Desenhar Jogador (Nave)
	ebitenutil.DrawRect(screen, screenWidth/2-8, screenHeight/2-8, 16, 16, color.White)

	// HUD - Interface
	info := fmt.Sprintf("COORDS: %.0f, %.0f\nMISSAO ATUAL: %d/20\nOBJETIVO: %s",
		g.playerX, g.playerY, g.currentMission, g.sectors[g.currentMission].Name)
	ebitenutil.DebugPrint(screen, info)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Gerar 150 Setores Proceduralmente
	sectors := []PointOfInterest{}
	for i := 0; i < 150; i++ {
		sectors = append(sectors, PointOfInterest{
			X:    float64(rand.Intn(20000) - 10000), // Mapa de 20.000 pixels
			Y:    float64(rand.Intn(20000) - 10000),
			Name: fmt.Sprintf("Setor-%03d", i),
		})
	}

	// Gerar Estrelas de fundo
	stars := []Star{}
	for i := 0; i < 150; i++ {
		stars = append(stars, Star{
			X: rand.Float64() * screenWidth,
			Y: rand.Float64() * screenHeight,
		})
	}

	game := &Game{
		sectors: sectors,
		stars:   stars,
	}

	ebiten.SetWindowTitle("Go Space Open World - 150 Níveis")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
