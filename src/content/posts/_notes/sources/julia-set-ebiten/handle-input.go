func (g *Game) handleInput() {
	const scaleStep = 0.3
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.CenterX -= scaleStep * g.ScaleX
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.CenterX += scaleStep * g.ScaleX
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.CenterY -= scaleStep * g.ScaleY
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.CenterY += scaleStep * g.ScaleY
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.ScaleX *= 1 + scaleStep
		g.ScaleY *= 1 + scaleStep
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.ScaleX *= 1 - scaleStep
		g.ScaleY *= 1 - scaleStep
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.needsUpdate = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		log.Println("Saving image ...")
		go func() {
			if err := saveImage(g.CenterX, g.CenterY, g.ScaleX, g.ScaleY, g.width, g.height); err != nil {
				log.Printf("Failed to save image: %v", err)
			} else {
				log.Printf("Image saved")
			}
		}()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		os.Exit(0)
	}
}
