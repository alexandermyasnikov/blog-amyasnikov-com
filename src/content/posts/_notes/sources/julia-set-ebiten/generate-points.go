func generatePoints(toCenterX, toCenterY, toScaleX, toScaleY float64, width, height int) []Point {
	points := make([]Point, 0, width*height)

	var (
		fromCenterX = float64(width) / 2
		fromCenterY = float64(height) / 2
		fromScaleX  = float64(width) / 2
		fromScaleY  = float64(height) / 2
	)

	for x := range width {
		for y := range height {
			nx := remapPoint(float64(x), fromCenterX, fromScaleX, toCenterX, toScaleX)
			ny := remapPoint(float64(y), fromCenterY, fromScaleY, toCenterY, toScaleY)
			col := computeColor(nx, ny, maxIters, cCoeff)
			points = append(points, Point{x, y, col})
		}
	}

	return points
}

func remapPoint(x, fromCenter, fromScale, toCenter, toScale float64) float64 {
	return (x-fromCenter)*(toScale/fromScale) + toCenter
}

func computeColor(x, y float64, maxIters int, c complex128) color.RGBA {
	z := complex(x, y)

	for i := range maxIters {
		if cmplx.Abs(z) > 2 {
			r := uint8((3*i + 17) % 0xFF)
			g := uint8((2*i + 01) % 0xFF)
			b := uint8((1*i + 01) % 0xFF)
			return color.RGBA{r, g, b, 0xFF}
		}

		z = z*z + c
	}

	return color.RGBA{}
}
