package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	totalCount, successCount int
	isPage                   *bool
	folderPath               = "images"
	urlRegex                 = regexp.MustCompile(`(http(s?):)([/|.|\w|\s|-])*\.(?:jpg|jpeg|webp|png)`)
)

func main() {
	timeStart := time.Now()
	os.MkdirAll(folderPath, os.ModePerm)

	// Read file
	file, err := readFile("example.txt")
	if err != nil {
		panic(err)
	}

	// Get list of URLs
	urlList := strings.Split(file, "\n")

	// Download each URL
	for idx, link := range urlList {
		if isRequestPage() {
			err = readPage(idx, link)
			if err != nil {
				logWithTag("P", idx, err.Error())
			}
		} else {
			err = downloadImage(idx, link)
			if err != nil {
				logWithTag("I", idx, err.Error())
			}
		}
	}

	str := "image"
	if successCount > 1 {
		str = "images"
	}
	fmt.Printf("Finished downloading: %d/%d %s in %v\n", successCount, totalCount, str, time.Now().Sub(timeStart))
}

func downloadImage(idx int, link string) error {
	if !validateLink(link) {
		return nil
	}
	totalCount++

	// Get file name and original extension decoded
	path := strings.Split(link, "/")
	name, err := url.QueryUnescape(path[len(path)-1])
	if err != nil {
		return err
	}

	// Download the file
	err = downloadFile(name, link)
	if err != nil {
		return errors.New("Error downloading: " + link + " err: " + err.Error())
	}
	logWithTag("I", idx, "Downloaded: "+link)

	// Get raw file name
	fileName := strings.Split(name, ".")

	// Convert if in webp format
	if ext := fileName[len(fileName)-1]; ext == "webp" {
		err = convertWebp(fileName[0], ext)
		if err != nil {
			return errors.New("Error converting: " + link + " err: " + err.Error())
		}
		logWithTag("I", idx, "Converted: "+link)
	}

	successCount++
	return nil
}

func readPage(idx int, link string) error {
	logWithTag("P", idx, "Fetching from page: "+link)

	request := NewRequest()
	request.URL = link
	request.Method = "GET"

	response, body, err := request.doRequest()
	if err != nil {
		return err
	}
	defer response.Body.Close()

	imgLinks := urlRegex.FindAllString(string(body), -1)

	for idx, imgLink := range imgLinks {
		err = downloadImage(idx, imgLink)
		if err != nil {
			continue
		}
	}

	return nil
}
