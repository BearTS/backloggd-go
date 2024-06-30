package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func (sdk *BackloggdSDK) GetPlaythroughDetails(id string) (PlaythroughDetails, error) {
	var playthroughDetails PlaythroughDetails
	req, err := http.NewRequest("GET", fmt.Sprintf(playthroughDetailsURL, id), nil)
	if err != nil {
		return playthroughDetails, err
	}

	req.Header.Set("Referer", baseURL)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.backloggd.com/u/beartstest/logs/lego-batman-the-videogame")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("X-CSRF-Token", sdk.csrfToken)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return playthroughDetails, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&playthroughDetails)
	if err != nil {
		return playthroughDetails, err
	}

	return playthroughDetails, nil
}

func (sdk *BackloggdSDK) GetPlaythroughIds(gameslug string) ([]Logs, string, error) { // returns logs, csrfToken, error
	var logs []Logs
	req, err := http.NewRequest("GET", fmt.Sprintf(logsURL, sdk.Username, gameslug), nil)
	if err != nil {
		return logs, "", err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return logs, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return logs, "", err
	}

	doc.Find(".delete-log").Each(func(i int, s *goquery.Selection) {
		logs = append(logs, Logs{
			ID: s.AttrOr("playthrough_id", ""),
		})
	})

	// get csrf token
	csrfToken, _ := doc.Find("meta[name='csrf-token']").Attr("content")

	return logs, csrfToken, nil
}

type Logs struct {
	ID string `json:"id"`
}

type PlaythroughDetails struct {
	GameLog struct {
		ID                    *int        `json:"id"` // log id
		Status                *string     `json:"status"`
		Rating                *int        `json:"rating"`
		IsPlay                *bool       `json:"is_play"`
		IsPlaying             *bool       `json:"is_playing"`
		IsBacklog             *bool       `json:"is_backlog"`
		IsWishlist            *bool       `json:"is_wishlist"`
		TimeSource            *int        `json:"time_source"`
		TotalHours            *float64    `json:"total_hours"`
		TotalMinutes          *int        `json:"total_minutes"`
		LibraryEntries        interface{} `json:"library_entries"`
		TotalHoursFormatted   *string     `json:"total_hours_formatted"`
		TotalMinutesFormatted *string     `json:"total_minutes_formatted"`
	} `json:"game_log"`
	Playthrough struct {
		ID             *int          `json:"id"`
		CreatedAt      *string       `json:"created_at"`
		Rating         *int          `json:"rating"`
		Review         *string       `json:"review"`
		ReviewSpoilers *bool         `json:"review_spoilers"`
		MediumID       interface{}   `json:"medium_id"`
		Platform       *int          `json:"platform"`
		PlayedPlatform interface{}   `json:"played_platform"`
		Hours          *float64      `json:"hours"`
		Minutes        *int          `json:"minutes"`
		IsMaster       *bool         `json:"is_master"`
		IsReplay       *bool         `json:"is_replay"`
		Title          *string       `json:"title"`
		PlayDates      []interface{} `json:"play_dates"`
		StartDate      interface{}   `json:"start_date"`
		FinishDate     interface{}   `json:"finish_date"`
	} `json:"playthrough"`
}
