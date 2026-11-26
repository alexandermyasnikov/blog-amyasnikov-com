type Game struct {
	width  int
	height int

	CenterX, CenterY float64
	ScaleX, ScaleY   float64

	points      []Point
	needsUpdate bool
}

func (g *Game) Update() error {
	g.handleInput()

	if g.needsUpdate {
		g.points = generatePoints(g.CenterX, g.CenterY, g.ScaleX, g.ScaleY, g.width, g.height)
		g.needsUpdate = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.points {
		screen.Set(p.X, p.Y, p.Color)
	}

	msg := fmt.Sprintf("TPS: %0.2f FPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += fmt.Sprintf("Center: (%v,%v), Scale: (%v,%v)\n", g.CenterX, g.CenterY, g.ScaleX, g.ScaleY)
	msg += fmt.Sprintf("maxIters: %v, cCoeff: %v\n", maxIters, cCoeff)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("2D Fractal Viever - Demo - amyasnikov.com")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	game := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
