package sdk

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BearTS/backloggd-go/enums"
	"github.com/PuerkitoBio/goquery"
)

type GamesQueryReq struct {
	Username string
	Filter   GamesQueryFilter
	PageSort UserGamesQueryPageSort
}

type GamesQueryFilter struct {
	ListType          []GamesListType
	ReleaseYearFilter GamesQueryFilterReleaseYear
	Genre             enums.GameCategory
	Category          []enums.GameCategory
	ReleasePlatform   enums.GamesPlatform
	NoPlatformLogged  bool
	PlayedPlatform    enums.GamesPlatform
	GameStatus        enums.GameStatus
	Rating            int
}

func (sdk *BackloggdSDK) GetGamesListFromUserPage(queryReq GamesQueryReq) (*[]Game, error) {
	var games []Game
	var doc *goquery.Document
	var listType string
	// listType add "," only if there are more than one list type
	for i, t := range queryReq.Filter.ListType {
		if i > 0 {
			listType += ","
		}
		listType += t.String()
	}

	if queryReq.Filter.ListType == nil {
		listType = UserGamesListTypePlaying.String()
	}

	filter := "type:" + listType

	// Now apply other filtesr
	if queryReq.Filter.ReleaseYearFilter != 0 {
		filter += ";release_year:"
		switch queryReq.Filter.ReleaseYearFilter {
		case -1:
			filter += "released"
		case -2:
			filter += "unreleased"
		default:
			if queryReq.Filter.ReleaseYearFilter < 0 {
				return nil, fmt.Errorf("invalid release year filter")
			}
			filter += strconv.Itoa(int(queryReq.Filter.ReleaseYearFilter))
		}
	}

	if queryReq.Filter.Genre != enums.GameCategoryDefaultNone {
		filter += ";genre:" + queryReq.Filter.Genre.String()
	}

	if len(queryReq.Filter.Category) > 0 {
		filter += ";category:"
		for i, c := range queryReq.Filter.Category {
			if i > 0 {
				filter += ","
			}
			filter += c.String()
		}
	}

	if queryReq.Filter.ReleasePlatform != enums.GamePlatformDefaultNone {
		filter += ";platform:" + queryReq.Filter.ReleasePlatform.String()
	}

	if queryReq.Filter.NoPlatformLogged {
		filter += ";no_platforms_logged:true"
	}

	if queryReq.Filter.PlayedPlatform != enums.GamePlatformDefaultNone {
		filter += ";played_platform:" + queryReq.Filter.PlayedPlatform.String()
	}

	if queryReq.Filter.GameStatus != enums.GameStatusDefaultNone {
		filter += ";game_status:" + queryReq.Filter.GameStatus.String()
	}

	if queryReq.Filter.Rating > 0 && queryReq.Filter.Rating <= 10 {
		filter += ";rating:" + strconv.Itoa(queryReq.Filter.Rating)
	}

	requestUrl := fmt.Sprintf(userGamesURL, queryReq.Username, queryReq.PageSort.String(), filter)

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

type GamesQueryFilterReleaseYear int // pass -1 for released only and -2 for unreleased

type GameQueryFilterOwnership int

const (
	GameQueryFilterOwnershipDefaultNone GameQueryFilterOwnership = iota
	GameQueryFilterOwnershipPhysical
	GameQueryFilterOwnershipDigital
	GameQueryFilterOwnershipLostSold
)

func (o GameQueryFilterOwnership) String() string {
	switch o {
	case GameQueryFilterOwnershipPhysical:
		return "physical"
	case GameQueryFilterOwnershipDigital:
		return "digital"
	case GameQueryFilterOwnershipLostSold:
		return "lost"
	case GameQueryFilterOwnershipDefaultNone:
		return ""
	default:
		return GameQueryFilterOwnershipDefaultNone.String()
	}
}
