package handler

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/YoavIsaacs/url_shortener/internal/config"
	"github.com/google/uuid"
)

func createShortDomain(original_domain string) string {
	hasher := md5.New()
	hasher.Write([]byte(original_domain))
	hashBytes := hasher.Sum(nil)

	base64Str := base64.URLEncoding.EncodeToString(hashBytes)
	shortened_domain := base64Str[:8]

	shortened_domain = strings.ReplaceAll(shortened_domain, "+", "a")
	shortened_domain = strings.ReplaceAll(shortened_domain, "/", "b")
	shortened_domain = strings.ReplaceAll(shortened_domain, "=", "c")

	return shortened_domain
}

func HandleAddURL(c config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type ExpectedData struct {
		Created_at      time.Time `json:"created_at"`
		Updated_at      time.Time `json:"updated_at"`
		Original_domain string    `json:"original_domain"`
	}

	type CreatedRow struct {
		Id               uuid.UUID `json:"id"`
		Created_at       time.Time `json:"created_at"`
		Updated_at       time.Time `json:"updated_at"`
		Original_domain  string    `json:"original_domain"`
		Shortened_domain string    `json:"shortened_domain"`
		Hits             int       `json:"hits"`
	}

	type Params struct {
		Id               uuid.UUID `json:"id"`
		Created_at       time.Time `json:"created_at"`
		Updated_at       time.Time `json:"updated_at"`
		Original_domain  string    `json:"original_domain"`
		Shortened_domain string    `json:"shortened_domain"`
		Hits             int       `json:"hits"`
	}

	if r.Method != http.MethodPost {
		return
	}

	paramId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error: could not create new id for this url...")
		w.WriteHeader(500)
		w.Write([]byte("error: could not create new id for this url..."))
	}

	paramTime := time.Now()

	decoder := json.NewDecoder(r.Body)

	receivedParamsDecoded := ExpectedData{}

	err = decoder.Decode(&receivedParamsDecoded)
	if err != nil {
		fmt.Printf("error: error decoding json: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("error: error decoding json"))
		return
	}

	shortened_domain := createShortDomain(receivedParamsDecoded.Original_domain)

	paramsToSend := Params{
		Id:               paramId,
		Created_at:       paramTime,
		Updated_at:       paramTime,
		Original_domain:  receivedParamsDecoded.Original_domain,
		Shortened_domain: shortened_domain,
		Hits:             0,
	}

	ctx := r.Context()
}
