package color

import "fmt"
import "math"

type RGB struct {
	R, G, B float64
}

func hexToRGB(s string) (*RGB, error) {
	format := "#%02x%02x%02x"
	factor := 1.0 / 255.0

	var r, g, b uint8
	n, err := fmt.Sscanf(s, format, &r, &g, &b)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, fmt.Errorf("conversion failed: %v is not a hex color", s)
	}

	return &RGB{float64(r) * factor, float64(g) * factor, float64(b) * factor}, nil
}

func applyGammaCorrection(value float64) float64 {
	if value > 0.04045 {
		return math.Pow((value+0.055)/(1.0+0.055), 2.4)
	}
	return value / 12.92
}

func RGBToXY(s string) (float64, float64, error) {
	rgb, err := hexToRGB(s)
	if err != nil {
		return -1.0, -1.0, err
	}

	// https://github.com/PhilipsHue/PhilipsHueSDK-iOS-OSX/blob/master/ApplicationDesignNotes/RGB%20to%20xy%20Color%20conversion.md
	red := applyGammaCorrection(rgb.R)
	green := applyGammaCorrection(rgb.G)
	blue := applyGammaCorrection(rgb.B)

	X := red*0.649926 + green*0.103455 + blue*0.197109
	Y := red*0.234327 + green*0.743075 + blue*0.022598
	Z := red*0.0000000 + green*0.053077 + blue*1.035763
	x := X / (X + Y + Z)
	y := Y / (X + Y + Z)

	return x, y, nil
}
