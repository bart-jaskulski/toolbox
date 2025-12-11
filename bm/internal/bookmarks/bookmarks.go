package bookmarks

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/bart-jaskulski/toolbox/bm/internal/utils"
)

var bookmarkFileMutex = &sync.Mutex{}

func AddBookmark(urlStr string) error {
	if !utils.IsValidURL(urlStr) {
		return fmt.Errorf("invalid URL: %s", urlStr)
	}

	// Fetch the title from the URL
	title, err := utils.FetchTitleFromURL(urlStr)
	if err != nil {
		return fmt.Errorf("error fetching title from URL: %v", err)
	}

	bookmarkFileMutex.Lock()
	defer bookmarkFileMutex.Unlock()

	file, err := os.OpenFile(utils.GetBookmarksFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	bookmarkEntry := fmt.Sprintf("%s \"%s\"\n", urlStr, title)
	_, err = file.WriteString(bookmarkEntry)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func FindBookmarks(searchPattern string) ([]string, error) {
	var matchingBookmarks []string
	file, err := os.Open(utils.GetBookmarksFilePath())
	if err != nil {
		return nil, fmt.Errorf("error opening .bookmark file: %v", err)
	}
	defer file.Close()

	pattern, err := regexp.Compile(`(?i)` + searchPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if pattern.MatchString(scanner.Text()) {
			matchingBookmarks = append(matchingBookmarks, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning .bookmark file: %v", err)
	}

	return matchingBookmarks, nil
}

func SearchAndChooseBookmark(searchTerm string) (string, error) {
	matches, err := FindBookmarks(searchTerm)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no bookmarks found")
	} else if len(matches) > 1 {
		chosenBookmark := utils.ChooseBookmark(matches)
		if chosenBookmark == "" {
			return "", fmt.Errorf("no bookmark chosen")
		}
		return chosenBookmark, nil
	}
	return matches[0], nil
}
