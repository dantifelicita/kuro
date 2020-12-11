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
	"time"

	"golang.org/x/image/webp"
)

func main() {
	timeStart := time.Now()

	// Read file
	file, err := readFile("example.txt")
	if err != nil {
		panic(err)
	}

	// Get list of URLs
	urlList := strings.Split(file, "\n")

	var totalCount, successCount int
	fmt.Println("Starting...")

	// Download each URL
	for _, link := range urlList {
		totalCount++

		// Get file name and original extension decoded
		path := strings.Split(link, "/")
		name, err := url.QueryUnescape(path[len(path)-1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Download the file
		err = downloadFile(name, link)
		if err != nil {
			fmt.Println("Error downloading: " + link + " err: " + err.Error())
			continue
		}
		fmt.Println("Downloaded: " + link)

		// Get raw file name
		fileName := strings.Split(name, ".")

		// Convert if in webp format
		if ext := fileName[len(fileName)-1]; ext == "webp" {
			err = convertWebp(fileName[0], ext)
			if err != nil {
				fmt.Println("Error converting: " + link + " err: " + err.Error())
				continue
			}
			fmt.Println("Converted: " + link)
		}

		successCount++
	}

	str := "file"
	if successCount > 1 {
		str = "files"
	}
	fmt.Printf("Finish downloading: %d/%d %s in %v\n", successCount, totalCount, str, time.Now().Sub(timeStart))
}

// readFile will read from a text file to string
func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

// convertWebp will convert a webp file into png
func convertWebp(fileName, ext string) error {
	// Open the file
	file, err := os.Open(fileName + "." + ext)
	if err != nil {
		return err
	}

	// Decode webp file
	img, err := webp.Decode(file)
	if err != nil {
		return err
	}

	// Create new png file
	f, err := os.Create(fileName + ".png")
	if err != nil {
		return err
	}

	// Encode to png
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
