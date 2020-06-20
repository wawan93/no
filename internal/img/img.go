package img

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func Generate(url string) (io.Reader, error) {
	// read uploaded photo
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	mark, err := os.Open("marks/box.png")
	if err != nil {
		return nil, err
	}

	// generate new picture
	buf, err := generate(response.Body, mark)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func generate(background, watermark io.Reader) (io.Reader, error) {
	first, err := jpeg.Decode(background)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err)
	}

	mark, err := png.Decode(watermark)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err)
	}

	b := first.Bounds()

	mark = resizeMark(mark, b)

	offsetY := b.Dy() - mark.Bounds().Dy()
	offset := image.Pt(0, offsetY)

	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, first, image.Point{}, draw.Src)
	draw.Draw(image3, mark.Bounds().Add(offset), mark, image.Point{}, draw.Over)

	var buf bytes.Buffer
	jpeg.Encode(&buf, image3, &jpeg.Options{jpeg.DefaultQuality})

	return &buf, nil
}

func resizeMark(img image.Image, b image.Rectangle) image.Image {
	var w, h uint

	if b.Dx() > b.Dy() {
		h = uint(b.Dy() / 2)
		w = h
	}

	if b.Dx() <= b.Dy() {
		w = uint(b.Dx() / 2)
		h = w
	}

	return resize.Resize(w, h, img, resize.Lanczos2)
}
