package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func BytesToPNG(rawBytes []byte, banner_id string){
	img, _, _ := image.Decode(bytes.NewReader(rawBytes))
	out, err := os.Create(fmt.Sprintf("./%s.jpeg", banner_id))
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}

func BytesToJPG(rawBytes []byte, banner_id string){
	img, _, _ := image.Decode(bytes.NewReader(rawBytes))
	out, err := os.Create(fmt.Sprintf("./%s.jpeg", banner_id))
	if err != nil {
		log.Fatal(err)
	}


	var opts jpeg.Options
	opts.Quality = 1
	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Fatal(err)
	}
}
