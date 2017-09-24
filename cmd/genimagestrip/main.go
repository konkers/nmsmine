package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"

	"github.com/konkers/nmsmine"
	"github.com/nfnt/resize"
)

var jsonFile = flag.String("data", "", "Json file containing item data")
var assetDir = flag.String("assets", "", "Path to a directory containing PNG assets.")
var imgFilename = flag.String("image", "", "File to write output PNG.")
var mapFilename = flag.String("map", "", "File to write output JSON map.")
var width = flag.Int("width", 32, "Json file containing item data")

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	flag.Parse()

	if *jsonFile == "" {
		fmt.Fprintf(os.Stderr, "-data must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if *assetDir == "" {
		fmt.Fprintf(os.Stderr, "-assets must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if *imgFilename == "" {
		fmt.Fprintf(os.Stderr, "-image must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if *mapFilename == "" {
		fmt.Fprintf(os.Stderr, "-map must be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	db, err := nmsmine.LoadItemDb(*jsonFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	imageSet := map[string]bool{}
	for _, item := range db {
		imageSet[item.Icon] = true
	}

	images := []string{}
	for image, _ := range imageSet {
		images = append(images, image)
	}

	sort.Strings(images)

	stride := int(math.Floor(math.Sqrt(float64(len(images)))))

	yDim := int(math.Ceil(float64(len(images))/float64(stride))) * *width
	xDim := stride * *width

	strip := image.NewRGBA(image.Rect(0, 0, xDim, yDim))
	imageMap := map[string]Point{}

	for index, imgName := range images {
		srcFilename := path.Join(*assetDir, imgName+".png")
		srcFile, err := os.Open(srcFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't open %s: %s\n", srcFilename, err.Error())
			continue
		}
		defer srcFile.Close()

		src, err := png.Decode(srcFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't decode %s: %s\n", srcFilename, err.Error())
			continue
		}

		srcResized := resize.Resize(uint(*width), 0, src, resize.Bilinear)

		x := index % stride
		y := (index - x) / stride
		xPos := x * *width
		yPos := y * *width

		imageMap[imgName] = Point{X: x, Y: y}
		destRect := image.Rect(xPos, yPos, xPos+*width, yPos+*width)
		draw.Draw(strip, destRect, srcResized, image.Pt(0, 0), draw.Src)
	}

	outFile, err := os.OpenFile(*imgFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s: %s\n", *imgFilename, err.Error())
		os.Exit(1)
	}
	defer outFile.Close()
	png.Encode(outFile, strip)

	b, err := json.MarshalIndent(imageMap, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %s\n", err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(*mapFilename, b, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to %s: %s\n", *mapFilename, err.Error())
		os.Exit(1)
	}
}
