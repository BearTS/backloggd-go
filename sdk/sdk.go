package sdk

import (
	"errors"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"go.nhat.io/cookiejar"
	"golang.org/x/net/publicsuffix"
)

var (
	baseURL             = "https://backloggd.com"
	signInURL           = baseURL + "/users/sign_in"
	settingsURL         = baseURL + "/settings"
	usersURL            = baseURL + "/users/%s"                         // %s = username
	autocompleteJsonURL = baseURL + "/autocomplete.json"                // + "?query=" + query
	userGamesURL        = baseURL + "/u/" + "%s" + "/games/%s/type:%s/" // %s = username, %s = sort, %s = type
	gamesURL            = baseURL + "/games/%s/"                        // %s = slug of the game
	logURL              = baseURL + "/log/"
	logStatusURL        = baseURL + "/log/status"
)

// BackloggdSDK provides methods to interact with the Backloggd website
type BackloggdSDK struct {
	Client    *http.Client
	Jar       *cookiejar.PersistentJar
	username  string
	userID    string
	csrfToken string
}

// NewBackloggdSDK creates a new instance of the Backloggd SDK
func NewBackloggdSDK(username, password string) (*BackloggdSDK, error) {
	var c BackloggdSDK
	cookiesFile := "cookies.json"
	file, err := os.Open(cookiesFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(cookiesFile)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer file.Close()

	jar := cookiejar.NewPersistentJar(
		cookiejar.WithFilePath(cookiesFile),
		cookiejar.WithAutoSync(true),
		cookiejar.WithPublicSuffixList(publicsuffix.List),
	)

	c.Client = &http.Client{
		Jar: jar,
	}
	c.Jar = jar
	c.username = username

	setUserId := func(c *BackloggdSDK) error {
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

		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		c.csrfToken, _ = doc.Find("meta[name='csrf-token']").Attr("content")

		doc.Find("button#save-profile-btn").Each(func(i int, s *goquery.Selection) {
			c.userID, _ = s.Attr("user_id")
		})

		if c.userID == "" {
			return errors.New("could not find user id")
		}
		if c.csrfToken == "" {
			return errors.New("could not find csrf token")
		}

		return nil
	}

	err = setUserId(&c)
	if err != nil {
		// try login
		err = c.Login(username, password)
		if err != nil {
			return nil, err
		}
		err = setUserId(&c)
		if err != nil {
			return nil, errors.New("could not initialize user, tried to login but failed")
		}
	}

	return &c, nil
}
