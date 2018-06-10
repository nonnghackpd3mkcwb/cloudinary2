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
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

func writeError(w http.ResponseWriter, msg string, httpCode int) {
	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	fmt.Fprintf(w, `{error":%q}`, msg)
}

type CloudinaryPanic struct {
	error    string
	httpCode uint
}

func triggerCloudinaryPanic(error string, httpCode uint) {
	panic(CloudinaryPanic{
		error:    error,
		httpCode: httpCode,
	})
}

func getThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	defer (func() {
		if r := recover(); r != nil {
			log.Println("ERROR", r)
			cdnryRec, ok := r.(CloudinaryPanic)
			if ok == true {
				writeError(w, fmt.Sprintf("%v", cdnryRec.error), int(cdnryRec.httpCode))
			} else {
				writeError(w, fmt.Sprintf("%v", r), http.StatusInternalServerError)
			}
		}
	})()

	log.Println("Handling a new request:")
	queryValues := r.URL.Query()
	url := queryValues.Get("url")
	width64, widthOk := strconv.ParseInt(queryValues.Get("width"), 10, 32)
	height64, heightOk := strconv.ParseInt(queryValues.Get("height"), 10, 32)

	if len(url) == 0 {
		triggerCloudinaryPanic("missing/bad url parameter", http.StatusBadRequest)
	}

	if widthOk != nil {
		triggerCloudinaryPanic("missing/bad width parameter", http.StatusBadRequest)
	}

	if heightOk != nil {
		triggerCloudinaryPanic("missing/bad height parameter", http.StatusBadRequest)
	}

	if width64 <= 0 || height64 <= 0 {
		triggerCloudinaryPanic("width and height can't be zero or negative", http.StatusBadRequest)
	}

	width, height := int(width64), int(height64)

	log.Println("parameters are (url, width,height) => ", url, width, height)

	data, err := downloadImageData(url)

	if err == nil && http.DetectContentType(data) != "image/jpeg" {
		triggerCloudinaryPanic("image type unsupposrted", http.StatusUnsupportedMediaType)
	}

	if err != nil {
		triggerCloudinaryPanic(
			fmt.Sprintf("unable to fetch image from url: %v", err),
			http.StatusBadRequest)
	}

	//save_data_to_disk(data, "/data/00)original.jpg")

	// Decoding data into an image
	_image, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		triggerCloudinaryPanic("unable to parse image", http.StatusUnsupportedMediaType)
	}

	// Decode image config
	config, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		triggerCloudinaryPanic("unable to parse image configuration", http.StatusUnsupportedMediaType)
	}

	newImage := _image
	trimmedWidth, trimmedHeight := config.Width, config.Height

	// In case the post dimentions are smaller than the original we need to resize the image
	if width < config.Width ||
		height < config.Height {

		ratio := float64(config.Width) / float64(config.Height)

		trimmedWidth, trimmedHeight = int(math.Min(float64(width), float64(config.Width))),
			int(math.Min(float64(height), float64(config.Height)))

		// apply the ratio restoration for width (maybe override in the height restoration)
		if trimmedWidth < config.Width {
			trimmedHeight = int(math.Min(float64(trimmedHeight), float64(trimmedWidth)/ratio))
		}

		// apply the ratio restoration for height
		if trimmedHeight < config.Height {
			trimmedWidth = int(math.Min(float64(trimmedWidth), float64(trimmedHeight)*ratio))
		}

		if trimmedWidth == 0 || trimmedHeight == 0 {
			triggerCloudinaryPanic("invalid image dimentions", http.StatusInternalServerError)
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
		triggerCloudinaryPanic("invalid image dimentions", http.StatusInternalServerError)
	}

	finalBytes := buf.Bytes()

	// Encode the image data and write as an http response
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(finalBytes)))
	if _, err := w.Write(finalBytes); err != nil {
		triggerCloudinaryPanic("unable to reconstruct image", http.StatusInternalServerError)
	}
	log.Println("Image processing complete.")
}

// http://localhost:80/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=50&height=40

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // TBD
	}

	http.HandleFunc("/thumbnail", getThumbnailHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
