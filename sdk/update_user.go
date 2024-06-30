package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type User struct {
	Username           *string
	Bio                *string
	TwitterUrl         *string
	LetterBoxdUrl      *string
	WebsiteUrl         *string
	DisplayQuickAccess *bool
	Favourites         []Favorite
}

type Favorite struct {
	GameID int  `json:"game_id"`
	Crown  bool `json:"crown"`
}

// UpdateUser updates the user profile with the provided data
func (sdk *BackloggdSDK) UpdateUser(userData User) error {
	req, err := http.NewRequest("GET", settingsURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", signInURL)
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	csrfToken, _ := doc.Find("meta[name='csrf-token']").Attr("content")

	doc.Find("button#save-profile-btn").Each(func(i int, s *goquery.Selection) {
		sdk.userID, _ = s.Attr("user_id")
	})

	form := new(bytes.Buffer)
	w := multipart.NewWriter(form)

	if userData.Username != nil {
		w.WriteField("user[username]", *userData.Username)
	}
	if userData.Bio != nil {
		w.WriteField("user[bio]", *userData.Bio)
	}
	if userData.TwitterUrl != nil {
		w.WriteField("user[url_twitter]", *userData.TwitterUrl)
	}
	if userData.LetterBoxdUrl != nil {
		w.WriteField("user[url_letterboxd]", *userData.LetterBoxdUrl)
	}
	if userData.WebsiteUrl != nil {
		w.WriteField("user[url_website]", *userData.WebsiteUrl)
	}
	if userData.DisplayQuickAccess != nil {
		w.WriteField("user_settings[display_quick_access]", fmt.Sprintf("%t", *userData.DisplayQuickAccess))
	}
	if len(userData.Favourites) > 0 {
		// json unmarshal
		// w.WriteField("favorites", userData.Favourites)
		favs, err := json.Marshal(userData.Favourites)
		if err == nil {
			w.WriteField("favorites", string(favs))
		} else {
			return err
		}
	}

	w.Close()

	// Step 2: Perform POST request to update user
	req, err = http.NewRequest("PATCH", userURL+sdk.userID, form)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", signInURL)
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = sdk.Client.Do(req)
	if err != nil {
		return err
	}

	// print fmt.Println(resp.Status)
	fmt.Println("User updated successfully")
	fmt.Println("status: ", resp.Status)

	return nil
}
