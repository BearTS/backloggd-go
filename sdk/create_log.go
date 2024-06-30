package sdk

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// game_id: 2738
// playthroughs[0][id]: -1
// playthroughs[0][title]: Log
// playthroughs[0][rating]: 8
// playthroughs[0][review]: test good
// playthroughs[0][review_spoilers]: false
// playthroughs[0][platform]: 9
// playthroughs[0][hours]: 23
// playthroughs[0][minutes]: 2
// playthroughs[0][is_master]: false
// playthroughs[0][is_replay]: true
// playthroughs[0][start_date]:
// playthroughs[0][finish_date]:
// playthroughs[0][medium_id]: 0
// playthroughs[0][played_platform]: 9
// log[is_play]: false
// log[is_playing]: true
// log[is_backlog]: false
// log[is_wishlist]: false
// log[status]: completed
// log[id]:
// log[library_entries][-5][platform_id]: 9
// log[library_entries][-5][ownership_id]: 1
// log[total_hours]:
// log[total_minutes]:
// log[time_source]: 1
// modal_type: full

type CreateLogReq struct {
	GameID         string
	PlayThroughID  string // in case of new playthrough, pass -1
	Title          string
	Rating         int
	Review         string
	ReviewSpoilers bool
	Platform       int
	Hours          int
	Minutes        int
	IsMaster       bool
	IsReplay       bool
	// StartDate      string
	// FinishDate     string
	MediumID       int
	PlayedPlatform int
	IsPlay         bool
	IsPlaying      bool
	IsBacklog      bool
	IsWishlist     bool
	Status         string
	// LibraryEntiresPlatformID  int
	// LibraryEntiresOwnershipID int
	TotalHours   float64
	TotalMinutes int
	TimeSource   int
}


// This breaks the page of backloggd oops

func (sdk *BackloggdSDK) CreateLog(cl CreateLogReq) error {
	form := url.Values{}
	form.Set("game_id", cl.GameID)
	form.Set("playthroughs[0][id]", cl.PlayThroughID)
	form.Set("playthroughs[0][title]", cl.Title)
	form.Set("playthroughs[0][rating]", fmt.Sprintf("%d", cl.Rating))
	form.Set("playthroughs[0][review]", cl.Review)
	form.Set("playthroughs[0][review_spoilers]", fmt.Sprintf("%t", cl.ReviewSpoilers))
	form.Set("playthroughs[0][platform]", fmt.Sprintf("%d", cl.Platform))
	form.Set("playthroughs[0][hours]", fmt.Sprintf("%d", cl.Hours))
	form.Set("playthroughs[0][minutes]", fmt.Sprintf("%d", cl.Minutes))
	form.Set("playthroughs[0][is_master]", fmt.Sprintf("%t", cl.IsMaster))
	form.Set("playthroughs[0][is_replay]", fmt.Sprintf("%t", cl.IsReplay))
	// form.Set("playthroughs[0][start_date]", cl.StartDate)
	// form.Set("playthroughs[0][finish_date]", cl.FinishDate)
	form.Set("playthroughs[0][medium_id]", fmt.Sprintf("%d", cl.MediumID))
	form.Set("playthroughs[0][played_platform]", fmt.Sprintf("%d", cl.PlayedPlatform))
	form.Set("log[is_play]", fmt.Sprintf("%t", cl.IsPlay))
	form.Set("log[is_playing]", fmt.Sprintf("%t", cl.IsPlaying))
	form.Set("log[is_backlog]", fmt.Sprintf("%t", cl.IsBacklog))
	form.Set("log[is_wishlist]", fmt.Sprintf("%t", cl.IsWishlist))
	form.Set("log[status]", cl.Status)
	// form.Set("log[id]", cl.LogID)
	// form.Set("log[library_entries][-5][platform_id]", string(cl.LibraryEntiresPlatformID))
	// form.Set("log[library_entries][-5][ownership_id]", string(cl.LibraryEntiresOwnershipID))
	form.Set("log[total_hours]", fmt.Sprintf("%f", cl.TotalHours))
	form.Set("log[total_minutes]", fmt.Sprintf("%d", cl.TotalMinutes))
	form.Set("log[time_source]", fmt.Sprintf("%d", cl.TimeSource))
	form.Set("modal_type", "full")

	data := form.Encode()
	req, err := http.NewRequest("POST", fmt.Sprintf(creatLogAPIURL, sdk.UserID, cl.GameID), strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL)
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
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create log: %d", resp.StatusCode)
	}

	return nil
}
