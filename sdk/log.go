package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// LogType can be "playing", "backlog" "wishlist" "play"

// buttons are Playing Played Backlog Wishlist

type LogReq struct {
	Slug    string
	LogType LogType
	Enable  bool
}

type LogType string

const (
	Played   LogType = "play"
	Playing  LogType = "playing"
	Backlog  LogType = "backlog"
	Wishlist LogType = "wishlist"
)

func (l LogType) String() string {
	return string(l)
}

func (l LogType) ButtonClass() ButtonClass {
	switch l {
	case Played:
		return PlayedButton
	case Playing:
		return PlayingButton
	case Backlog:
		return BacklogButton
	case Wishlist:
		return WishlistButton
	}

	return ""
}

// <div class="col pr-0 play-btn-container" id="play-15698">
// <div class="col px-0 playing-btn-container " id="playing-15698">
// <div class="col px-0 backlog-btn-container " id="backlog-15698">
// <div class="col pl-0 wishlist-btn-container " id="wishlist-15698">

type ButtonClass string

const (
	PlayingButton ButtonClass = ".playing-btn-container"

	// enabled if btn-play-fill present
	PlayedButton   ButtonClass = ".play-btn-container"
	BacklogButton  ButtonClass = ".backlog-btn-container"
	WishlistButton ButtonClass = ".wishlist-btn-container"
)

func (b ButtonClass) String() string {
	return string(b)
}

func (sdk *BackloggdSDK) LogGame(logReq LogReq) error {
	// Step 1: Perform GET request to get details for the game

	req, err := http.NewRequest("GET", fmt.Sprintf(gamesURL, logReq.Slug), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/html, application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", baseURL)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Turbolinks-Referrer", baseURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Step 2: now fetch the current status of the game by parsing the body
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// <div class="col-auto pl-2 pr-1 play-btn-container ml-auto" id="play-25076">
	// 	<button class="button-link btn-play btn-unplayed top-tooltip" data-tippy-content="Played" game_id="25076">
	// 		<i class="fa-kit fa-gamepad-classic"></i>
	// 	</button>
	// </div>

	doc.Find(".col" + logReq.LogType.ButtonClass().String()).Each(func(i int, s *goquery.Selection) {
		// the buttons are always present
		// if the same class contains btn-fill then it is enabled
		// game_id for when the div class has button-link and logReq.LogType.ButtonClass().String()
		gameID := s.Find("button").AttrOr("game_id", "")
		if logReq.Enable {
			if !s.HasClass("btn-play-fill") && logReq.LogType != Played {
				err = sdk.LogRequest(logReq.LogType.String(), gameID)
				if err != nil {
					return
				}
			}
		} else {
			if s.HasClass("btn-play-fill") && logReq.LogType != Played {
				err = sdk.LogRequest(logReq.LogType.String(), gameID)
				if err != nil {
					return
				}
			}
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func (sdk *BackloggdSDK) LogRequest(logType string, gameId string) error {

	// application/x-www-form-urlencoded

	// its formdata
	fmt.Println("Logging game", gameId, logType)

	var reqByte = []byte(fmt.Sprintf("type=%s&game_id=%s", logType, gameId))

	req, err := http.NewRequest("POST", logURL, bytes.NewBuffer(reqByte))
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

	fmt.Println(resp.Status)

	//
	// read the response
	var logResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&logResp)
	if err != nil {
		return err
	}

	if logResp["status"] != "completed" {
		return fmt.Errorf("failed to log the game")
	}

	return nil
}
