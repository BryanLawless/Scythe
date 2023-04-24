package youtube

import (
	"Scythe/core/common"
	"Scythe/core/providers"
	"Scythe/core/utility"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kkdai/youtube/v2"
)

var videoUrl string = "https://www.youtube.com/watch?v=%s"
var searchUrl string = "https://www.youtube.com/results?search_query=%s&sp=EgIQAQ%253D%253D"

type youtubeProvider struct {
	Client *youtube.Client
}

func init() {
	providers.Register("youtube", "video", New())
}

func New() providers.Provider {
	return &youtubeProvider{
		Client: new(youtube.Client),
	}
}

func (y *youtubeProvider) Start(ctx context.Context) ([]common.Media, error) {
	urlSafeSearchTerm := utility.UrlSafeSearchPrompt()

	videos, err := y.scrapeYoutube(ctx, urlSafeSearchTerm)
	if err != nil {
		return nil, fmt.Errorf("error scraping youtube: %w", err)
	}

	selectedResults := utility.ChooseResultsPrompt(videos, []string{
		"ID",
		"Title",
		"Channel",
		"Views",
		"Length",
		"Uploaded",
	}, map[string]int{"Title": 40, "Channel": 15})

	if len(selectedResults) == 0 {
		return nil, nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Gathering selected resources..."
	s.Start()

	youtubeMedia := []common.Media{}
	for _, result := range selectedResults {
		videoLink := fmt.Sprintf(videoUrl, result)
		video, err := y.Client.GetVideo(videoLink)
		if err != nil {
			return nil, fmt.Errorf("error getting video: %w", err)
		}

		formats := video.Formats.WithAudioChannels()
		extension := utility.GetExtensionFromMime(formats[0].MimeType)
		sourceLink, err := y.Client.GetStreamURL(video, &formats[0])
		if err != nil {
			return nil, fmt.Errorf("error getting video link: %w", err)
		}

		youtubeMedia = append(youtubeMedia, common.Media{
			URL:       videoLink,
			Name:      video.Title,
			Provider:  "youtube",
			Category:  "video",
			Extension: extension,
			SourceURL: sourceLink,
		})
	}

	s.Stop()

	utility.WriteToConsole("Added selected links to download list", "success")

	return youtubeMedia, nil
}

func (y *youtubeProvider) scrapeYoutube(ctx context.Context, urlSafeSearchTerm string) ([]map[string]string, error) {
	videos := []map[string]string{}

	request, body, err := common.MakeRequest(ctx, common.Request{
		URL:       fmt.Sprintf(searchUrl, urlSafeSearchTerm),
		Method:    "GET",
		ParseBody: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if request.Response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get search results")
	}

	searchResponse := string(body)
	searchResponse = strings.Join(strings.Split(searchResponse, "\n"), "")
	searchResponse = strings.Split(searchResponse, "var ytInitialData")[1]

	chunk := strings.Split(searchResponse, "=")[1:]
	rawData := strings.Split(strings.Join(chunk, "="), ";</script>")[0]

	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(rawData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %s", err)
	}

	found := data["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"]

	if val, ok := found.(map[string]interface{})["sectionListRenderer"]; ok {
		for _, v := range val.(map[string]interface{})["contents"].([]interface{}) {
			if val, ok := v.(map[string]interface{})["itemSectionRenderer"]; ok {
				for _, v := range val.(map[string]interface{})["contents"].([]interface{}) {
					if val, ok := v.(map[string]interface{})["videoRenderer"]; ok {

						uploaded := "Unknown"
						if val.(map[string]interface{})["publishedTimeText"] != nil {
							uploaded = val.(map[string]interface{})["publishedTimeText"].(map[string]interface{})["simpleText"].(string)
						}

						videos = append(videos, map[string]string{
							"ID":       val.(map[string]interface{})["videoId"].(string),
							"Title":    val.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string),
							"Channel":  val.(map[string]interface{})["ownerText"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string),
							"Views":    val.(map[string]interface{})["viewCountText"].(map[string]interface{})["simpleText"].(string),
							"Length":   val.(map[string]interface{})["lengthText"].(map[string]interface{})["simpleText"].(string),
							"Uploaded": uploaded,
						})
					}
				}
			}
		}
	}

	return videos, nil
}
