package sdk

import (
	"net/http"
	"os"
	"path/filepath"

	"go.nhat.io/cookiejar"
	"golang.org/x/net/publicsuffix"
)

// BackloggdSDK provides methods to interact with the Backloggd website
type BackloggdSDK struct {
	Client *http.Client
	Jar    *cookiejar.PersistentJar
}

// NewBackloggdSDK creates a new instance of the Backloggd SDK
func NewBackloggdSDK() (*BackloggdSDK, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "example")
	if err != nil {
		return nil, err
	}

	cookiesFile := filepath.Join(tempDir, "cookies")

	jar := cookiejar.NewPersistentJar(
		cookiejar.WithFilePath(cookiesFile),
		cookiejar.WithAutoSync(true),
		cookiejar.WithPublicSuffixList(publicsuffix.List),
	)

	client := &http.Client{
		Jar: jar,
	}

	return &BackloggdSDK{
		Client: client,
		Jar:    jar,
	}, nil
}
