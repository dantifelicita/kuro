package main

import (
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/image/webp"
)

func main() {
	file, err := readFile("example.txt")
	if err != nil {
		panic(err)
	}

	urlList := strings.Split(file, "\n")

	for _, link := range urlList {
		path := strings.Split(link, "/")
		filePath, err := url.QueryUnescape(path[len(path)-1])
		if err != nil {
			panic(err)
		}

		err = downloadFile(filePath, link)
		if err != nil {
			panic(err)
		}
		fmt.Println("Downloaded: " + link)

		fileName := strings.Split(filePath, ".")

		if ext := fileName[len(fileName)-1]; ext == "webp" {
			err = convertWebp(fileName[0], ext)
			if err != nil {
				panic(err)
			}
			fmt.Println("Converted: " + link)
		}
	}
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func convertWebp(fileName, ext string) error {
	file, err := os.Open(fileName + "." + ext)
	if err != nil {
		return err
	}

	img, err := webp.Decode(file)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName + ".png")
	if err != nil {
		return err
	}

	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(name, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
