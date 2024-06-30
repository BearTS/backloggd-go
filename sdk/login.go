package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseURL   = "https://backloggd.com"
	signInURL = baseURL + "/users/sign_in"
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

	// Prepare POST data
	postData := strings.NewReader(
		"authenticity_token=" + authenticityToken +
			"&user%5Blogin%5D=" + username +
			"&user%5Bpassword%5D=" + password +
			"&user%5Bremember_me%5D=0" + // Remember me
			"&user%5Bremember_me%5D=1" + // Remember me
			"&utf8=%E2%9C%93" + // ✓ symbol URL encoded
			"&commit=",
	)

	// Step 2: Perform POST request to sign in
	req, err = http.NewRequest("POST", signInURL, postData)
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

	// Successfully logged in
	return nil
}

// ExportCookies exports all cookies currently stored in the SDK's cookie jar to a JSON file
func (sdk *BackloggdSDK) ExportCookies() error {
	// Parse URL
	parsedUrl, err := netUrl.Parse(baseURL)
	if err != nil {
		return err
	}
	cookies := sdk.Jar.Cookies(parsedUrl)

	// Export cookies to a JSON file
	cookiesFile := "cookies.json"
	file, err := os.Create(cookiesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, cookie := range cookies {
		if err := json.NewEncoder(file).Encode(cookie); err != nil {
			return err
		}
	}

	fmt.Printf("Cookies exported to %s\n", cookiesFile)
	return nil
}
