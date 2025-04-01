package image

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/nfnt/resize"

	"github.com/fogleman/gg"
)

const fontSize = 16

func (s *imgFile) TextWater(Text, fontPath string) (*imgFile, error) {

	bound := s.Img.Bounds()
	W := bound.Dx()
	H := bound.Dy()
	rd := W / 375
	if rd == 0 {
		rd = 1
	}

	dc := gg.NewContextForImage(s.Img)

	if err := dc.LoadFontFace(fontPath, float64(rd*fontSize)); err != nil {
		return nil, err
	}
	fontColor := color.RGBA{R: 255, G: 255, B: 255, A: 0xff}
	dc.SetColor(fontColor)
	sw, sh := dc.MeasureString(Text)

	dc.DrawString(Text, float64(W)-sw-float64(rd*fontSize), float64(H)-sh)
	outImg := dc.Image()

	s.Img = outImg

	return s, nil

}
func (s *imgFile) ImgWater(logoPath string, size uint) (*imgFile, error) {
	logoFile, err := os.Open(logoPath)
	if err != nil {
		return nil, err
	}
	defer logoFile.Close()
	logoImg, _, _ := image.Decode(logoFile)
	outImg := image.NewRGBA(s.Img.Bounds())
	resizedLogo := resize.Resize(size, 0, logoImg, resize.Lanczos3)
	offsetX := (outImg.Bounds().Max.X - resizedLogo.Bounds().Max.X) / 2
	offsetY := (outImg.Bounds().Max.Y - resizedLogo.Bounds().Max.Y) / 2
	// Draw the logo onto the QR code image
	draw.Draw(outImg, image.Rect(offsetX, offsetY, offsetX+resizedLogo.Bounds().Max.X, offsetY+resizedLogo.Bounds().Max.Y), resizedLogo, image.Point{}, draw.Over)

	// Draw the QR code onto the output image
	draw.Draw(outImg, outImg.Bounds(), s.Img, image.Point{}, draw.Over)

	s.Img = outImg

	return s, nil

}
