package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	urlRegex    = regexp.MustCompile(`(http(s?):)([/|.|\w|\s|%|-])*\.(?:jpg|jpeg|webp|png)`)
	rawUrlRegex = regexp.MustCompile(`(www)([/|.|\w|\s|-])*\.(?:jpg|jpeg|webp|png)`)
)

func logWithTag(prefix string, idx int, msg string) {
	text := fmt.Sprintf("[%s%d] %s\n", prefix, idx+1, msg)
	fmt.Print(text)
	fileLog.WriteString(text)
}

func validateLink(link string) bool {
	return urlRegex.MatchString(link)
}

func getFilePath(name string) string {
	return folderPath + "/" + name
}

func getMode() string {
	if mode != nil {
		return *mode
	}

	var m string
	if len(os.Args) > 1 {
		m = os.Args[1][1:]
	}
	mode = &m

	return *mode
}

func getPageQuery() (bool, int, int) {
	if pageFrom != nil && pageUntil != nil {
		return *usePageQuery, *pageFrom, *pageUntil
	}

	var (
		exist  bool
		p1, p2 int
	)

	if len(os.Args) > 2 {
		p1, _ = strconv.Atoi(os.Args[2])
		exist = true
	}

	if len(os.Args) > 3 {
		p2, _ = strconv.Atoi(os.Args[3])
	}

	if p2 == 0 {
		p2 = p1
	}

	pageFrom = &p1
	pageUntil = &p2
	usePageQuery = &exist

	return *usePageQuery, *pageFrom, *pageUntil
}
