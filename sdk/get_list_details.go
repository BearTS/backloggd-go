package sdk

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type ListDetails struct {
	Slug                    string
	ID                      string
	AuthencityToken         string // exposing it for editing
	CsrfToken               string // exposing it for editing
	Name                    string
	Description             string
	Privacy                 string
	Ranked                  bool
	Style                   string
	DefaultSorting          string
	DefaultSortingDirection int //0 for desc and 1 for asc
	CurrentOrder            []ListGameDetails
}

type ListGameDetails struct {
	GameID  string
	EntryID string
	Note    string
	Name    string
}

func (sdk *BackloggdSDK) GetListDetails(slug string) (ListDetails, error) {
	var listDetails ListDetails
	requestUrl := fmt.Sprintf(editListURL, sdk.Username, slug)

	// Get Authencity token
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return listDetails, err
	}

	req.Header.Set("Accept", "text/html, application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", requestUrl)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Turbolinks-Referrer", requestUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return listDetails, err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return listDetails, err
	}

	listDetails.CsrfToken, _ = doc.Find("meta[name='csrf-token']").Attr("content")
	listDetails.ID = doc.Find("#list-submit").First().AttrOr("list_id", "")
	if listDetails.ID == "" {
		return listDetails, fmt.Errorf("list not found")
	}
	listDetails.Slug = slug
	listDetails.AuthencityToken = doc.Find("form[action='/api/list/"+listDetails.ID+"']"+" input[name='authenticity_token']").First().AttrOr("value", "")
	listDetails.Name = doc.Find("input[name='list[name]']").First().AttrOr("value", "")
	listDetails.Description = doc.Find("textarea[name='list[desc]']").First().Text()
	if doc.Find("#list_privacy").First().HasClass("checked") {
		listDetails.Privacy = "public"
	} else {
		listDetails.Privacy = "private"

	}
	listDetails.Ranked = doc.Find("#list_ranked").First().HasClass("checked")
	if doc.Find("#list_style").First().HasClass("checked") {
		listDetails.Style = "grid"
	} else {
		listDetails.Style = "detail"

	}
	listDetails.DefaultSorting = doc.Find("#default_list_sorting").First().AttrOr("selected-value", "")
	listDetails.DefaultSortingDirection, _ = strconv.Atoi(doc.Find("#default_list_sorting_dir").First().AttrOr("selected-value", "0"))

	doc.Find("#list-grid .grid-list-entry").Each(func(i int, s *goquery.Selection) {
		var gameDetails ListGameDetails
		gameDetails.GameID = s.AttrOr("game_id", "")
		gameDetails.EntryID = s.AttrOr("entry_id", "")
		gameDetails.Note = s.AttrOr("note", "")
		gameDetails.Name = s.Find(".card-img").First().AttrOr("alt", "")
		listDetails.CurrentOrder = append(listDetails.CurrentOrder, gameDetails)
	})

	return listDetails, nil
}

// <select data-icon-base="fas" class="selectpicker slim-selectpicker" id="default_list_sorting" title="Sort method" data-size="6" name="default_list_sorting" data-width="100%" selected-value="user">
// <option value="user" id="user" default_dir="desc">User Order</option>
// <option value="title" id="title" default_dir="asc">Game Title</option>
// <option value="popularity" id="popularity" default_dir="desc">Popularity</option>
// <option value="release" id="release" default_dir="desc">Release Date</option>
// <option value="trending" id="trending" default_dir="desc">Trending</option>
// <option value="user_rating" id="user_rating" default_dir="desc">User&#39;s Rating</option>
// <option value="rating" id="rating" default_dir="desc">Avg Rating</option>
