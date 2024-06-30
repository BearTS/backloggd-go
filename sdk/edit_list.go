package sdk

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type EditListReq struct {
	Name                        *string
	Privacy                     *string
	Ranked                      *bool
	Style                       *string // details or grid
	DefaultListSorting          *string
	DefaultListSortingDirection *int // 0 for desc and 1 for asc
	Description                 *string
	EditOrder                   []ListGameDetails // need to pass the full list of games in the order you want them to be, emit the ones you wish to be removed
	AddGameByGameId             []string
}

func (sdk *BackloggdSDK) EditList(listDetails ListDetails, editReq EditListReq) (ListDetails, error) {
	var err error
	if len(editReq.EditOrder) > 0 || len(editReq.AddGameByGameId) > 0 {
		err = sdk.EditListGamesOrder(editReq.EditOrder, editReq.AddGameByGameId, listDetails)
		if err != nil {
			return listDetails, err
		}
	}

	form := url.Values{}

	if editReq.Name != nil {
		form.Set("list[name]", *editReq.Name)
	}
	if editReq.Privacy != nil {
		form.Set("list[privacy]", *editReq.Privacy)
	}
	if editReq.Ranked != nil {
		form.Set("list[ranked]", fmt.Sprintf("%t", *editReq.Ranked))
	}

	if editReq.Style != nil {
		form.Set("list[style]", *editReq.Style)
	}

	if editReq.DefaultListSorting != nil {
		form.Set("default_list_sorting", *editReq.DefaultListSorting)
	}

	if editReq.DefaultListSortingDirection != nil {
		form.Set("default_list_sorting_dir", fmt.Sprintf("%d", *editReq.DefaultListSortingDirection))
	}

	if editReq.Description != nil {
		form.Set("list[desc]", *editReq.Description)
	}

	// check if the form is empty
	if len(form) == 0 {
		listDetails, err = sdk.GetListDetails(listDetails.Slug)
		if err != nil {
			return listDetails, err
		}
		return listDetails, nil
	}

	form.Set("utf8", "✓")
	form.Set("_method", "put")
	form.Set("authenticity_token", listDetails.AuthencityToken)
	data := form.Encode()

	// Edit the list
	req, err := http.NewRequest("POST", fmt.Sprintf(editListApiURL, listDetails.ID), strings.NewReader(data))
	if err != nil {
		return listDetails, err
	}

	req.Header.Set("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", fmt.Sprintf(editListURL, sdk.Username, listDetails.Slug))
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-CSRF-Token", sdk.csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return listDetails, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listDetails, err
	}
	listLink := extractVisitURL(string(body))
	linkSlug := extractLinkSlugFromLink(listLink)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return listDetails, fmt.Errorf("edit list failed with status code %d", resp.StatusCode)
	}

	listDetails, err = sdk.GetListDetails(linkSlug)
	if err != nil {
		return listDetails, err
	}

	return listDetails, nil
}

func (sdk *BackloggdSDK) EditListGamesOrder(orderedArray []ListGameDetails, addGameID []string, originalList ListDetails) error {
	var data string

	lengthOfOrderedArray := len(orderedArray)
	for i, game := range originalList.CurrentOrder {
		if i != 0 {
			if data[len(data)-1] != '&' {
				data += "&"
			}
		}
		data += fmt.Sprintf("prev_order%%5B%%5D=%s", game.EntryID)
	}

	if lengthOfOrderedArray > 0 {
		for _, game := range orderedArray {
			if data[len(data)-1] != '&' {
				data += "&"
			}
			data += fmt.Sprintf("new_order%%5B%%5D=%s", game.EntryID)
		}
	} else {
		for _, game := range originalList.CurrentOrder {
			if data[len(data)-1] != '&' {
				data += "&"
			}
			data += fmt.Sprintf("new_order%%5B%%5D=%s", game.EntryID)
		}
	}

	for i, gameID := range addGameID {
		if data[len(data)-1] != '&' {
			data += "&"
		}
		data += fmt.Sprintf("new_entries%%5B-%d%%5D=%s", i+1, gameID)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf(editListApiURL, originalList.ID)+"update-entries-2/", strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", fmt.Sprintf(editListURL, sdk.Username, originalList.Slug))
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-CSRF-Token", originalList.CsrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("edit list failed with status code %d", resp.StatusCode)
	}

	return nil
}
