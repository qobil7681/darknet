package testing

import (
	"fmt"
	"github.com/r00tk3y/prying-deep/configs"
	"github.com/r00tk3y/prying-deep/models"
	"github.com/r00tk3y/prying-deep/pkg/logger"
	"github.com/r00tk3y/prying-deep/pkg/pryingtools/cryptoscanner"
	"github.com/r00tk3y/prying-deep/pkg/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var client *http.Client
var url string
var cryptoScanner *cryptoscanner.CryptoScanner

func TestSetup(t *testing.T) {
	configs.SetupEnvironment()
	cfg := configs.GetConfig().DbConf

	logger.InitLogger()
	defer logger.Logger.Sync()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbTestName)
	models.SetupDatabase(dbURL)
	var err error
	torProxy := fmt.Sprintf("socks5://localhost:9050")
	client, err = utils.SetupNewTorClient(torProxy)
	if err != nil {
		t.Fatal(err)
	}
	cryptoScanner = cryptoscanner.New()

}
func TestCryptoScanner(t *testing.T) {
	assert := assert.New(t)
	url = "http://btcmixer2e3pkn64eb5m65un5nypat4mje27er4ymltzshkmujmxlmyd.onion/pgp-and-bitcoin"
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	cryptoScanner.Search(string(body), 1)

	crypto, err := models.GetCrypto(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(crypto), 2)
	t.Cleanup(func() {
		models.DeleteCryptoByWebPageId(1)
	})

}