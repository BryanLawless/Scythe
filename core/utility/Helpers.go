package utility

import (
	"Scythe/core/common"
	"fmt"
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

func GetExtensionFromMime(mime string) string {
	types := MatchOne(`(?i)(\w+)/(\w+);`, mime)
	if types == nil {
		return ""
	}

	return types[2]
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
