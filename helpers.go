package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func save_data_to_disk(data []byte, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return out.Close()
}

func downloadImageData(url string) ([]byte, error) {
	response, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	if response.StatusCode != 200 {
		return nil, errors.New(
			fmt.Sprintf("bad response from url, status code %v", response.StatusCode))
	}

	defer response.Body.Close()
	rawData, e := ioutil.ReadAll(response.Body)

	if http.DetectContentType(rawData) != "image/jpeg" {
		return nil, errors.New("bad image format")
	}

	if e != nil {
		return nil, e
	}

	return rawData, nil
}
