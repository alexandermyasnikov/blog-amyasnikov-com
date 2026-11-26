func saveImage(toCenterX, toCenterY, toScaleX, toScaleY float64, width, height int) error {
	width *= scaleSaveImage
	height *= scaleSaveImage

	points := generatePoints(toCenterX, toCenterY, toScaleX, toScaleY, width, height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, p := range points {
		img.Set(p.X, p.Y, p.Color)
	}

	filename := fmt.Sprintf("/tmp/fractal_%s.png", time.Now().Format("20060102-150405"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
