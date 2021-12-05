package main

import (
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (g *generator) fontToImage(filePathName string, font string) (err error) {

	img := image.NewRGBA(image.Rect(0, 0, 640, 120))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	g.create(img, 20, 30, font)

	f, err := os.Create(filePathName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	return nil
}

func (g *generator) create(img *image.RGBA, x, y int, fontWord string) (err error) {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile("./luxisr.ttf")
	if err != nil {
		log.Println(err)
		return
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	d := &font.Drawer{
		Dst: img,
		Src: image.Black,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    24,
			Hinting: font.HintingNone,
		}),
	}

	d.Dot = fixed.Point26_6{
		X: fixed.Int26_6(x * 50),
		Y: fixed.Int26_6(y * 140),
	}
	d.DrawString(fontWord)

	return nil
}
