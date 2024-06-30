package sdk

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Game struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Link string `json:"link"`
}

func (sdk *BackloggdSDK) Wishlist() (*[]Game, error) {
	var games []Game
	// Step 1: Perform GET request to obtain wishlist JSON
	req, err := http.NewRequest("GET", fmt.Sprintf(wishlistURL, sdk.username), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
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
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Request.URL)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".card.mx-auto.game-cover").Each(func(i int, s *goquery.Selection) {
		game := Game{}
		game.Name = s.Find(".card-img").AttrOr("alt", "")
		game.Id, _ = strconv.Atoi(s.AttrOr("game_id", "0"))
		game.Link, _ = s.Find(".cover-link").Attr("href")
		games = append(games, game)
	})

	return &games, nil
}
