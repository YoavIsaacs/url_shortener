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
	"github.com/YoavIsaacs/url_shortener/internal/internal/db/sqlc"
	"github.com/google/uuid"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server is working properly")
}

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
		w.Write([]byte("error: could not create new id for this url..."))
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
		w.Write([]byte("error: error decoding json"))
		return
	}

	shortenedDomain := createShortDomain(receivedParamsDecoded.OriginalDomain)

	paramsToSend := sqlc.CreateNewURLParams{
		ID:             paramId,
		CreatedAt:      paramTime,
		UpdatedAt:      paramTime,
		OriginalDomain: receivedParamsDecoded.OriginalDomain,
		ShortDomain:    shortenedDomain,
	}

	ctx := r.Context()

	createdURL, err := c.Database.CreateNewURL(ctx, paramsToSend)
	if err != nil {
		fmt.Printf("error: error creating new URL: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("error: error creating new URL"))
		return
	}

	responseData, err := json.Marshal(createdURL)
	if err != nil {
		fmt.Printf("error: error decoding response: %s", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("error: error decoding response: %s"))
		return
	}

	fmt.Printf("created new url: %+v\n", createdURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(responseData)
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
		w.Write([]byte("error: error reseting database...\n"))
		fmt.Println("error: error reseting database...")
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("successfully reset database"))
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
		w.Write([]byte("error: error decoding json"))
		return
	}

	ctx := r.Context()

	result, err := c.Database.DeleteSingleURL(ctx, receivedParamsDecoded.ShortDomain)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("error: error removing single url"))
		fmt.Println("error: error removing single url")
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("error: error with internal response from query"))
		fmt.Println("error: error with internal response from query")
		return
	}

	if affected == 0 {
		responseString := receivedParamsDecoded.ShortDomain + " does not exist"
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(responseString))
		fmt.Println(responseString)
		return
	}

	responseString := "successfully removed: " + receivedParamsDecoded.ShortDomain
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(responseString))
	fmt.Println(responseString)
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
