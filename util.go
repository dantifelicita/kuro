package main

import (
	"fmt"
	"os"
	"strings"
)

func logWithTag(prefix string, idx int, msg string) {
	fmt.Printf("[%s%d] %s\n", prefix, idx+1, msg)
}

func validateLink(link string) bool {
	return strings.HasPrefix(link, "http")
}

func getFilePath(name string) string {
	return folderPath + "/" + name
}

func isRequestPage() bool {
	if isPage != nil {
		return *isPage
	}
	p := false

	if len(os.Args) > 1 {
		if len(os.Args[1]) >= 5 {
			if os.Args[1][:5] == "-page" {
				p = true
			}

		}
	}
	isPage = &p

	return *isPage
}
