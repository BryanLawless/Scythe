package utility

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/olekukonko/tablewriter"
)

func exitWrapper(err error) {
	if err != nil && err == terminal.InterruptErr {
		os.Exit(0)
	}
}

func SearchPrompt() string {
	text := ""
	prompt := &survey.Question{
		Prompt: &survey.Input{
			Message: "Search terms:",
			Help:    "Enter the search terms you would like to use when searching for links.",
		},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); !ok || len(str) > 64 {
				return fmt.Errorf("search terms must be less than 64 characters long")
			}

			return nil
		},
	}

	exitWrapper(survey.Ask([]*survey.Question{prompt}, &text))

	return text
}

func UrlSafeSearchPrompt() string {
	return strings.ReplaceAll(SearchPrompt(), " ", "+")
}

func StartActionPrompt() string {
	option := ""
	prompt := &survey.Select{
		Message: "Select an action:",
		Options: []string{"Browse", "Download", "Exit"},
	}

	exitWrapper(survey.AskOne(prompt, &option))

	return option
}

func MediaCategoryPrompt() string {
	media := ""
	prompt := &survey.Select{
		Message: "Select a category:",
		Options: []string{"Videos"},
	}

	exitWrapper(survey.AskOne(prompt, &media))

	return media
}

func VideoProviderPrompt(providers []string) string {
	provider := ""
	prompt := &survey.Select{
		Message: "Select a provider:",
		Options: providers,
	}

	exitWrapper(survey.AskOne(prompt, &provider))

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
			Help:    "Enter the number related to the resource you want to download. You can use multiple numbers separated by commas. Type c to cancel.",
		},
		Validate: func(val interface{}) error {
			if val.(string) == "c" {
				return nil
			}

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

	exitWrapper(survey.Ask([]*survey.Question{question}, &selected))

	if selected == "c" {
		return choices
	}

	for _, num := range strings.Split(selected, ",") {
		num, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			continue
		}

		choices = append(choices, num)
	}

	return choices
}

func ConfirmSelectDownloadsPrompt(options []string) []string {
	mediaSelected := []string{}
	question := survey.Question{
		Prompt: &survey.MultiSelect{
			Message: "Confirm the content you want to download:",
			Options: options,
		},
		Validate: func(val interface{}) error {
			if len(val.([]survey.OptionAnswer)) == 0 {
				return fmt.Errorf("you must select at least one option")
			}

			return nil
		},
	}

	exitWrapper(survey.Ask([]*survey.Question{&question}, &mediaSelected))

	return mediaSelected
}
