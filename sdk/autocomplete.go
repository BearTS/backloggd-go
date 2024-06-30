package sdk

import (
	"encoding/json"
	"net/http"
)

type AutoComplete struct {
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
			Year  string `json:"year"`
			ID    int    `json:"id"`
		} `json:"data"`
	} `json:"suggestions"`
}

// Returns a list of autocomplete suggestions based on the query
// Example usage: data, err := client.Autocomplete("Spiderman")
func (sdk *BackloggdSDK) Autocomplete(query string) (*AutoComplete, error) {
	var ac AutoComplete
	// Step 1: Perform GET request to obtain autocomplete JSON
	req, err := http.NewRequest("GET", autocompleteJsonURL+"?query="+query, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("X-CSRF-Token", sdk.csrfToken)
	req.Header.Set("Referer", settingsURL)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ac)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}
