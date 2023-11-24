package testing

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/iudicium/pryingdeep/models"
	"github.com/iudicium/pryingdeep/pkg/pryingtools/cryptoscanner"
	"github.com/iudicium/pryingdeep/pkg/utils"
	"github.com/iudicium/pryingdeep/test/helpers"
)

var client *http.Client
var cryptoScanner *cryptoscanner.CryptoScanner
var db *gorm.DB

type cryptoTestCase struct {
	name string
	url  string
	iter int
	tor  bool
}

func RequestHelper(url string, useTor bool) string {
	if !useTor {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body)
}

func runCryptoTest(t *testing.T, url string, tor bool) {
	assert := assert.New(t)

	body := RequestHelper(url, tor)
	cryptoScanner.Search(body, 1)

	crypto, err := models.GetCrypto(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(len(crypto), 1)

	t.Cleanup(func() {
		models.DeleteCryptoByWebPageId(1)
	})
}

func TestSetup(t *testing.T) {
	err := helpers.InitTestConfig()
	torProxy := fmt.Sprintf("socks5://localhost:9050")
	client, err = utils.SetupNewTorClient(torProxy)
	if err != nil {
		t.Fatal(err)
	}
	cryptoScanner = cryptoscanner.NewCryptoScanner()
	db = models.GetDB()

	sql := "INSERT INTO web_pages (id, url, title, status_code, body, headers) VALUES (1, 'https://example.com', 'Example Page', 200, '<html><body>Hello, world!</body></html>', '{\"Content-Type\": \"text/html\"}')"

	result := db.Exec(sql)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

}

func TestCryptoCases(t *testing.T) {
	testCases := []cryptoTestCase{
		{
			name: "BTCAndPGPAndCertificates",
			url:  "http://btcmixer2e3pkn64eb5m65un5nypat4mje27er4ymltzshkmujmxlmyd.onion/pgp-and-bitcoin",
			tor:  true,
		},
		{
			name: "Bitcoin",
			url:  "https://blockchair.com/bitcoin/address/bc1q34aq5drpuwy3wgl9lhup9892qp6svr8ldzyy7c",
			tor:  false,
		},
		{
			name: "Eth",
			url:  "https://etherscan.io/address/0xefa06f99dfecfc0236ed7398ce57656cab732780",
			tor:  false,
		},
		{
			name: "XMR",
			url:  "https://www.getmonero.org/resources/moneropedia/address.html",
			tor:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runCryptoTest(t, tc.url, tc.tor)
		})
	}
}

func TestTearDown(t *testing.T) {
	result := db.Exec("DELETE FROM web_pages WHERE id = ?", 1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
}
