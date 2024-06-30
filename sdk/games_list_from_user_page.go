package sdk

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// <div class="row mx-0 my-2">
// 	<nav class="pagy-nav pagination"><span class="page prev disabled">&lsaquo;&nbsp;Prev</span> <span class="page active">1</span> <span class="page"><a href="/u/stovetop/games/title:asc/type:played?page=2"   rel="next" >2</a></span> <span class="page"><a href="/u/stovetop/games/title:asc/type:played?page=3"   >3</a></span> <span class="page"><a href="/u/stovetop/games/title:asc/type:played?page=4"   >4</a></span> <span class="page gap">&hellip;</span> <span class="page"><a href="/u/stovetop/games/title:asc/type:played?page=25"   >25</a></span> <span class="page next"><a href="/u/stovetop/games/title:asc/type:played?page=2"   rel="next" aria-label="next">Next&nbsp;&rsaquo;</a></span></nav>
// </div>

type GamesListType int

const (
	UserGamesListTypePlayed GamesListType = iota
	UserGamesListTypePlaying
	UserGamesListTypeBacklog
	UserGamesListTypeWishlist
)

func (t GamesListType) String() string {
	switch t {
	case UserGamesListTypePlayed:
		return "played"
	case UserGamesListTypePlaying:
		return "playing"
	case UserGamesListTypeBacklog:
		return "backlog"
	case UserGamesListTypeWishlist:
		return "wishlist"
	default:
		return UserGamesListTypeWishlist.String()
	}
}

type UserGamesQueryPageSort int

const (
	UserGamesQueryPageSortGameTitle UserGamesQueryPageSort = iota
	UserGamesQueryPageSortWhenAdded
	UserGamesQueryPageSortTrending
	UserGamesQueryPageSortReleaseDate
	UserGamesQueryPageSortUserRating
	UserGamesQueryPageSortPopularity
	UserGamesQueryPageSortAverageRating
	UserGamesQueryPageSortTimePlayed
	UserGamesQueryPageSortRandom
)

func (s UserGamesQueryPageSort) String() string {
	switch s {
	case UserGamesQueryPageSortGameTitle:
		return "title"
	case UserGamesQueryPageSortWhenAdded:
		return "added"
	case UserGamesQueryPageSortTrending:
		return "popular"
	case UserGamesQueryPageSortReleaseDate:
		return "release"
	case UserGamesQueryPageSortUserRating:
		return "user-rating"
	case UserGamesQueryPageSortPopularity:
		return "played"
	case UserGamesQueryPageSortAverageRating:
		return "rating"
	case UserGamesQueryPageSortTimePlayed:
		return "time-played"
	case UserGamesQueryPageSortRandom:
		return "shuffle"
	default:
		return UserGamesQueryPageSortWhenAdded.String()
	}
}

type GamesQueryReq struct {
	Username string
	ListType GamesListType
	PageSort UserGamesQueryPageSort
}

func (sdk *BackloggdSDK) GetGamesListFromUserPage(queryReq GamesQueryReq) (*[]Game, error) {
	var games []Game
	var doc *goquery.Document
	requestUrl := fmt.Sprintf(userGamesURL, queryReq.Username, queryReq.PageSort.String(), queryReq.ListType.String())
	for {
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			return nil, err
		}
		fmt.Println("Requesting", requestUrl)

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

		// Load the HTML document
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		fmt.Println("Parsing", requestUrl)

		doc.Find(".card.mx-auto.game-cover").Each(func(i int, s *goquery.Selection) {
			game := Game{}
			game.Name = s.Find(".card-img").AttrOr("alt", "")
			game.Id, _ = strconv.Atoi(s.AttrOr("game_id", "0"))
			game.Link, _ = s.Find(".cover-link").Attr("href")
			games = append(games, game)
		})
		nextPage := doc.Find(".page.next")
		if nextPage.HasClass("disabled") {
			break
		}

		fmt.Println("Next page found")

		requestUrl, _ = nextPage.Find("a").Attr("href")
		if requestUrl == "" {
			fmt.Println("No next page link found")
		}
		requestUrl = baseURL + requestUrl
	}

	return &games, nil
}
