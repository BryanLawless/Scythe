package utility

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func Prompt(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintf(color.Output, "%s%s[INPUT]%s %s: ", BlackBg, Purple, ColorReset, prompt)

	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		return ""
	}

	return strings.Trim(scanner.Text(), "")
}

func ExitWrapper(err error) {
	if err != nil && err == terminal.InterruptErr {
		os.Exit(0)
	}
}

func StartActionPrompt() string {
	option := ""
	prompt := &survey.Select{
		Message: "Select an action:",
		Options: []string{"Browse", "Download", "Exit"},
	}

	ExitWrapper(survey.AskOne(prompt, &option))

	return option
}

func MediaCategoryPrompt() string {
	media := ""
	prompt := &survey.Select{
		Message: "Select a category:",
		Options: []string{"Videos", "Pictures"},
	}

	ExitWrapper(survey.AskOne(prompt, &media))

	return media
}

func VideoProviderPrompt(providers []string) string {
	provider := ""
	prompt := &survey.Select{
		Message: "Select a provider:",
		Options: providers,
	}

	ExitWrapper(survey.AskOne(prompt, &provider))

	return provider
}

func ChooseResultsPrompt(results []map[string]string, headers []string, limits map[string]int) []string {
	tableHeaders := AddToStart(headers, "#")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader(tableHeaders)

	for index, result := range results {
		for limitKey, limitVal := range limits {
			result[limitKey] = CutLongText(result[limitKey], limitVal)
		}

		data := SortedMapValuesByKeys(headers, result)
		data = AddToStart(data, strconv.Itoa(index+1))

		table.Append(data)
	}

	table.Render()

	selectedVideos := []string{}
	choices := SelectResults(len(results))
	for _, choice := range choices {
		selectedVideos = append(selectedVideos, results[choice-1]["ID"])
	}

	return selectedVideos
}

func SelectResults(results int) []int {
	selected := ""
	choices := []int{}
	question := &survey.Question{
		Prompt: &survey.Input{
			Message: "Select resources (#):",
			Help:    "Enter the number corresponding to the resource you want to download. You can use multiple numbers separated by commas.",
		},
		Validate: func(val interface{}) error {
			numbers := strings.Split(val.(string), ",")

			for _, num := range numbers {
				num, err := strconv.Atoi(strings.TrimSpace(num))

				if err != nil {
					return fmt.Errorf("input must be numbers")
				} else if num < 1 || num > results {
					return fmt.Errorf("input range exceeded (1-%d)", results)
				}

			}

			return nil
		},
	}

	ExitWrapper(survey.Ask([]*survey.Question{question}, &selected))

	for _, num := range strings.Split(selected, ",") {
		num, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			continue
		}

		choices = append(choices, num)
	}

	return choices
}

func UrlSafeSearchPrompt() string {
	return strings.ReplaceAll(Prompt("Enter a search term"), " ", "+")
}
