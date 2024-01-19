package utility

import (
	"Scythe/core/common"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func NotPrev(current, previous string) bool {
	return current != previous && current != ""
}

func MatchOne(pattern string, text string) []string {
	re := regexp.MustCompile(pattern)
	value := re.FindStringSubmatch(text)

	if value != nil {
		return value
	}

	return nil
}

func Contains(elements []string, value string) bool {
	for _, search := range elements {
		if value == search {
			return true
		}
	}

	return false
}

func RemoveValFromSlice(slice []string, value string) []string {
	for index, search := range slice {
		if value == search {
			return RemoveIndexFromSlice(slice, index)
		}
	}

	return slice
}

func RemoveIndexFromSlice(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	slice[len(slice)-1] = ""

	return slice[:len(slice)-1]
}

func RemoveDuplicates(str_slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range str_slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}

	return list
}

func PartialDirectory(targetDirectory, filename string, processes int) string {
	if targetDirectory == "" {
		return fmt.Sprintf("_%s.%d", filename, processes)
	}

	return filepath.Join(targetDirectory, fmt.Sprintf("_%s.%d", filename, processes))
}

func DirectorySize(directory string) (int64, error) {
	var size int64
	err := filepath.Walk(directory, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}

		return err
	})

	return size, err
}

func SortedMapValuesByKeys(keys []string, m map[string]string) []string {
	sorted := make([]string, len(keys))

	for index, key := range keys {
		if value, ok := m[key]; ok {
			sorted[index] = value
		}
	}

	return sorted
}

func AddToStart(slice []string, value string) []string {
	return append([]string{value}, slice...)
}

func CutLongText(text string, max int) string {
	clean := strings.TrimSpace(text)
	if len(clean) > max {
		return clean[0:max-3] + "..."
	}

	return text
}

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	random := make([]rune, length)
	for i := range random {
		random[i] = letters[rand.Intn(len(letters))]
	}

	return string(random)
}

func FilenameSafe(title string) string {
	title = strings.TrimSpace(title)
	filter := []string{" ", ":", "/", "\\", "?", "*", "\"", "<", ">", "|", "(", ")", "'", "!", "."}

	for _, item := range filter {
		switch item {
		case " ":
			title = strings.ReplaceAll(title, item, "_")
		case ".":
			title = strings.ReplaceAll(title, item, "")
		default:
			title = strings.ReplaceAll(title, item, "")
		}
	}

	return title
}

func MediaToResults(media []common.Media) []string {
	results := make([]string, len(media))

	for index, item := range media {
		results[index] = fmt.Sprintf("%d-[%s](%s)", index, item.Name, item.URL)
	}

	return results
}

func RemovePeriods(str string) string {
	return strings.ReplaceAll(str, ".", "")
}

func GetFileNameFromUrl(rawUrl string) (string, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	path := url.Path
	segments := strings.Split(path, "/")

	if len(segments) == 0 {
		return "", errors.New("invalid URL")
	}

	return segments[len(segments)-1], nil
}

func GetFileExtensionFromMime(contentType string) string {
	parts := strings.Split(contentType, ";")
	mediaType := strings.TrimSpace(parts[0])

	typeParts := strings.Split(mediaType, "/")

	extension := typeParts[len(typeParts)-1]

	return extension
}

func GetCategoryFromMimeType(mimeType string) string {
	extension := GetFileExtensionFromMime(mimeType)

	switch extension {
	case "mp4", "webm", "mkv", "flv", "avi", "mov", "wmv", "mpg", "mpeg":
		return "video"
	case "mp3", "wav", "ogg", "flac", "m4a", "wma", "aac":
		return "audio"
	case "jpg", "jpeg", "png", "gif", "bmp", "svg", "webp":
		return "image"
	case "pdf", "epub", "mobi", "azw", "azw3", "djvu", "fb2", "prc", "doc", "docx", "txt":
		return "document"
	default:
		return "unknown"
	}
}
