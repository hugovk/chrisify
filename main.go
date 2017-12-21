package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/zikes/chrisify/facefinder"

	"github.com/disintegration/imaging"
)

var haarCascade = flag.String("haar", "haarcascade_frontalface_alt.xml", "The location of the Haar Cascade XML configuration to be provided to OpenCV.")
var facesDir = flag.String("faces", "faces", "The directory to search for faces.")
var bodiesDir = flag.String("bodies", "bodies", "The directory to search for bodies.")

func main() {
	flag.Parse()

	var chrisFaces FaceList
	var chrisBodies FaceList

	var facesPath string
	var bodiesPath string
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
	err = chrisBodies.Load(bodiesPath)
	if err != nil {
		panic(err)
	}
	if len(chrisBodies) == 0 {
		panic("no faces found")
	}

	file := flag.Arg(0)

	finder := facefinder.NewFinder(*haarCascade)

	baseImage := loadImage(file)

	faces := finder.Detect(baseImage)

	bounds := baseImage.Bounds()

	canvas := canvasFromImage(baseImage)

	for _, face := range faces {
		rect := rectMargin(30.0, face)

		newFace := chrisFaces.Random()
		if newFace == nil {
			panic("nil face")
		}
		chrisFace := imaging.Fit(newFace, rect.Dx(), rect.Dy(), imaging.Lanczos)

		draw.Draw(
			canvas,
			rect,
			chrisFace,
			bounds.Min,
			draw.Over,
		)
	}

	if len(faces) == 0 {
		body := imaging.Resize(
			chrisBodies.Random(),
			bounds.Dx()/3,
			0,
			imaging.Lanczos,
		)
		body_bounds := body.Bounds()
		draw.Draw(
			canvas,
			bounds,
			body,
			bounds.Min.Add(image.Pt(-2*bounds.Max.X/3+body_bounds.Max.X/2, -bounds.Max.Y+int(float64(body_bounds.Max.Y)))),
			draw.Over,
		)
	}

	jpeg.Encode(os.Stdout, canvas, &jpeg.Options{jpeg.DefaultQuality})
}
