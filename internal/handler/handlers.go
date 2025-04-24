package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/YoavIsaacs/url_shortener/internal/config"
	"github.com/YoavIsaacs/url_shortener/internal/internal/db/sqlc"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server is working properly")
}

func createShortDomain(original_domain string) string {
	hasher := sha256.New()
	hasher.Write([]byte(original_domain))
	hashBytes := hasher.Sum(nil)

	base64Str := base64.URLEncoding.EncodeToString(hashBytes)
	shortened_domain := base64Str[:8]

	shortened_domain = strings.ReplaceAll(shortened_domain, "+", "a")
	shortened_domain = strings.ReplaceAll(shortened_domain, "/", "b")
	shortened_domain = strings.ReplaceAll(shortened_domain, "=", "c")

	return shortened_domain
}

func addURL(c config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type ExpectedData struct {
		OriginalDomain string `json:"original_domain"`
	}

	if r.Method != http.MethodPost {
		return
	}

	paramId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error: could not create new id for this url...")
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: could not create new id for this url..."))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	paramTime := time.Now()

	decoder := json.NewDecoder(r.Body)
	receivedParamsDecoded := ExpectedData{}

	err = decoder.Decode(&receivedParamsDecoded)
	if err != nil {
		fmt.Printf("error: error decoding json: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error decoding json"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	shortenedDomain := createShortDomain(receivedParamsDecoded.OriginalDomain)

	ctx := r.Context()
	// chekc if the url already exists
	checkOriginalDomain, err := c.Database.GetURLViaShortURL(ctx, shortenedDomain)
	if err == nil {
		type responseJson struct {
			Message string `json:"msg"`
			URL     string `json:"url"`
		}

		responseString := "the short url for " + receivedParamsDecoded.OriginalDomain + " already exists"

		paramsPreMarshal := responseJson{
			Message: responseString,
			URL:     checkOriginalDomain,
		}

		params, err := json.Marshal(paramsPreMarshal)
		if err != nil {
			fmt.Printf("error: error marshalling repsosne data: %s", err)
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, err := w.Write([]byte("error: error marshalling repsosne data"))
			if err != nil {
				fmt.Println("Failed to write response:", err)
				return
			}
			return
		}

		fmt.Println(responseString)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(params)
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	paramsToSend := sqlc.CreateNewURLParams{
		ID:             paramId,
		CreatedAt:      paramTime,
		UpdatedAt:      paramTime,
		OriginalDomain: receivedParamsDecoded.OriginalDomain,
		ShortDomain:    shortenedDomain,
	}

	createdURL, err := c.Database.CreateNewURL(ctx, paramsToSend)
	if err != nil {
		fmt.Printf("error: error creating new URL: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error creating new URL"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	responseData, err := json.Marshal(createdURL)
	if err != nil {
		fmt.Printf("error: error decoding response: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error decoding response: %s"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	fmt.Printf("created new url: %+v\n", createdURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_, err = w.Write(responseData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		return
	}
}

func deleteAllURLs(c config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	ctx := r.Context()
	err := c.Database.DeleteAllURLs(ctx)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error reseting database...\n"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error reseting database...")
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte("successfully reset database"))
	if err != nil {
		fmt.Println("Failed to write response:", err)
		return
	}
	fmt.Println("successfully reset database")
}

func deleteSingleURL(c config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type ExpectedData struct {
		ShortDomain string `json:"short_domain"`
	}

	if r.Method != http.MethodPost {
		return
	}

	decoder := json.NewDecoder(r.Body)

	receivedParamsDecoded := ExpectedData{}

	err := decoder.Decode(&receivedParamsDecoded)
	if err != nil {
		fmt.Printf("error: error decoding json: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error decoding json"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}

	ctx := r.Context()

	result, err := c.Database.DeleteSingleURL(ctx, receivedParamsDecoded.ShortDomain)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error removing single url"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error removing single url")
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error with internal response from query"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error with internal response from query")
		return
	}

	if affected == 0 {
		responseString := receivedParamsDecoded.ShortDomain + " does not exist"
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte(responseString))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println(responseString)
		return
	}

	responseString := "successfully removed: " + receivedParamsDecoded.ShortDomain
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(responseString))
	if err != nil {
		fmt.Println("Failed to write response:", err)
		return
	}
	fmt.Println(responseString)
}

func handleRedirect(c config.ApiConfig, short string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	originalDomain, err := c.Database.GetURLViaShortURL(ctx, short)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error with internal response from query"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error with internal response from query")
		return
	}

	if originalDomain == "" {
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: this short URL does not exist yet"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: this short URL does not exist yet")
		return
	}
	if !strings.HasPrefix(originalDomain, "http://") && !strings.HasPrefix(originalDomain, "https://") {
		originalDomain = "https://" + originalDomain
	}

	http.Redirect(w, r, originalDomain, http.StatusFound)
	fmt.Printf("successfully redirected to: %s", originalDomain)
	c.Database.AddOneToHits(ctx, short)
}

func handleCreateQRCode(c config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type expected struct {
		Short string `json:"short_domain"`
	}

	decoder := json.NewDecoder(r.Body)
	receivedParamsDecoded := expected{}

	err := decoder.Decode(&receivedParamsDecoded)
	if err != nil {
		fmt.Printf("error: error decoding json: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error decoding json"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		return
	}
	ctx := r.Context()
	short := receivedParamsDecoded.Short

	originalDomain, err := c.Database.GetURLViaShortURL(ctx, short)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error with internal response from query"))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error with internal response from query")
		return
	}

	png, err := createQRCode(originalDomain)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error creating qr code..."))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error creating qr code...")
		return
	}

	params := sqlc.AddQRCodeParams{
		ShortDomain: short,
		QrCode:      png,
	}

	qr, err := c.Database.AddQRCode(ctx, params)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte("error: error storing qr code..."))
		if err != nil {
			fmt.Println("Failed to write response:", err)
			return
		}
		fmt.Println("error: error storing qr code...")
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(qr.QrCode)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		return
	}
}

func createQRCode(originalDomain string) ([]byte, error) {
	var png []byte
	originalDomain = ensureHTTPS(originalDomain)
	png, err := qrcode.Encode(originalDomain, qrcode.Medium, 256)
	if err != nil {
		return []byte{}, err
	}
	return png, nil
}

func ensureHTTPS(url string) string {
	if strings.HasPrefix(url, "https://") {
		return url
	}
	// also handle http:// if you want to normalize everything to https
	if strings.HasPrefix(url, "http://") {
		return "https://" + strings.TrimPrefix(url, "http://")
	}
	return "https://" + url
}

func HandleAddQR(c config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleCreateQRCode(c, w, r)
	}
}

func HandleRedirect(c config.ApiConfig, u string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRedirect(c, u, w, r)
	}
}

func HandleDeleteSingleURL(c config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteSingleURL(c, w, r)
	}
}

func HandleDeleteAllURLs(c config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteAllURLs(c, w, r)
	}
}

func HandleAddURL(c config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addURL(c, w, r)
	}
}
