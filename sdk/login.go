package sdk

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Login performs the login to Backloggd using the provided credentials
func (sdk *BackloggdSDK) Login(username, password string) error {
	// Step 1: Perform GET request to obtain initial session and HTML

	req, err := http.NewRequest("GET", signInURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Referer", signInURL)
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

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

	// Find the authenticity_token
	var authenticityToken string
	doc.Find("form[action='/users/sign_in'] input[name='authenticity_token']").Each(func(i int, s *goquery.Selection) {
		authenticityToken, _ = s.Attr("value")
	})

	// Construct form data
	form := url.Values{}
	form.Set("authenticity_token", authenticityToken)
	form.Set("user[login]", username)
	form.Set("user[password]", password)
	form.Set("user[remember_me]", "0") // Remember me
	form.Add("user[remember_me]", "1") // Remember me
	form.Set("utf8", "✓")              // Check symbol URL encoded
	form.Set("commit", "")

	postData := form.Encode()

	// Step 2: Perform POST request to sign in
	req, err = http.NewRequest("POST", signInURL, strings.NewReader(postData))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Referer", signInURL)
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = sdk.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}

	sdk.username = username

	// Successfully logged in
	return nil
}
