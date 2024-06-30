package sdk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/BearTS/backloggd-go/enums"
	"github.com/PuerkitoBio/goquery"
)

type LogStatusReq struct {
	Slug   string
	Status enums.GameStatus
}

func (sdk *BackloggdSDK) LogStatus(sr LogStatusReq) error {
	req, err := http.NewRequest("GET", fmt.Sprintf(gamesURL, sr.Slug), nil)
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

	var gameID string
	doc.Find(".col" + PlayedButton.String()).Each(func(i int, s *goquery.Selection) {
		gameID = s.Find("button").AttrOr("game_id", "")
		if !s.HasClass("btn-play-fill") {
			err = sdk.LogRequest(Played, gameID)
			if err != nil {
				return
			}
		}
	})

	// Now change the status
	var reqByte = []byte(fmt.Sprintf("game_id=%s&status_id=%d", gameID, sr.Status.Int()))

	// POST To get the status
	req, err = http.NewRequest("PATCH", logStatusURL, bytes.NewBuffer(reqByte))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", fmt.Sprintf(gamesURL, sr.Slug))
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("X-CSRF-Token", sdk.csrfToken)

	resp, err = sdk.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to log status: %d", resp.StatusCode)
	}

	return nil
}
