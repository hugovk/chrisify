package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

var facesDir = flag.String("faces", "faces", "The directory to search for faces.")

func main() {
	flag.Parse()

	var chrisFaces FaceList

	var facesPath string
	var err error

	if *facesDir != "" {
		facesPath, err = filepath.Abs(*facesDir)
		if err != nil {
			panic(err)
		}
	}

	err = chrisFaces.Load(facesPath)
	if err != nil {
		panic(err)
	}
	if len(chrisFaces) == 0 {
		panic("no faces found")
	}

	file := flag.Arg(0)

	baseImage := loadImage(file)

	bounds := baseImage.Bounds()

	canvas := canvasFromImage(baseImage)

	face := imaging.Resize(
		chrisFaces[0],
		bounds.Dx()/3,
		0,
		imaging.Lanczos,
	)
	face_bounds := face.Bounds()
	draw.Draw(
		canvas,
		bounds,
		face,
		bounds.Min.Add(image.Pt(-2*bounds.Max.X/3+face_bounds.Max.X/2, -bounds.Max.Y+int(float64(face_bounds.Max.Y)))),
		draw.Over,
	)

	jpeg.Encode(os.Stdout, canvas, &jpeg.Options{jpeg.DefaultQuality})
}
