package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/nfnt/resize"
)

func writeError500(w http.ResponseWriter, msg string) {
	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{error":%q}`, msg)
}

func getThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	defer (func() {
		if r := recover(); r != nil {
			log.Println("ERROR", r)
			writeError500(w, fmt.Sprintf("%v", r))
		}
	})()

	log.Println("Handling a new request:")
	queryValues := r.URL.Query()
	url := queryValues.Get("url")
	width64, widthOk := strconv.ParseInt(queryValues.Get("width"), 10, 32)
	height64, heightOk := strconv.ParseInt(queryValues.Get("height"), 10, 32)

	if len(url) == 0 {
		panic("missing/bad url parameter")
	}

	if widthOk != nil {
		panic("missing/bad width parameter")
	}

	if heightOk != nil {
		panic("missing/bad height parameter")
	}

	if width64 <= 0 || height64 <= 0 {
		panic("width and height can't be zero or negative")
	}

	width, height := int(width64), int(height64)

	log.Println("parameters are (url, width,height) => ", url, width, height)

	data, err := downloadImageData(url)
	if err != nil {
		panic(fmt.Sprintf("unable to fetch image from url: %v", err))
	}

	//save_data_to_disk(data, "/data/00)original.jpg")

	// Decoding data into an image
	_image, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic("unable to parse image")
	}

	// Decode image config
	config, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		panic("unable to parse image config")
	}

	newImage := _image
	trimmedWidth, trimmedHeight := config.Width, config.Height

	// In case the post dimentions are smaller than the original we need to resize the image
	if width < config.Width ||
		height < config.Height {

		trimmedWidth, trimmedHeight = int(math.Min(float64(width), float64(config.Width))),
			int(math.Min(float64(height), float64(config.Height)))

		if trimmedWidth == 0 || trimmedHeight == 0 {
			panic("invalid image dimentions")
		}

		log.Println("image is resized to", trimmedWidth, trimmedHeight)

		newImage = resize.Resize(
			uint(trimmedWidth),
			uint(trimmedHeight),
			_image,
			resize.Lanczos3)
	}

	relX, relY := (width-trimmedWidth)/2, (height-trimmedHeight)/2
	buf := new(bytes.Buffer)
	if relX != 0 || relY != 0 {
		// if we a relative position to place our image
		// thus, creating horizontal or vertical padding
		newImageBuffer := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(newImageBuffer, newImageBuffer.Bounds(), newImage, image.Point{-relX, -relY}, draw.Src)
		err = jpeg.Encode(buf, newImageBuffer, nil)
	} else {
		err = jpeg.Encode(buf, newImage, nil)
	}

	if err != nil {
		panic("unable to encode padded image after processing")
	}

	finalBytes := buf.Bytes()

	// Encode the image data and write as an http response
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(finalBytes)))
	if _, err := w.Write(finalBytes); err != nil {
		panic("unable to reconstruct image")
	}
	log.Println("Image processing complete.")
}

// http://localhost:80/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=50&height=40

func main() {
	http.HandleFunc("/thumbnail", getThumbnailHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
