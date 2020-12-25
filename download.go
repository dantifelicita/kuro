package main

import (
	"errors"
	"net/url"
	"strings"
)

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

func getImageLinks(link string) ([]string, []string, error) {
	request := NewRequest()
	request.URL = link
	request.Method = "GET"

	response, body, err := request.doRequest()
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	arr1 := urlRegex.FindAllString(string(body), -1)
	arr2 := rawUrlRegex.FindAllString(string(body), -1)

	return arr1, arr2, nil
}

func readPage(idx int, link string) error {
	logWithTag("P", idx, "Fetching from page: "+link)

	imgLinks, rawImgLinks, err := getImageLinks(link)
	if err != nil {
		return err
	}

	for idx, imgLink := range imgLinks {
		err = downloadImage(idx, imgLink)
		if err != nil {
			logWithTag("P", idx, "Error downloading: "+imgLink+" "+err.Error())
		}
	}

	for idx, imgLink := range rawImgLinks {
		err = downloadImage(idx, "https://"+imgLink)
		if err != nil {
			logWithTag("P", idx, "Error downloading: https://"+imgLink+" "+err.Error())
		}
	}

	return nil
}

func fetchLinks(link string) error {
	imgLinks, rawImgLinks, err := getImageLinks(link)
	if err != nil {
		return err
	}

	text := strings.Join(imgLinks, "\n")

	text2 := strings.Join(rawImgLinks, "\nhttps://")
	text2 = "https://" + text2

	outLog.WriteString(text + text2)

	return nil
}
