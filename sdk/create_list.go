package sdk

import (
	"io"
	"net/http"
	"strings"
)

func (sdk *BackloggdSDK) CreateList() (string, string, error) {
	var listLink string
	var linkSlug string
	req, err := http.NewRequest("POST", newListURL, nil)
	if err != nil {
		return listLink, linkSlug, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-CSRF-Token", sdk.csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return listLink, linkSlug, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listLink, linkSlug, err
	}

	listLink = extractVisitURL(string(body))
	linkSlug = extractLinkSlugFromLink(listLink)
	return listLink, linkSlug, err
}

func extractVisitURL(body string) string {
	// Find index of "Turbolinks.visit("
	index := strings.Index(body, `Turbolinks.visit("`)
	if index == -1 {
		return "" // Return empty string if "Turbolinks.visit(" not found
	}

	// Move index to start after "Turbolinks.visit("
	index += len(`Turbolinks.visit("`)

	// Find closing quote after the URL
	endIndex := strings.Index(body[index:], `"`)
	if endIndex == -1 {
		return "" // Return empty string if closing quote not found
	}

	// Extract the URL
	url := body[index : index+endIndex]
	return url
}

func extractLinkSlugFromLink(link string) string {
	// Find index of "/list/"
	index := strings.Index(link, "/list/")
	if index == -1 {
		return "" // Return empty string if "/list/" not found
	}

	// Move index to start after "/list/"
	index += len("/list/")

	// Find closing quote after the URL
	endIndex := strings.Index(link[index:], "/")
	if endIndex == -1 {
		return "" // Return empty string if closing quote not found
	}

	// Extract the URL
	url := link[index : index+endIndex]
	return url
}
